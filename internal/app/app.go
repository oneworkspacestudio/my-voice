package app

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

type App struct {
	H                *apphost.Host
	mu               sync.Mutex
	inference        sync.Mutex
	cancel           map[string]context.CancelFunc
	wake             chan struct{}
	Engines          map[string]SpeechEngine
	modelsMu         sync.Mutex
	installing       map[string]string
	runtimePath      string
	mossRuntimeReady bool
}
type Voice struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Engine    string `json:"engine"`
	Language  string `json:"language"`
	Tags      string `json:"tags"`
	Reference string `json:"-"`
	Owner     string `json:"owner,omitempty"`
	Consent   string `json:"consent,omitempty"`
	Preset    bool   `json:"preset"`
	Favorite  bool   `json:"favorite"`
	Shared    bool   `json:"shared"`
}
type Segment struct {
	ID       string  `json:"id"`
	Text     string  `json:"text"`
	Voice    string  `json:"voice"`
	Output   string  `json:"output,omitempty"`
	Duration float64 `json:"duration"`
}
type Script struct {
	Defaults   []Rule    `json:"accountDictionary,omitempty"`
	RequestID  string    `json:"requestId,omitempty"`
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Text       string    `json:"text"`
	InputMode  string    `json:"inputMode,omitempty"`
	Voice      string    `json:"voice"`
	Model      string    `json:"model"`
	Language   string    `json:"language"`
	Speed      float64   `json:"speed"`
	Pitch      float64   `json:"pitch"`
	Volume     float64   `json:"volume"`
	Pause      float64   `json:"pause"`
	Dictionary []Rule    `json:"dictionary"`
	Segments   []Segment `json:"segments"`
	Parent     string    `json:"parent,omitempty"`
	Replace    string    `json:"replace,omitempty"`
}
type Rule struct {
	From string `json:"from"`
	To   string `json:"to"`
}
type Job struct {
	ID       string    `json:"id"`
	Status   string    `json:"status"`
	Progress int       `json:"progress"`
	Error    string    `json:"error,omitempty"`
	Created  string    `json:"created"`
	Script   Script    `json:"script"`
	Duration float64   `json:"duration"`
	Wall     float64   `json:"wall"`
	Waveform []float64 `json:"waveform,omitempty"`
}

