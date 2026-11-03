package app

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

var sentences = regexp.MustCompile(`[^.!?\n]+[.!?]*|[.!?]+`)

func segmentText(segments []Segment) string {
	parts := make([]string, len(segments))
	for i, segment := range segments {
		parts[i] = segment.Text
	}
	return strings.Join(parts, "\n")
}

func speechContent(text string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, text)
}

func splitScript(s Script) []Segment {
	if s.InputMode != "text" && len(s.Segments) > 0 {
		return s.Segments
	}
	var result []Segment
	for _, part := range sentences.FindAllString(s.Text, -1) {
		part = strings.TrimSpace(part)
		for len(part) > 0 {
			n := len(part)
			if n > 500 {
				n = 500
				for n > 0 && !utf8.RuneStart(part[n]) {
					n--
				}
				if last := strings.LastIndex(part[:n], " "); last > 150 {
					n = last
				}
			}
			text := strings.TrimSpace(part[:n])
			part = strings.TrimSpace(part[n:])
			if text != "" {
				result = append(result, Segment{ID: apphost.ID(), Text: text, Voice: s.Voice})
			}
		}
	}
	return result
}
func transform(text string, rules []Rule) (string, error) {
	if len(text) > 5000 {
		return "", errors.New("Pronunciation expansion is too long")
	}
	for _, r := range rules {
		if r.From == "" {
			return "", errors.New("Dictionary replacement is invalid")
		}
		next := int64(len(text)) + int64(strings.Count(text, r.From))*int64(len(r.To)-len(r.From))
		if next > 5000 {
			return "", errors.New("Pronunciation expansion is too long")
		}
		text = strings.ReplaceAll(text, r.From, r.To)
	}
	return text, nil
}
func mergeRules(defaults, project []Rule) []Rule {
	out := append([]Rule(nil), defaults...)
	positions := map[string]int{}
	for i, r := range out {
		positions[r.From] = i
	}
	for _, r := range project {
		if i, ok := positions[r.From]; ok {
			out[i] = r
		} else {
			positions[r.From] = len(out)
			out = append(out, r)
		}
	}
	return out
}
func (a *App) enqueue(s Script, user string) (Job, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if e := validateScript(&s); e != nil {
		return Job{}, e
	}
	jobID := apphost.ID()
	if s.RequestID != "" {
		if !safeID.MatchString(s.RequestID) {
			return Job{}, errors.New("Invalid request identifier")
		}
		jobID = apphost.Hash(user + "\x00" + s.RequestID)
		if old, e := a.readJob(jobID, user); e == nil {
			return old, nil
		}
	}
	s.Segments = splitScript(s)
	if len(s.Segments) == 0 || len(s.Segments) > 1000 {
		return Job{}, errors.New("Add a script with at most 1,000 segments")
	}
	engine := a.Engines[s.Model]
	if engine == nil {
		return Job{}, errors.New("Choose an available model")
	}
	if !engine.Available() {
		return Job{}, errors.New("Install this model in Settings before generating")
	}
	var defaults string
	_ = a.H.DB.QueryRow(`SELECT body FROM preferences WHERE owner=? AND key='dictionary'`, user).Scan(&defaults)
	s.Defaults = nil
	_ = decodeJSON(defaults, &s.Defaults)
	var count int
	if a.H.DB.QueryRow(`SELECT count(*) FROM jobs WHERE status IN('queued','running')`).Scan(&count) != nil || count >= 20 {
		return Job{}, errors.New("The generation queue is full; wait for an item to finish")
	}
	seen := map[string]bool{}
	for i := range s.Segments {
		seg := &s.Segments[i]
		if seg.Voice == "" {
			seg.Voice = s.Voice
		}
		if len(seg.Text) == 0 || len(seg.Text) > 2000 {
			return Job{}, errors.New("Keep each segment between 1 and 2,000 characters")
		}
		if seg.ID == "" {
			seg.ID = apphost.ID()
		}
		if !safeID.MatchString(seg.ID) {
			return Job{}, errors.New("Invalid segment")
		}
		if seen[seg.ID] {
			return Job{}, errors.New("Duplicate segment identifier")
		}
		seen[seg.ID] = true
		seg.Output = ""
		v, e := a.findVoice(seg.Voice, user)
		if e != nil {
			return Job{}, errors.New("Voice unavailable; it may have been revoked")
		}
		if !engine.Compatible(v, s.Language) {
			return Job{}, errors.New("The model, language and selected voice are incompatible")
		}
	}
	// Only the server may reuse immutable audio from an owner's completed job.
	if s.Parent != "" {
		old, e := a.readJob(s.Parent, user)
		if e != nil || old.Status != "complete" {
			return Job{}, errors.New("Original recording unavailable")
		}
		s.Defaults = old.Script.Defaults
		if s.Model != old.Script.Model || s.Speed != old.Script.Speed || s.Pitch != old.Script.Pitch || s.Volume != old.Script.Volume || jsonText(s.Dictionary) != jsonText(old.Script.Dictionary) {
			return Job{}, errors.New("Keep model and audio controls unchanged when regenerating one segment")
		}
		if len(s.Segments) != len(old.Script.Segments) {
			return Job{}, errors.New("Segment count changed; generate a new recording")
		}
		for i, seg := range s.Segments {
			prior := old.Script.Segments[i]
			if seg.ID != prior.ID {
				return Job{}, errors.New("Segment order changed")
			}
			if seg.ID != s.Replace && (seg.Text != prior.Text || seg.Voice != prior.Voice) {
				return Job{}, errors.New("Only the selected segment may change")
			}
		}
	}
	j := Job{ID: jobID, Status: "queued", Created: now(), Script: s}
	_, e := a.H.DB.Exec(`INSERT INTO jobs(id,owner,status,created,body)VALUES(?,?,?,?,?)`, j.ID, user, j.Status, j.Created, jsonText(j))
	if e == nil {
		a.signal()
	}
	return j, e
}
func (a *App) createJob(w http.ResponseWriter, r *http.Request) {
	var s Script
	if !apphost.Decode(w, r, &s) {
		return
	}
	j, e := a.enqueue(s, uid(r))
	if e != nil {
		apphost.Problem(w, 400, e.Error())
		return
	}
	apphost.JSON(w, 202, j)
}
func (a *App) previewVoice(w http.ResponseWriter, r *http.Request) {
	v, e := a.findVoice(r.PathValue("id"), uid(r))
	if e != nil {
		apphost.Problem(w, 404, "Voice not found")
		return
	}
	model := "pocket18"
	if v.Engine == "moss" {
		model = "moss"
	}
	j, e := a.enqueue(Script{Name: "Voice preview · " + v.Name, Text: "Welcome to My Voice. Bring your words to life, one idea at a time.", Voice: v.ID, Model: model, Language: v.Language}, uid(r))
	if e != nil {
		apphost.Problem(w, 400, e.Error())
		return
	}
	apphost.JSON(w, 202, j)
}
func (a *App) processNext(parent context.Context) {
	a.mu.Lock()
	var id, user string
	e := a.H.DB.QueryRow(`SELECT id,owner FROM jobs WHERE status='queued' ORDER BY created LIMIT 1`).Scan(&id, &user)
	if e != nil {
		a.mu.Unlock()
		return
	}
	j, e := a.readJob(id, user)
	if e != nil {
		a.mu.Unlock()
		return
	}
	ctx, cancel := context.WithCancel(parent)
	a.cancel[id] = cancel
	j.Status = "running"
	a.storeJob(j)
	a.mu.Unlock()
	start := time.Now()
	e = a.generate(ctx, &j, user)
	cancel()
	a.mu.Lock()
	delete(a.cancel, id)
	if old, err := a.readJob(id, user); err == nil && old.Status == "cancelled" {
		j.Status = "cancelled"
	} else if parent.Err() != nil {
		j.Status = "queued"
	} else if e != nil {
		j.Status = "failed"
		j.Error = e.Error()
	} else {
		j.Status = "complete"
		j.Progress = 100
	}
	j.Wall = time.Since(start).Seconds()
	a.storeJob(j)
	a.mu.Unlock()
}
func (a *App) generate(ctx context.Context, j *Job, user string) error {
	a.inference.Lock()
	defer a.inference.Unlock()
	dir := a.work(j.ID)
	if e := os.MkdirAll(dir, 0700); e != nil {
		return e
	}
	var size int64
	_ = filepath.WalkDir(a.path("projects"), func(p string, d os.DirEntry, e error) error {
		if e == nil && !d.IsDir() {
			if i, e := d.Info(); e == nil {
				size += i.Size()
			}
		}
		return nil
	})
	if size > 10<<30 {
		return errors.New("My Voice storage reached its 10 GB limit; export and remove recordings before continuing")
	}
	engine := a.Engines[j.Script.Model]
	if engine == nil {
		return errors.New("Model not available")
	}
	// Per-user defaults are applied before project replacements.
	rules := mergeRules(j.Script.Defaults, j.Script.Dictionary)
	files := []string{}
	for i := range j.Script.Segments {
		seg := &j.Script.Segments[i]
		if e := ctx.Err(); e != nil {
			return errors.New("Generation stopped")
		}
		v, e := a.findVoice(seg.Voice, user)
		if e != nil {
			return errors.New("Voice permission was revoked; generation stopped")
		}
		target := filepath.Join(dir, seg.ID+".wav")
		reused := false
		if j.Script.Parent != "" && seg.ID != j.Script.Replace {
			reused = copyFile(filepath.Join(a.work(j.Script.Parent), seg.ID+".wav"), target) == nil
		}
		if !reused && seg.Output != "" {
			if _, e := inspectWAV(target); e == nil {
				reused = true
			}
		}
		if !reused {
			raw := filepath.Join(dir, seg.ID+"-raw.wav")
			text, e := transform(seg.Text, rules)
			if e != nil {
				return e
			}
			if e = engine.Generate(ctx, text, v, raw); e != nil {
				return e
			}
			if e = a.processAudio(ctx, raw, target, j.Script); e != nil {
				return e
			}
			_ = os.Remove(raw)
		}
		info, e := inspectWAV(target)
		if e != nil {
			return e
		}
		seg.Duration = info.Duration
		seg.Output = seg.ID
		files = append(files, target)
		j.Progress = (i + 1) * 95 / len(j.Script.Segments)
		a.mu.Lock()
		old, e := a.readJob(j.ID, user)
		if e == nil && old.Status == "cancelled" {
			a.mu.Unlock()
			return errors.New("Generation cancelled")
		}
		a.storeJob(*j)
		a.mu.Unlock()
	}
	if e := joinWAV(files, filepath.Join(dir, "output.wav"), j.Script.Pause); e != nil {
		return e
	}
	if e := a.convertMP3(ctx, filepath.Join(dir, "output.wav"), filepath.Join(dir, "output.mp3")); e != nil {
		return e
	}
	info, e := inspectWAV(filepath.Join(dir, "output.wav"))
	if e != nil {
		return e
	}
	j.Duration = info.Duration
	j.Waveform = info.Waveform
	return nil
}
func copyFile(src, dst string) error {
	in, e := os.Open(src)
	if e != nil {
		return e
	}
	defer in.Close()
	out, e := os.OpenFile(dst+".part", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if e != nil {
		return e
	}
	_, e = io.Copy(out, in)
	ce := out.Close()
	if e != nil {
		return e
	}
	if ce != nil {
		return ce
	}
	return os.Rename(dst+".part", dst)
}
func decodeJSON(b string, v any) error { return json.Unmarshal([]byte(b), v) }
