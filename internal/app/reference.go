package app

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strings"
)

// ReferencePreparation describes edits to the copy sent to the encoder. The
// uploaded original is retained separately and is never overwritten.
type ReferencePreparation struct {
	InputSeconds    float64 `json:"inputSeconds"`
	UsedSeconds     float64 `json:"usedSeconds"`
	TrimmedSilence  bool    `json:"trimmedSilence"`
	SelectedExcerpt bool    `json:"selectedExcerpt"`
	LevelAdjusted   bool    `json:"levelAdjusted"`
	Summary         string  `json:"summary"`
}

func prepareReference(src, dst string) (AudioInfo, ReferencePreparation, error) {
	data, rate, err := pcm(src)
	var prep ReferencePreparation
	if err != nil || rate != 24000 || len(data) == 0 || len(data) > 24000*121*2 {
		return AudioInfo{}, prep, errors.New("This recording could not be prepared. Use WAV, MP3, M4A or a browser recording")
	}
	n := len(data) / 2
	prep.InputSeconds = float64(n) / float64(rate)
	if n < 3*rate {
		return AudioInfo{}, prep, fmt.Errorf("This recording is only %.1f seconds long. Record at least 3 seconds of speech; your reference is still here", prep.InputSeconds)
	}
	// Remove DC/low-frequency drift before measuring or changing volume. A
	// DC-blocking filter also handles recordings padded with zero-valued silence;
	// subtracting one whole-file mean would turn that silence into false signal.
	// This does not invent speech or reconstruct flattened/clipped peaks.
	values := make([]float64, n)
	var mean float64
	for i := range values {
		values[i] = float64(int16(binary.LittleEndian.Uint16(data[2*i:]))) / 32768
		mean += values[i]
	}
	mean /= float64(n)
	decay := math.Exp(-2 * math.Pi * 20 / float64(rate))
	previous, filtered := values[0], 0.0
	for i, raw := range values {
		filtered = raw - previous + decay*filtered
		previous = raw
		values[i] = filtered
	}
	const window = 480 // 20ms at the canonical 24kHz rate.
	frames := (n + window - 1) / window
	energy, flat := make([]float64, frames), make([]int, frames)
	var loudest float64
	run := 0
	for i, value := range values {
		energy[i/window] += value * value
		raw := int(int16(binary.LittleEndian.Uint16(data[2*i:])))
		if i > 0 {
			previous := int(int16(binary.LittleEndian.Uint16(data[2*(i-1):])))
			if absInt(raw) >= 32600 && absInt(raw-previous) <= 1 && (raw > 0) == (previous > 0) {
				run++
				if run >= 2 {
					flat[i/window]++
				}
			} else {
				run = 0
			}
		}
	}
	for i := range energy {
		energy[i] = math.Sqrt(energy[i] / float64(min(window, n-i*window)))
		loudest = math.Max(loudest, energy[i])
	}
	threshold := math.Max(.00012, math.Min(.006, loudest*.035))
	active := make([]bool, frames)
	activeCount := 0
	for i, rms := range energy {
		if rms >= threshold {
			active[i] = true
			activeCount++
		}
	}
	if activeCount < 100 { // At least two seconds of actual signal, not padding.
		return AudioInfo{}, prep, errors.New("There is not enough audible speech to prepare a reference. Check the microphone and record at least 3 seconds of clear speech; your reference is still here")
	}
	// Keep consonants and natural short pauses around the energy envelope.
	keep := make([]bool, frames)
	for i, yes := range active {
		if yes {
			for k := max(0, i-5); k <= min(frames-1, i+5); k++ {
				keep[k] = true
			}
		}
	}
	first, last := 0, frames-1
	for !keep[first] {
		first++
	}
	for !keep[last] {
		last--
	}
	// Retain enough real context for a three-second near-minimum recording.
	for min(n, (last+1)*window)-first*window < 3*rate && (first > 0 || last < frames-1) {
		if first > 0 {
			first--
		}
		if last < frames-1 {
			last++
		}
	}
	indices := make([]int, 0, last-first+1)
	for i := first; i <= last; {
		if keep[i] {
			indices = append(indices, i)
			i++
			continue
		}
		end := i
		for end <= last && !keep[end] {
			end++
		}
		count := end - i
		if count > 30 {
			count = 12
		} // Compress pauses >600ms to 240ms.
		for k := 0; k < count; k++ {
			indices = append(indices, i+k)
		}
		i = end
	}
	prep.TrimmedSilence = len(indices) < frames-2
	if len(indices) > 1500 {
		// Choose a contiguous thirty-second excerpt with little clipping and a
		// useful amount of speech. Never loop or pad missing speech.
		speechPrefix, clipPrefix := make([]int, len(indices)+1), make([]int, len(indices)+1)
		for i, frame := range indices {
			speechPrefix[i+1] = speechPrefix[i]
			if active[frame] {
				speechPrefix[i+1]++
			}
			clipPrefix[i+1] = clipPrefix[i] + flat[frame]
		}
		best, bestScore := 0, math.Inf(-1)
		for i := 0; i+1500 <= len(indices); i++ {
			speech := speechPrefix[i+1500] - speechPrefix[i]
			clipped := clipPrefix[i+1500] - clipPrefix[i]
			score := float64(speech) - float64(clipped)*20/window
			if score > bestScore {
				best, bestScore = i, score
			}
		}
		indices = indices[best : best+1500]
		prep.SelectedExcerpt = true
	}
	var selected []float64
	var clipped, signalSamples int
	for _, frame := range indices {
		end := min(n, (frame+1)*window)
		selected = append(selected, values[frame*window:end]...)
		clipped += flat[frame]
		if active[frame] {
			signalSamples += end - frame*window
		}
	}
	if len(selected) < 3*rate || signalSamples < 2*rate {
		return AudioInfo{}, prep, errors.New("Less than 3 seconds of usable speech remains after trimming. Record a little more speech; your reference is still here")
	}
	if float64(clipped)/float64(max(1, signalSamples)) > .03 {
		return AudioInfo{}, prep, errors.New("Too much of the speech is clipped or distorted. Lower the microphone input volume or move farther away, then record again; your reference is still here")
	}
	var power, peak float64
	for _, v := range selected {
		power += v * v
		peak = math.Max(peak, math.Abs(v))
	}
	rms := math.Sqrt(power / float64(len(selected)))
	gain := math.Min(100, math.Min(.1/math.Max(rms, 1e-9), .92/math.Max(peak, 1e-9)))
	if rms*gain < .005 {
		return AudioInfo{}, prep, errors.New("This recording is too faint to recover clearly. Move closer to the microphone or increase its input volume, then record again; your reference is still here")
	}
	prep.LevelAdjusted = math.Abs(20*math.Log10(gain)) > 1 || math.Abs(mean) > .001
	out := make([]byte, len(selected)*2)
	for i, v := range selected {
		// A short outer fade avoids clicks after trimming, without changing pace.
		fade := math.Min(1, math.Min(float64(i)/120, float64(len(selected)-1-i)/120))
		v = math.Max(-.98, math.Min(.98, v*gain*fade))
		binary.LittleEndian.PutUint16(out[2*i:], uint16(int16(math.Round(v*32767))))
	}
	if err := writeWAV(dst, out, rate); err != nil {
		return AudioInfo{}, prep, err
	}
	info, err := inspectWAV(dst)
	prep.UsedSeconds = info.Duration
	var edits []string
	if prep.TrimmedSilence {
		edits = append(edits, "trimmed silence")
	}
	if prep.SelectedExcerpt {
		edits = append(edits, "selected a 30-second excerpt")
	}
	if prep.LevelAdjusted {
		edits = append(edits, "balanced volume")
	}
	if len(edits) == 0 {
		prep.Summary = "Reference format prepared"
	} else {
		prep.Summary = "Automatically " + strings.Join(edits, ", ")
	}
	return info, prep, err
}

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