var safeID = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,64}$`)

var appMigrations = []string{`CREATE TABLE voices(id TEXT PRIMARY KEY,owner TEXT NOT NULL,body TEXT NOT NULL); CREATE TABLE grants(voice TEXT NOT NULL REFERENCES voices(id) ON DELETE CASCADE,user TEXT NOT NULL,PRIMARY KEY(voice,user)); CREATE TABLE projects(id TEXT PRIMARY KEY,owner TEXT NOT NULL,body TEXT NOT NULL); CREATE TABLE jobs(id TEXT PRIMARY KEY,owner TEXT NOT NULL,status TEXT NOT NULL,created TEXT NOT NULL,body TEXT NOT NULL); CREATE INDEX job_queue ON jobs(status,created); CREATE TABLE preferences(owner TEXT NOT NULL,key TEXT NOT NULL,body TEXT NOT NULL,PRIMARY KEY(owner,key));`, `ALTER TABLE projects ADD COLUMN updated TEXT NOT NULL DEFAULT '';`}

func New(h *apphost.Host) (*App, error) {
	a := &App{H: h, cancel: map[string]context.CancelFunc{}, wake: make(chan struct{}, 1), installing: map[string]string{}}
	var current int
	if e := h.DB.QueryRow(`SELECT coalesce(max(version),0) FROM host_migrations`).Scan(&current); e != nil {
		return nil, e
	}
	if current > len(appMigrations) {
		return nil, errors.New("App data needs a newer My Voice executable")
	}
	for i := current; i < len(appMigrations); i++ {
		if e := h.Migrate(i+1, appMigrations[i]); e != nil {
			return nil, e
		}
	}
	if _, e := h.DB.Exec(`UPDATE jobs SET status='queued' WHERE status='running'`); e != nil {
		return nil, e
	}
	a.Engines = map[string]SpeechEngine{"pocket18": &PocketEngine{A: a, ID: "pocket18"}, "pocket7": &PocketEngine{A: a, ID: "pocket7"}, "moss": &MossEngine{A: a}}
	routes := []struct {
		p     string
		f     http.HandlerFunc
		owner bool
	}{
		{"GET /api/models", a.models, false}, {"POST /api/models/{id}/install", a.installModel, true}, {"DELETE /api/models/{id}", a.removeModel, true},
		{"GET /api/voices", a.voices, false}, {"POST /api/voices", a.createVoice, false}, {"PATCH /api/voices/{id}", a.editVoice, false}, {"DELETE /api/voices/{id}", a.deleteVoice, false}, {"GET /api/voices/{id}/reference", a.reference, false}, {"POST /api/voices/{id}/preview", a.previewVoice, false}, {"POST /api/voices/{id}/share", a.shareVoice, false}, {"DELETE /api/voices/{id}/share/{user}", a.revokeVoice, false},
		{"POST /api/voices/{id}/favorite", a.favoriteVoice, false}, {"DELETE /api/history/{id}", a.deleteHistory, false},
		{"GET /api/voices/{id}/shares", a.voiceShares, false},
		{"GET /api/projects", a.projects, false}, {"POST /api/projects", a.saveProject, false}, {"GET /api/dictionaries", a.getPreference, false}, {"POST /api/dictionaries", a.setPreference, false}, {"GET /api/settings", a.settings, false}, {"PATCH /api/settings", a.setSettings, true},
		{"POST /api/tts/jobs", a.createJob, false}, {"GET /api/tts/jobs/{id}", a.getJob, false}, {"POST /api/tts/jobs/{id}/cancel", a.cancelJob, false}, {"POST /api/tts/jobs/{id}/retry", a.retryJob, false}, {"GET /api/history", a.history, false}, {"GET /api/outputs/{id}", a.output, false}, {"GET /api/outputs/{id}/download", a.output, false}, {"GET /api/outputs/{id}/segments/{segment}", a.segmentOutput, false}, {"GET /api/exports/zip", a.exportZIP, false},
	}
	for _, r := range routes {
		f := h.Require(r.f)
		if r.owner {
			f = h.Owner(r.f)
		}
		h.Mux.HandleFunc(r.p, f)
	}
	h.BackupPaths = append(h.BackupPaths, "models", "projects")
	return a, nil
}
func uid(r *http.Request) string { return apphost.Current(r).User.ID }
func jsonText(v any) string      { b, _ := json.Marshal(v); return string(b) }
func now() string                { return time.Now().UTC().Format(time.RFC3339Nano) }
func (a *App) path(parts ...string) string {
	return filepath.Join(append([]string{a.H.Config.Data}, parts...)...)
}
func (a *App) work(id string) string { return a.path("projects", id) }
func (a *App) signal() {
	select {
	case a.wake <- struct{}{}:
	default:
	}
}
func (a *App) Run(ctx context.Context) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			a.mu.Lock()
			for _, cancel := range a.cancel {
				cancel()
			}
			a.mu.Unlock()
			return nil
		case <-ticker.C:
		case <-a.wake:
		}
		if ctx.Err() != nil {
			return nil
		}
		a.processNext(ctx)
	}
}
func (a *App) presetVoices() []Voice {
	return []Voice{
		{ID: "alba", Name: "Alba", Engine: "pocket", Language: "en", Tags: "Warm · Conversational · CC BY 4.0", Preset: true},
		{ID: "moss-open", Name: "Open Voice", Engine: "moss", Language: "en", Tags: "CC0 donated reference · MOSS", Preset: true},
	}
}
func (a *App) findVoice(id, user string) (Voice, error) {
	for _, v := range a.presetVoices() {
		if v.ID == id {
			return v, nil
		}
	}
	var b string
	e := a.H.DB.QueryRow(`SELECT body FROM voices WHERE id=? AND (owner=? OR id IN(SELECT voice FROM grants WHERE user=?))`, id, user, user).Scan(&b)
	var v Voice
	if e == nil {
		e = json.Unmarshal([]byte(b), &v)
	}
	return v, e
}
func (a *App) voices(w http.ResponseWriter, r *http.Request) {
	out := a.presetVoices()
	for i := range out {
		var body string
		_ = a.H.DB.QueryRow(`SELECT body FROM preferences WHERE owner=? AND key=?`, uid(r), "favorite-"+out[i].ID).Scan(&body)
		out[i].Favorite = body == "true"
	}
	rows, e := a.H.DB.Query(`SELECT body FROM voices WHERE owner=? OR id IN(SELECT voice FROM grants WHERE user=?)`, uid(r), uid(r))
	if e != nil {
		a.H.DBError(w, e)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var b string
		var v Voice
		if rows.Scan(&b) != nil {
			continue
		}
		if json.Unmarshal([]byte(b), &v) == nil {
			v.Shared = v.Owner != uid(r)
			out = append(out, v)
		}
	}
	apphost.JSON(w, 200, out)
}
func (a *App) editVoice(w http.ResponseWriter, r *http.Request) {
	var p struct {
		Name     string `json:"name"`
		Tags     string `json:"tags"`
		Favorite bool   `json:"favorite"`
	}
	if !apphost.Decode(w, r, &p) {
		return
	}
	var b string
	if a.H.DB.QueryRow(`SELECT body FROM voices WHERE id=? AND owner=?`, r.PathValue("id"), uid(r)).Scan(&b) != nil {
		apphost.Problem(w, 404, "Voice not found")
		return
	}
	var v Voice
	_ = json.Unmarshal([]byte(b), &v)
	if len(p.Name) < 1 || len(p.Name) > 100 || len(p.Tags) > 300 {
		apphost.Problem(w, 400, "Choose a name under 100 characters")
		return
	}
	v.Name = p.Name
	v.Tags = p.Tags
	v.Favorite = p.Favorite
	_, e := a.H.DB.Exec(`UPDATE voices SET body=? WHERE id=? AND owner=?`, jsonText(v), v.ID, uid(r))
	if e != nil {
		a.H.DBError(w, e)
		return
	}
	apphost.JSON(w, 200, v)
}
func (a *App) deleteVoice(w http.ResponseWriter, r *http.Request) {
	a.mu.Lock()
	defer a.mu.Unlock()
	var b string
	if a.H.DB.QueryRow(`SELECT body FROM voices WHERE id=? AND owner=?`, r.PathValue("id"), uid(r)).Scan(&b) != nil {
		apphost.Problem(w, 404, "Voice not found")
		return
	}
	var v Voice
	_ = json.Unmarshal([]byte(b), &v)
	if _, e := a.H.DB.Exec(`DELETE FROM voices WHERE id=? AND owner=?`, v.ID, uid(r)); e != nil {
		a.H.DBError(w, e)
		return
	}
	// Reference and model-specific state are app-generated names. No user paths.
	if safeID.MatchString(v.ID) {
		_ = os.RemoveAll(a.path("projects", "voice-"+v.ID))
	}
	a.H.Audit(uid(r), "voice-delete", v.ID)
	apphost.JSON(w, 200, map[string]bool{"deleted": true})
}
func (a *App) shareVoice(w http.ResponseWriter, r *http.Request) {
	var p struct {
		User string `json:"user"`
	}
	if !apphost.Decode(w, r, &p) {
		return
	}
	var n int
	if a.H.DB.QueryRow(`SELECT count(*) FROM voices WHERE id=? AND owner=?`, r.PathValue("id"), uid(r)).Scan(&n) != nil || n != 1 {
		apphost.Problem(w, 404, "Voice not found")
		return
	}
	var user string
	if a.H.DB.QueryRow(`SELECT id FROM host_users WHERE email=? AND role!='disabled'`, p.User).Scan(&user) != nil {
		apphost.Problem(w, 400, "Ask this person to open My Voice with their OneAccount first")
		return
	}
	_, e := a.H.DB.Exec(`INSERT OR IGNORE INTO grants(voice,user)VALUES(?,?)`, r.PathValue("id"), user)
	if e != nil {
		a.H.DBError(w, e)
		return
	}
	a.H.Audit(uid(r), "voice-share", r.PathValue("id"))
	apphost.JSON(w, 200, map[string]string{"user": user})
}
func (a *App) revokeVoice(w http.ResponseWriter, r *http.Request) {
	result, e := a.H.DB.Exec(`DELETE FROM grants WHERE voice=? AND user=? AND voice IN(SELECT id FROM voices WHERE owner=?)`, r.PathValue("id"), r.PathValue("user"), uid(r))
	if e != nil {
		a.H.DBError(w, e)
		return
	}
	n, _ := result.RowsAffected()
	a.H.Audit(uid(r), "voice-revoke", r.PathValue("id"))
	apphost.JSON(w, 200, map[string]bool{"revoked": n > 0})
}
func (a *App) reference(w http.ResponseWriter, r *http.Request) {
	v, e := a.findVoice(r.PathValue("id"), uid(r))
	if e != nil || v.Preset || v.Owner != uid(r) {
		apphost.Problem(w, 404, "Reference not found")
		return
	}
	http.ServeFile(w, r, a.path("projects", "voice-"+v.ID, "reference.wav"))
}
func (a *App) projects(w http.ResponseWriter, r *http.Request) {
	rows, e := a.H.DB.Query(`SELECT body FROM projects WHERE owner=? ORDER BY updated DESC,rowid DESC LIMIT 100`, uid(r))
	if e != nil {
		a.H.DBError(w, e)
		return
	}
	defer rows.Close()
	out := []json.RawMessage{}
	for rows.Next() {
		var b string
		if rows.Scan(&b) == nil {
			out = append(out, json.RawMessage(b))
		}
	}
	apphost.JSON(w, 200, out)
}
func validateScript(s *Script) error {
	if len(s.Segments) > 1000 {
		return errors.New("Script exceeds project limits")
	}
	switch s.InputMode {
	case "text":
		// Cached/recorded chunks must never override the current editor text.
		s.Segments = nil
	case "segments":
		s.Text = segmentText(s.Segments)
	case "":
		if len(s.Segments) == 0 {
			s.InputMode = "text"
		} else {
			if s.Parent == "" && strings.TrimSpace(s.Text) != "" && speechContent(s.Text) != speechContent(segmentText(s.Segments)) {
				return errors.New("The script and saved segments differ. Refresh My Voice and generate the full script again")
			}
			s.InputMode = "segments"
			s.Text = segmentText(s.Segments)
		}
	default:
		return errors.New("Choose full-script or edited-segment generation")
	}
	if s.Name == "" {
		s.Name = "Untitled speech"
	}
	if len(s.Name) > 200 || len(s.Text) > 200000 || len(s.Segments) > 1000 || len(s.Dictionary) > 100 {
		return errors.New("Script exceeds project limits")
	}
	if s.Voice == "" {
		s.Voice = "alba"
	}
	if s.Model == "" || s.Model == "auto" {
		s.Model = "pocket18"
	}
	if s.Language == "" {
		s.Language = "en"
	}
	if s.Speed == 0 {
		s.Speed = 1
	}
	if s.Volume == 0 {
		s.Volume = 1
	}
	if s.Speed < 0.5 || s.Speed > 2 || s.Pitch < -12 || s.Pitch > 12 || s.Volume < 0.1 || s.Volume > 2 || s.Pause < 0 || s.Pause > 5 {
		return errors.New("Audio controls are out of range")
	}
	for _, rule := range s.Dictionary {
		if rule.From == "" || len(rule.From) > 200 || len(rule.To) > 500 {
			return errors.New("Dictionary replacement is invalid")
		}
	}
	return nil
}
func (a *App) saveProject(w http.ResponseWriter, r *http.Request) {
	var s Script
	if !apphost.Decode(w, r, &s) {
		return
	}
	if e := validateScript(&s); e != nil {
		apphost.Problem(w, 400, e.Error())
		return
	}
	if s.ID == "" {
		s.ID = apphost.ID()
	}
	if !safeID.MatchString(s.ID) {
		apphost.Problem(w, 400, "Invalid project")
		return
	}
	result, e := a.H.DB.Exec(`INSERT INTO projects(id,owner,body,updated)VALUES(?,?,?,?) ON CONFLICT(id) DO UPDATE SET body=excluded.body,updated=excluded.updated WHERE projects.owner=excluded.owner`, s.ID, uid(r), jsonText(s), now())
	if e != nil {
		a.H.DBError(w, e)
		return
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		apphost.Problem(w, 404, "Project not found")
		return
	}
	apphost.JSON(w, 200, s)
}
func (a *App) getPreference(w http.ResponseWriter, r *http.Request) {
	var b string
	if a.H.DB.QueryRow(`SELECT body FROM preferences WHERE owner=? AND key='dictionary'`, uid(r)).Scan(&b) != nil {
		b = "[]"
	}
	apphost.JSON(w, 200, json.RawMessage(b))
}
func (a *App) setPreference(w http.ResponseWriter, r *http.Request) {
	var rules []Rule
	if !apphost.Decode(w, r, &rules) {
		return
	}
	s := Script{Dictionary: rules}
	if e := validateScript(&s); e != nil {
		apphost.Problem(w, 400, e.Error())
		return
	}
	_, e := a.H.DB.Exec(`INSERT INTO preferences(owner,key,body)VALUES(?,'dictionary',?) ON CONFLICT(owner,key) DO UPDATE SET body=excluded.body`, uid(r), jsonText(rules))
	if e != nil {
		a.H.DBError(w, e)
		return
	}
	apphost.JSON(w, 200, rules)
}
func (a *App) history(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if len(query) > 200 {
		apphost.Problem(w, 400, "Keep the search under 200 characters")
		return
	}
	rows, e := a.H.DB.Query(`SELECT body,status FROM jobs WHERE owner=? AND instr(lower(body),lower(?))>0 ORDER BY created DESC LIMIT 200`, uid(r), query)
	if e != nil {
		a.H.DBError(w, e)
		return
	}
	defer rows.Close()
	out := []Job{}
	q := strings.ToLower(r.URL.Query().Get("q"))
	for rows.Next() {
		var b, status string
		var j Job
		if rows.Scan(&b, &status) == nil && json.Unmarshal([]byte(b), &j) == nil {
			j.Status = status
			if q == "" || strings.Contains(strings.ToLower(j.Script.Name+" "+j.Script.Text), q) {
				out = append(out, j)
			}
		}
	}
	apphost.JSON(w, 200, out)
}
func (a *App) readJob(id, user string) (Job, error) {
	var j Job
	var b, status string
	e := a.H.DB.QueryRow(`SELECT body,status FROM jobs WHERE id=? AND owner=?`, id, user).Scan(&b, &status)
	if e == nil {
		e = json.Unmarshal([]byte(b), &j)
		j.Status = status
	}
	return j, e
}
func (a *App) getJob(w http.ResponseWriter, r *http.Request) {
	j, e := a.readJob(r.PathValue("id"), uid(r))
	if e != nil {
		apphost.Problem(w, 404, "Generation not found")
		return
	}
	apphost.JSON(w, 200, j)
}
func (a *App) storeJob(j Job) {
	_, e := a.H.DB.Exec(`UPDATE jobs SET body=?,status=? WHERE id=?`, jsonText(j), j.Status, j.ID)
	if e != nil {
		a.H.Log.Error("generation state save failed")
	}
}
func (a *App) cancelJob(w http.ResponseWriter, r *http.Request) {
	a.mu.Lock()
	defer a.mu.Unlock()
	j, e := a.readJob(r.PathValue("id"), uid(r))
	if e != nil {
		apphost.Problem(w, 404, "Generation not found")
		return
	}
	if j.Status == "running" || j.Status == "queued" {
		j.Status = "cancelled"
		if f := a.cancel[j.ID]; f != nil {
			f()
		}
		a.storeJob(j)
	}
	apphost.JSON(w, 200, j)
}
func (a *App) retryJob(w http.ResponseWriter, r *http.Request) {
	a.mu.Lock()
	defer a.mu.Unlock()
	j, e := a.readJob(r.PathValue("id"), uid(r))
	if e != nil {
		apphost.Problem(w, 404, "Generation not found")
		return
	}
	if j.Status != "failed" && j.Status != "cancelled" {
		apphost.Problem(w, 409, "Only stopped generations can be retried")
		return
	}
	j.Status = "queued"
	j.Error = ""
	a.storeJob(j)
	a.signal()
	apphost.JSON(w, 202, j)
}
func (a *App) output(w http.ResponseWriter, r *http.Request) {
	j, e := a.readJob(r.PathValue("id"), uid(r))
	if e != nil || j.Status != "complete" {
		apphost.Problem(w, 404, "Audio is not ready")
		return
	}
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "wav"
	}
	if format != "wav" && format != "mp3" {
		apphost.Problem(w, 400, "Choose WAV or MP3")
		return
	}
	w.Header().Set("Content-Type", map[string]string{"wav": "audio/wav", "mp3": "audio/mpeg"}[format])
	if strings.HasSuffix(r.URL.Path, "/download") {
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="my-voice-%s.%s"`, j.ID[:12], format))
	}
	http.ServeFile(w, r, a.work(j.ID)+"/output."+format)
}
func (a *App) segmentOutput(w http.ResponseWriter, r *http.Request) {
	j, e := a.readJob(r.PathValue("id"), uid(r))
	s := r.PathValue("segment")
	if e != nil || !safeID.MatchString(s) {
		apphost.Problem(w, 404, "Segment not found")
		return
	}
	found := false
	for _, v := range j.Script.Segments {
		if v.ID == s && v.Output != "" {
			found = true
		}
	}
	if !found {
		apphost.Problem(w, 404, "Segment not ready")
		return
	}
	http.ServeFile(w, r, filepath.Join(a.work(j.ID), s+".wav"))
}
func (a *App) exportZIP(w http.ResponseWriter, r *http.Request) {
	ids := strings.Split(r.URL.Query().Get("ids"), ",")
	if len(ids) > 20 {
		apphost.Problem(w, 400, "Export up to 20 recordings at once")
		return
	}
	for _, id := range ids {
		j, e := a.readJob(id, uid(r))
		if e != nil || j.Status != "complete" {
			apphost.Problem(w, 404, "A recording is not ready")
			return
		}
	}
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", `attachment; filename="my-voice-recordings.zip"`)
	z := zip.NewWriter(w)
	defer z.Close()
	for _, id := range ids {
		f, e := os.Open(filepath.Join(a.work(id), "output.wav"))
		if e != nil {
			return
		}
		dst, e := z.Create(id + ".wav")
		if e == nil {
			_, _ = io.Copy(dst, f)
		}
		f.Close()
	}
}
