package app

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Qwen uses the pinned native CPU runtime. Reference audio remains private and
// model-specific; MOSS state is never sent to a different engine.
type QwenEngine struct {
	A  *App
	ID string
}

func (q *QwenEngine) definition() Model {
	for _, m := range modelList() {
		if m.ID == q.ID {
			return m
		}
	}
	return Model{}
}
func (q *QwenEngine) Available() bool {
	m := q.definition()
	b, e := os.ReadFile(filepath.Join(q.A.modelDir(q.ID), "installed.json"))
	var marker struct {
		Version string `json:"version"`
	}
	if e != nil || json.Unmarshal(b, &marker) != nil || marker.Version != m.Version {
		return false
	}
	for _, a := range m.Assets {
		f, e := os.Stat(filepath.Join(q.A.modelDir(q.ID), filepath.FromSlash(a.Path)))
		if e != nil || f.Size() != a.Size {
			return false
		}
	}
	return m.ID != ""
}
func (q *QwenEngine) Cloning() bool { return q.Available() }
func (q *QwenEngine) Compatible(v Voice, lang string) bool {
	if v.Engine != q.ID {
		return false
	}
	for _, supported := range q.definition().Languages {
		if supported == lang {
			return true
		}
	}
	return false
}
func (q *QwenEngine) Generate(ctx context.Context, text string, v Voice, lang, out string) error {
	dir, e := q.A.runtime()
	if e != nil {
		return e
	}
	reference := filepath.Join(q.A.modelDir(q.ID), "presets", "open-voice.wav")
	if !v.Preset {
		reference = q.A.path("projects", "voice-"+v.ID, "reference.wav")
	}
	if _, e = os.Stat(reference); e != nil {
		return errors.New("This voice reference is missing. Create the voice again with this model")
	}
	// A UTF-8 file keeps text out of process arguments and handles non-Latin text
	// without depending on the Windows command-line code page.
	prompt := out + ".txt"
	if e = os.WriteFile(prompt, []byte(text), 0600); e != nil {
		return e
	}
	defer os.Remove(prompt)
	args := []string{"-m", filepath.Join(q.A.modelDir(q.ID), "model.gguf"), "-mm", filepath.Join(q.A.modelDir(q.ID), "codec.gguf"),
		"-f", prompt, "--no-escape", "--tts-lang", lang, "--tts-speaker-file", reference, "--output", out,
		"-t", q.A.threads(), "-tb", q.A.threads(), "-ngl", "0", "--no-mmproj-offload", "-n", "768", "-c", "2048"}
	return q.A.service(ctx, filepath.Join(dir, "qwen", "llama-tts.exe"), args, nil)
}
func (q *QwenEngine) Clone(_ context.Context, src, dst string) error {
	if _, e := inspectWAV(src); e != nil {
		return errors.New("This reference could not be read. Choose clear speech and retry")
	}
	return os.WriteFile(dst+".json", []byte(jsonText(map[string]string{"model": q.ID, "reference": "reference.wav"})), 0600)
}

func (q *QwenEngine) prepare(ctx context.Context) error {
	dir := q.A.modelDir(q.ID)
	raw, converted := filepath.Join(dir, "setup-raw.wav"), filepath.Join(dir, "setup-check.wav")
	defer os.Remove(raw)
	defer os.Remove(converted)
	if e := q.Generate(ctx, "This model is ready.", Voice{Preset: true, Engine: q.ID}, "en", raw); e != nil {
		return e
	}
	if e := q.A.ffmpeg(ctx, "-i", raw, "-ac", "1", "-ar", "24000", "-c:a", "pcm_s16le", converted); e != nil {
		return errors.New("The model ran but its audio could not be checked. Retry model setup")
	}
	info, e := inspectWAV(converted)
	if e != nil || info.Duration < .1 || info.RMS < .0001 {
		return errors.New("The model did not produce usable audio. Retry model setup before recording")
	}
	return nil
}
