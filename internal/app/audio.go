package app

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type AudioInfo struct {
	Duration float64   `json:"duration"`
	Peak     float64   `json:"peak"`
	RMS      float64   `json:"rms"`
	Clipped  float64   `json:"clipped"`
	Waveform []float64 `json:"waveform"`
	Rate     int       `json:"rate"`
}

func pcm(path string) ([]byte, int, error) {
	b, e := os.ReadFile(path)
	if e != nil {
		return nil, 0, e
	}
	if len(b) < 44 || string(b[:4]) != "RIFF" || string(b[8:12]) != "WAVE" {
		return nil, 0, errors.New("Invalid WAV audio")
	}
	var rate int
	valid := false
	for n := 12; n+8 <= len(b); {
		size := int(binary.LittleEndian.Uint32(b[n+4 : n+8]))
		if size < 0 || size > len(b)-n-8 {
			return nil, 0, errors.New("Truncated audio")
		}
		chunk := b[n+8 : n+8+size]
		switch string(b[n : n+4]) {
		case "fmt ":
			if size >= 16 {
				valid = binary.LittleEndian.Uint16(chunk) == 1 && binary.LittleEndian.Uint16(chunk[2:]) == 1 && binary.LittleEndian.Uint16(chunk[14:]) == 16
				rate = int(binary.LittleEndian.Uint32(chunk[4:]))
			}
		case "data":
			if valid && rate > 0 {
				return chunk, rate, nil
			}
		}
		n += 8 + size + (size % 2)
	}
	return nil, 0, errors.New("Use mono 16-bit PCM WAV audio")
}
func inspectWAV(path string) (AudioInfo, error) {
	data, rate, e := pcm(path)
	if e != nil {
		return AudioInfo{}, e
	}
	count := len(data) / 2
	if count == 0 {
		return AudioInfo{}, errors.New("Audio is empty")
	}
	out := AudioInfo{Rate: rate, Duration: float64(count) / float64(rate), Waveform: make([]float64, 160)}
	var power, clipped float64
	for i := 0; i < count; i++ {
		v := math.Abs(float64(int16(binary.LittleEndian.Uint16(data[2*i:]))) / 32768)
		power += v * v
		if v > out.Peak {
			out.Peak = v
		}
		if v > .99 {
			clipped++
		}
		bucket := i * 160 / count
		if v > out.Waveform[bucket] {
			out.Waveform[bucket] = v
		}
	}
	out.RMS = math.Sqrt(power / float64(count))
	out.Clipped = clipped / float64(count)
	return out, nil
}
func writeWAV(dst string, data []byte, rate int) error {
	var b bytes.Buffer
	b.WriteString("RIFF")
	_ = binary.Write(&b, binary.LittleEndian, uint32(len(data)+36))
	b.WriteString("WAVEfmt ")
	for _, v := range []any{uint32(16), uint16(1), uint16(1), uint32(rate), uint32(rate * 2), uint16(2), uint16(16)} {
		_ = binary.Write(&b, binary.LittleEndian, v)
	}
	b.WriteString("data")
	_ = binary.Write(&b, binary.LittleEndian, uint32(len(data)))
	b.Write(data)
	if e := os.WriteFile(dst+".part", b.Bytes(), 0600); e != nil {
		return e
	}
	return os.Rename(dst+".part", dst)
}
func joinWAV(files []string, out string, pause float64) error {
	var data []byte
	var rate int
	for i, p := range files {
		chunk, r, e := pcm(p)
		if e != nil {
			return e
		}
		if rate != 0 && rate != r {
			return errors.New("Segments use different sample rates")
		}
		rate = r
		// Fade only the outer 5ms of each segment to avoid hard boundary clicks.
		chunk = append([]byte(nil), chunk...)
		fade := rate / 200
		for j := 0; j < fade && j < len(chunk)/4; j++ {
			factor := float64(j) / float64(fade)
			for _, n := range []int{j, len(chunk)/2 - 1 - j} {
				v := int16(binary.LittleEndian.Uint16(chunk[2*n:]))
				binary.LittleEndian.PutUint16(chunk[2*n:], uint16(int16(float64(v)*factor)))
			}
		}
		if i > 0 {
			data = append(data, make([]byte, int(float64(rate)*pause)*2)...)
		}
		data = append(data, chunk...)
		if len(data) > 256<<20 {
			return errors.New("Recording exceeds the 256 MB export limit; split the book into chapters")
		}
	}
	return writeWAV(out, data, rate)
}
func (a *App) ffmpeg(ctx context.Context, args ...string) error {
	runtime, e := a.runtime()
	if e != nil {
		return e
	}
	all := append([]string{"-hide_banner", "-loglevel", "error", "-nostdin", "-y", "-threads", "2"}, args...)
	return a.service(ctx, filepath.Join(runtime, "ffmpeg.exe"), all, nil)
}
func (a *App) processAudio(ctx context.Context, src, dst string, s Script) error {
	pitch := math.Pow(2, s.Pitch/12)
	tempo := math.Sqrt(s.Speed / pitch)
	filters := fmt.Sprintf("aresample=24000,asetrate=24000*%.6f,aresample=24000,atempo=%.6f,atempo=%.6f,volume=%.4f", pitch, tempo, tempo, s.Volume)
	return a.ffmpeg(ctx, "-i", src, "-ac", "1", "-ar", "24000", "-af", filters, "-c:a", "pcm_s16le", "-f", "wav", dst)
}
func (a *App) convertMP3(ctx context.Context, src, dst string) error {
	return a.ffmpeg(ctx, "-i", src, "-c:a", "libmp3lame", "-b:a", "160k", "-f", "mp3", dst)
}
func (a *App) createVoice(w http.ResponseWriter, r *http.Request) {
	if !a.inference.TryLock() {
		apphost.Problem(w, 409, "Wait for the current speech or reference preparation to finish")
		return
	}
	defer a.inference.Unlock()
	r.Body = http.MaxBytesReader(w, r.Body, 26<<20)
	if e := r.ParseMultipartForm(26 << 20); e != nil {
		apphost.Problem(w, 400, "Upload audio smaller than 25 MB")
		return
	}
	defer r.MultipartForm.RemoveAll()
	name := strings.TrimSpace(r.FormValue("name"))
	model := r.FormValue("model")
	if model == "" {
		model = "pocket18"
	}
	if name == "" || len(name) > 100 || len(r.FormValue("tags")) > 300 || r.FormValue("consent") != "true" {
		apphost.Problem(w, 400, "Name this voice and confirm permission to clone it")
		return
	}
	engine := a.Engines[model]
	if engine == nil || !engine.Available() || !engine.Cloning() {
		apphost.Problem(w, 400, "Install a cloning-enabled model in Settings first")
		return
	}
	language := r.FormValue("language")
	if language == "" {
		language = "en"
	}
	if !engine.Compatible(Voice{Engine: model}, language) {
		apphost.Problem(w, 400, "Choose a reference language supported by this model")
		return
	}
	a.mu.Lock()
	var count int
	_ = a.H.DB.QueryRow(`SELECT count(*) FROM jobs WHERE status IN('running','queued')`).Scan(&count)
	if count > 0 {
		a.mu.Unlock()
		apphost.Problem(w, 409, "Wait for the current speech queue before preparing a reference")
		return
	}
	a.mu.Unlock()
	file, header, e := r.FormFile("audio")
	if e != nil {
		apphost.Problem(w, 400, "Choose reference audio")
		return
	}
	defer file.Close()
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{".wav": true, ".mp3": true, ".m4a": true, ".webm": true, ".ogg": true}
	if !allowed[ext] {
		apphost.Problem(w, 400, "Use WAV, MP3, M4A or a browser recording")
		return
	}
	id := apphost.ID()
	dir := a.path("projects", "voice-"+id)
	if e = os.MkdirAll(dir, 0700); e != nil {
		a.H.DBError(w, e)
		return
	}
	saved := false
	defer func() {
		if !saved {
			_ = os.RemoveAll(dir)
		}
	}()
	source := filepath.Join(dir, "source"+ext)
	f, e := os.OpenFile(source, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if e != nil {
		a.H.DBError(w, e)
		return
	}
	n, e := io.Copy(f, io.LimitReader(file, 25<<20+1))
	ce := f.Close()
	if e != nil || ce != nil || n > 25<<20 {
		apphost.Problem(w, 400, "Audio upload failed or exceeds 25 MB")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Minute)
	defer cancel()
	reference := filepath.Join(dir, "reference.wav")
	if e = a.ffmpeg(ctx, "-i", source, "-t", "31", "-ac", "1", "-ar", "24000", "-c:a", "pcm_s16le", reference); e != nil {
		apphost.Problem(w, 400, "This file does not contain supported audio")
		return
	}
	info, e := inspectWAV(reference)
	if e != nil || info.Duration < 3 || info.Duration > 30.25 || info.RMS < .005 || info.Clipped > .03 {
		apphost.Problem(w, 400, "Use 3–30 seconds of clear speech without silence or heavy clipping")
		return
	}
	v := Voice{ID: id, Name: name, Engine: model, Language: language, Tags: r.FormValue("tags"), Owner: uid(r), Consent: now(), Reference: "reference.wav"}
	if model == "moss" {
		v.Engine = "moss"
	}
	e = engine.Clone(ctx, reference, filepath.Join(dir, "state"))
	if e != nil {
		apphost.Problem(w, 400, e.Error())
		return
	}
	_, e = a.H.DB.Exec(`INSERT INTO voices(id,owner,body)VALUES(?,?,?)`, id, uid(r), jsonText(v))
	if e != nil {
		a.H.DBError(w, e)
		return
	}
	saved = true
	a.H.Audit(uid(r), "voice-consent-save", id)
	apphost.JSON(w, 201, map[string]any{"voice": v, "audio": info})
}
