package app

import (
	"encoding/binary"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func referenceFixture(t *testing.T, seconds float64, sample func(float64) float64) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "input.wav")
	data := make([]byte, int(seconds*24000)*2)
	for i := 0; i < len(data)/2; i++ {
		v := math.Max(-1, math.Min(1, sample(float64(i)/24000)))
		binary.LittleEndian.PutUint16(data[2*i:], uint16(int16(v*32767)))
	}
	if err := writeWAV(path, data, 24000); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestReferenceAutomaticallyPreparesQuietPaddedAudio(t *testing.T) {
	src := referenceFixture(t, 45, func(s float64) float64 {
		if s < 20 || s > 30 {
			return 0
		}
		return .0015 * math.Sin(2*math.Pi*180*s)
	})
	before, _ := os.ReadFile(src)
	dst := filepath.Join(t.TempDir(), "prepared.wav")
	info, prep, err := prepareReference(src, dst)
	if err != nil {
		t.Fatal(err)
	}
	if info.Duration < 10 || info.Duration > 10.5 || info.RMS < .05 || info.Clipped != 0 || !prep.TrimmedSilence || !prep.LevelAdjusted {
		t.Fatal(info, prep)
	}
	after, _ := os.ReadFile(src)
	if string(before) != string(after) {
		t.Fatal("original changed")
	}
}

func TestReferenceCompactsLongPausesAndSelectsExcerpt(t *testing.T) {
	src := referenceFixture(t, 80, func(s float64) float64 {
		if s > 10 && s < 20 {
			return 0
		}
		return .15 * math.Sin(2*math.Pi*180*s)
	})
	info, prep, err := prepareReference(src, filepath.Join(t.TempDir(), "out.wav"))
	if err != nil || info.Duration != 30 || !prep.SelectedExcerpt || !prep.TrimmedSilence {
		t.Fatal(info, prep, err)
	}
}

func TestReferenceDCDriftDoesNotTurnPaddingIntoSpeech(t *testing.T) {
	src := referenceFixture(t, 45, func(s float64) float64 {
		if s < 20 || s > 30 {
			return 0
		}
		return .2 + .01*math.Sin(2*math.Pi*180*s)
	})
	info, prep, err := prepareReference(src, filepath.Join(t.TempDir(), "out.wav"))
	// The DC boundary transient also constrains safe gain: recover useful
	// loudness without amplifying that transient beyond the peak limit.
	if err != nil || info.Duration > 10.8 || !prep.TrimmedSilence || info.RMS < .02 || info.Peak > .93 {
		t.Fatal(info, prep, err)
	}
}

func TestReferenceChoosesCleanSpeechInsteadOfClippedBeginning(t *testing.T) {
	src := referenceFixture(t, 70, func(s float64) float64 {
		v := math.Sin(2 * math.Pi * 180 * s)
		if s < 32 {
			return v * 8
		}
		return v * .15
	})
	_, prep, err := prepareReference(src, filepath.Join(t.TempDir(), "out.wav"))
	if err != nil || !prep.SelectedExcerpt {
		t.Fatal(prep, err)
	}
}

func TestReferenceCannotInventMissingOrHeavilyDistortedSpeech(t *testing.T) {
	for _, c := range []struct {
		name    string
		seconds float64
		sample  func(float64) float64
		message string
	}{
		{"silent", 10, func(float64) float64 { return 0 }, "not enough audible"},
		{"too-short", 1.5, func(s float64) float64 { return .1 * math.Sin(2*math.Pi*180*s) }, "only 1.5 seconds"},
		{"one-word-and-silence", 10, func(s float64) float64 {
			if s < 1 {
				return .1 * math.Sin(2*math.Pi*180*s)
			}
			return 0
		}, "not enough audible"},
		{"clipped", 10, func(s float64) float64 { return 8 * math.Sin(2*math.Pi*180*s) }, "clipped or distorted"},
	} {
		t.Run(c.name, func(t *testing.T) {
			src := referenceFixture(t, c.seconds, c.sample)
			dst := filepath.Join(t.TempDir(), "out.wav")
			_, _, err := prepareReference(src, dst)
			if err == nil || !strings.Contains(err.Error(), c.message) {
				t.Fatal(err)
			}
			if _, err := os.Stat(dst); !os.IsNotExist(err) {
				t.Fatal("unusable output published")
			}
		})
	}
}

func TestReferenceAcceptsFullScaleUnclippedWaveAndThreeSecondBoundary(t *testing.T) {
	for _, seconds := range []float64{3, 3.001, 5} {
		src := referenceFixture(t, seconds, func(s float64) float64 { return .999 * math.Sin(2*math.Pi*180*s) })
		info, _, err := prepareReference(src, filepath.Join(t.TempDir(), "out.wav"))
		if err != nil || info.Duration < 3 || info.Duration > 30 || info.Peak > .93 {
			t.Fatal(info, err)
		}
	}
}
