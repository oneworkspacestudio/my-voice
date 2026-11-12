package app

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// These are build inputs locked by scripts/prepare-runtime.py, never downloaded
// during ordinary startup. The public source archive contains the recipe/locks.
//
//go:embed payload.zip models.json pocket18.yaml pocket7.yaml
var resources embed.FS
var extractMu sync.Mutex

type SpeechEngine interface {
	Available() bool
	Cloning() bool
	Compatible(Voice, string) bool
	Generate(context.Context, string, Voice, string, string) error
	Clone(context.Context, string, string) error
}
type Asset struct {
	URL   string `json:"url"`
	Path  string `json:"path"`
	Size  int64  `json:"size"`
	SHA   string `json:"sha256"`
	Gated bool   `json:"gated,omitempty"`
}
type Model struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Version      string   `json:"version"`
	Languages    []string `json:"languages"`
	License      string   `json:"license"`
	URL          string   `json:"url"`
	Runtime      string   `json:"runtime"`
	Cloning      bool     `json:"cloning"`
	CloningReady bool     `json:"cloningReady"`
	Installing   bool     `json:"installing"`
	Installed    bool     `json:"installed"`
	Status       string   `json:"status"`
	Size         int64    `json:"size"`
	Summary      string   `json:"summary,omitempty"`
	ShortName    string   `json:"shortName,omitempty"`
	RAMGB        int      `json:"ramGB,omitempty"`
	Assets       []Asset  `json:"assets,omitempty"`
}

func modelList() []Model {
	b, _ := resources.ReadFile("models.json")
	var out []Model
	_ = json.Unmarshal(b, &out)
	return out
}
func (a *App) models(w http.ResponseWriter, r *http.Request) {
	out := modelList()
	a.modelsMu.Lock()
	defer a.modelsMu.Unlock()
	for i := range out {
		out[i].Installed = a.Engines[out[i].ID].Available()
		// Capability must remain discoverable before installation. Readiness is
		// separate: MOSS prepares a real reference state before marking installed.
		out[i].CloningReady = out[i].Installed && a.Engines[out[i].ID].Cloning()
		out[i].Cloning = out[i].Cloning || out[i].CloningReady
		out[i].Status = a.installing[out[i].ID]
		out[i].Installing = modelInstalling(out[i].Status)
		out[i].Assets = nil
	}
	apphost.JSON(w, 200, out)
}
func modelInstalling(status string) bool {
	return strings.HasPrefix(status, "Downloading") || strings.HasPrefix(status, "Preparing")
}
func (a *App) modelDir(id string) string { return a.path("models", id) }
func (a *App) settings(w http.ResponseWriter, r *http.Request) {
	var b string
	_ = a.H.DB.QueryRow(`SELECT body FROM preferences WHERE owner='' AND key='settings'`).Scan(&b)
	threads := 2
	var cfg struct {
		Threads int `json:"threads"`
	}
	if json.Unmarshal([]byte(b), &cfg) == nil && cfg.Threads > 0 {
		threads = cfg.Threads
	}
	var usage int64
	_ = filepath.WalkDir(a.H.Config.Data, func(_ string, d fs.DirEntry, err error) error {
		if err == nil && d.Type().IsRegular() {
			if info, e := d.Info(); e == nil {
				usage += info.Size()
			}
		}
		return nil
	})
	available, _ := freeBytes(a.H.Config.Data)
	var podcastsConnected bool
	if apphost.Current(r).User.Role == "owner" {
		_ = a.H.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM preferences WHERE owner='' AND key=?)", podcastGrantKey).Scan(&podcastsConnected)
	}
	apphost.JSON(w, 200, map[string]any{"threads": threads, "storage": a.H.Config.Data, "usedBytes": usage, "availableBytes": available, "retention": "Recordings and reusable segments stay until you remove them in History. Raw successful generation intermediates are removed.", "network": "Authenticated Workspace proxy; loopback listener", "drive": "Export to My Drive using the signed-in Drive account", "version": a.H.Config.Version, "queueLimit": 20, "storageLimitGB": 10, "podcastsConnected": podcastsConnected})
}
func (a *App) setSettings(w http.ResponseWriter, r *http.Request) {
	var cfg struct {
		Threads int `json:"threads"`
	}
	if !apphost.Decode(w, r, &cfg) {
		return
	}
	if cfg.Threads < 1 || cfg.Threads > 8 {
		apphost.Problem(w, 400, "Choose 1–8 CPU threads")
		return
	}
	_, e := a.H.DB.Exec(`INSERT INTO preferences(owner,key,body)VALUES('','settings',?) ON CONFLICT(owner,key) DO UPDATE SET body=excluded.body`, jsonText(cfg))
	if e != nil {
		a.H.DBError(w, e)
		return
	}
	apphost.JSON(w, 200, cfg)
}
func (a *App) threads() string {
	var b string
	var cfg struct {
		Threads int `json:"threads"`
	}
	_ = a.H.DB.QueryRow(`SELECT body FROM preferences WHERE owner='' AND key='settings'`).Scan(&b)
	_ = json.Unmarshal([]byte(b), &cfg)
	if cfg.Threads < 1 || cfg.Threads > 8 {
		return "2"
	}
	return strconv.Itoa(cfg.Threads)
}
func (a *App) service(ctx context.Context, exe string, args, extra []string) error {
	env := append([]string{"OMP_NUM_THREADS=" + a.threads(), "MKL_NUM_THREADS=" + a.threads(), "RAYON_NUM_THREADS=" + a.threads(), "HF_HUB_OFFLINE=1", "HF_HOME=" + a.path("models", "hf"), "USERPROFILE=" + a.H.Config.Data}, extra...)
	if e := os.MkdirAll(a.path("tmp"), 0700); e != nil {
		return e
	}
	// Use apphost's suspended-start, kill-on-close Job Object; never inherit SSO secrets.
	var failure *workerFailure
	for attempt := 0; attempt < 3; attempt++ {
		e := apphost.Service(ctx, exe, a.path("tmp"), args, env, io.Discard)
		if e == nil {
			return nil
		}
		if ctx.Err() != nil {
			return errors.New("Generation stopped")
		}
		// The error contains only the executable/OS status. Never record argv,
		// standard output, speech text, references or inherited credentials.
		a.H.Log.Warn("worker failed", "worker", filepath.Base(exe), "attempt", attempt+1, "cause", e.Error())
		failure = describeWorkerFailure(exe, e)
		if failure.permanent() {
			return failure
		}
		select {
		case <-ctx.Done():
			return errors.New("Generation stopped")
		case <-time.After(time.Duration(attempt+1) * 200 * time.Millisecond):
		}
	}
	return failure
}

// Only fixed messages, the bundled executable's name and numeric OS status may
// reach a job response. Worker output and command lines can contain private text.
type workerFailure struct {
	worker string
	code   uint32
	kind   string
}

func describeWorkerFailure(exe string, err error) *workerFailure {
	f := &workerFailure{worker: filepath.Base(exe), kind: "start"}
	var exited *exec.ExitError
	switch {
	case errors.As(err, &exited):
		f.code, f.kind = uint32(exited.ExitCode()), "exit"
	case os.IsPermission(err):
		f.kind = "permission"
	case os.IsNotExist(err):
		f.kind = "missing"
	}
	return f
}
func (f *workerFailure) permanent() bool {
	return f.kind == "permission" || f.kind == "missing" || f.code == 0xc000001d || f.code == 0xc0000135 || f.code == 0xc000007b
}
func (f *workerFailure) Error() string {
	switch f.kind {
	case "permission":
		return "Windows blocked the speech worker (" + f.worker + "). Check its security protection or app-folder permissions, then retry"
	case "missing":
		return "The speech worker is missing (" + f.worker + "). Repair My Voice through Workspace Updates, then retry"
	case "exit":
		reason := "Check the installed model and audio, then retry"
		switch f.code {
		case 0xc000001d:
			reason = "This worker requires CPU instructions unavailable on this PC"
		case 0xc0000135, 0xc000007b:
			reason = "A Windows runtime dependency is missing or incompatible. Repair My Voice through Workspace Updates"
		case 0xc0000017, 0xc000009a:
			reason = "Windows ran out of memory. Close other apps or select the smaller Pocket model, then retry"
		}
		return fmt.Sprintf("Speech worker %s exited (0x%08X). %s", f.worker, f.code, reason)
	default:
		return "Windows could not start the speech worker (" + f.worker + "). Check the app log for the OS error, then retry"
	}
}

func (a *App) pocketService(ctx context.Context, dir string, args []string) error {
	err := a.service(ctx, filepath.Join(dir, "pocket-tts.exe"), args, nil)
	var failed *workerFailure
	if errors.As(err, &failed) && failed.code == 0xc000001d {
		// The upstream Windows release targets x86-64-v3. Fall back only for
		// STATUS_ILLEGAL_INSTRUCTION, retaining its fast kernels on newer PCs.
		a.H.Log.Info("using portable Pocket worker", "reason", "unsupported CPU instructions")
		return a.service(ctx, filepath.Join(dir, "pocket-tts-portable.exe"), args, nil)
	}
	return err
}
func (a *App) runtime() (string, error) {
	return a.runtimeFor(false)
}
func (a *App) runtimeFor(moss bool) (string, error) {
	extractMu.Lock()
	defer extractMu.Unlock()
	if a.runtimePath != "" && (!moss || a.mossRuntimeReady) {
		return a.runtimePath, nil
	}
	b, e := resources.ReadFile("payload.zip")
	if e != nil {
		return "", e
	}
	sum := sha256.Sum256(b)
	dir := a.path("runtime", hex.EncodeToString(sum[:8]))
	z, e := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if e != nil {
		return "", e
	}
	for _, m := range z.File {
		if !moss && strings.HasPrefix(m.Name, "moss/") {
			continue
		}
		rel := filepath.FromSlash(m.Name)
		if filepath.IsAbs(rel) || strings.Contains(rel, "..") || strings.Contains(rel, ":") || m.Mode()&os.ModeSymlink != 0 {
			return "", errors.New("Invalid bundled runtime")
		}
		dst := filepath.Join(dir, rel)
		if m.FileInfo().IsDir() {
			continue
		}
		src, e := m.Open()
		if e != nil {
			return "", e
		}
		raw, e := io.ReadAll(io.LimitReader(src, 256<<20))
		src.Close()
		if e != nil || uint64(len(raw)) != m.UncompressedSize64 {
			return "", errors.New("Invalid bundled file")
		}
		wanted := sha256.Sum256(raw)
		actual, e := os.ReadFile(dst)
		if e == nil && sha256.Sum256(actual) == wanted {
			continue
		}
		if e = os.MkdirAll(filepath.Dir(dst), 0700); e != nil {
			return "", e
		}
		if e = os.WriteFile(dst+".part", raw, 0600); e != nil {
			return "", e
		}
		if e = os.Rename(dst+".part", dst); e != nil {
			return "", e
		}
	}
	a.runtimePath = dir
	if moss {
		a.mossRuntimeReady = true
	}
	return dir, nil
}
func (a *App) installModel(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var selected *Model
	for _, m := range modelList() {
		if m.ID == id {
			copy := m
			selected = &copy
		}
	}
	if selected == nil {
		apphost.Problem(w, 404, "Unknown model")
		return
	}
	var p struct {
		License bool   `json:"license"`
		Token   string `json:"token"`
		Cloning bool   `json:"cloning"`
	}
	if !apphost.Decode(w, r, &p) {
		return
	}
	if !p.License {
		apphost.Problem(w, 400, "Review and accept the linked model licence first")
		return
	}
	if p.Cloning || p.Token != "" {
		apphost.Problem(w, 400, "These models use pinned public downloads. No access token or extra cloning access is needed")
		return
	}
	if len(p.Token) > 512 || strings.ContainsAny(p.Token, "\r\n") {
		apphost.Problem(w, 400, "Invalid access token")
		return
	}
	a.modelsMu.Lock()
	if modelInstalling(a.installing[id]) {
		a.modelsMu.Unlock()
		apphost.Problem(w, 409, "Download already running")
		return
	}
	a.installing[id] = "Downloading 0%"
	a.modelsMu.Unlock()
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 90*time.Minute)
		defer cancel()
		e := a.downloadModel(ctx, *selected, p.Token, p.Cloning)
		a.modelsMu.Lock()
		defer a.modelsMu.Unlock()
		if e != nil {
			a.installing[id] = e.Error()
		} else {
			a.installing[id] = "Ready"
		}
	}()
	a.H.Audit(uid(r), "model-install", id)
	apphost.JSON(w, 202, map[string]string{"status": "downloading"})
}
func (a *App) downloadModel(ctx context.Context, m Model, token string, cloning bool) error {
	dir := a.modelDir(m.ID)
	if e := os.MkdirAll(dir, 0700); e != nil {
		return e
	}
	client := &http.Client{Timeout: 45 * time.Minute}
	var total, done int64
	for _, v := range m.Assets {
		if !v.Gated || cloning {
			total += v.Size
		}
	}
	for _, v := range m.Assets {
		if v.Gated && !cloning {
			continue
		}
		dest := filepath.Join(dir, filepath.FromSlash(v.Path))
		if e := os.MkdirAll(filepath.Dir(dest), 0700); e != nil {
			return e
		}
		if verifiedModelFile(dest, v) {
			done += v.Size
			continue
		}
		available, e := freeBytes(dir)
		if e != nil {
			return errors.New("Cannot check model storage space")
		}
		remaining := v.Size
		if part, e := os.Stat(dest + ".part"); e == nil && part.Size() < v.Size {
			remaining -= part.Size()
		}
		if available < uint64(remaining)+(128<<20) {
			return errors.New("Not enough space for this model download; free disk space and retry to continue")
		}
		show := func(n int64) {
			a.modelsMu.Lock()
			a.installing[m.ID] = fmt.Sprintf("Downloading %d%%", (done+n)*100/total)
			a.modelsMu.Unlock()
		}
		if e := downloadModelAsset(ctx, client, v, dest, show); e != nil {
			return e
		}
		n := v.Size
		done += n
		a.modelsMu.Lock()
		a.installing[m.ID] = fmt.Sprintf("Downloading %d%%", done*100/total)
		a.modelsMu.Unlock()
	}
	if strings.HasPrefix(m.ID, "pocket") {
		name := "pocket18.yaml"
		if m.ID == "pocket7" {
			name = "pocket7.yaml"
		}
		raw, _ := resources.ReadFile(name)
		lines := strings.Split(string(raw), "\n")
		for i, line := range lines {
			trim := strings.TrimSpace(line)
			indent := line[:len(line)-len(strings.TrimLeft(line, " "))]
			if strings.HasPrefix(trim, "weights_path:") {
				weight := "model-tts.safetensors"
				if cloning {
					weight = "model-clone.safetensors"
				}
				lines[i] = indent + "weights_path: " + strconv.Quote(filepath.ToSlash(filepath.Join(dir, weight)))
			}
			if strings.HasPrefix(trim, "weights_path_without_voice_cloning:") {
				lines[i] = ""
			}
			if strings.HasPrefix(trim, "tokenizer_path:") {
				lines[i] = indent + "tokenizer_path: " + strconv.Quote(filepath.ToSlash(filepath.Join(dir, "tokenizer.model")))
			}
		}
		if e := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(strings.Join(lines, "\n")), 0600); e != nil {
			return e
		}
	}
	if m.ID == "moss" {
		a.modelsMu.Lock()
		a.installing[m.ID] = "Preparing cloning model…"
		a.modelsMu.Unlock()
		a.inference.Lock()
		e := a.Engines["moss"].Clone(ctx, filepath.Join(dir, "presets", "open-voice.wav"), filepath.Join(dir, "preset-state"))
		a.inference.Unlock()
		if e != nil {
			return e
		}
	}
	if q, ok := a.Engines[m.ID].(*QwenEngine); ok {
		a.modelsMu.Lock()
		a.installing[m.ID] = "Preparing cloning model and checking audio…"
		a.modelsMu.Unlock()
		a.inference.Lock()
		e := q.prepare(ctx)
		a.inference.Unlock()
		if e != nil {
			return e
		}
	}
	return os.WriteFile(filepath.Join(dir, "installed.json"), []byte(jsonText(map[string]any{"version": m.Version, "cloning": cloning || m.Cloning, "verified": now()})), 0600)
}
func (a *App) removeModel(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if a.Engines[id] == nil {
		apphost.Problem(w, 404, "Unknown model")
		return
	}
	a.inference.Lock()
	defer a.inference.Unlock()
	a.mu.Lock()
	defer a.mu.Unlock()
	var n int
	_ = a.H.DB.QueryRow(`SELECT count(*) FROM jobs WHERE status IN('queued','running')`).Scan(&n)
	if n > 0 {
		apphost.Problem(w, 409, "Stop or finish queued generations first")
		return
	}
	a.modelsMu.Lock()
	defer a.modelsMu.Unlock()
	if modelInstalling(a.installing[id]) {
		apphost.Problem(w, 409, "Wait for the download to finish")
		return
	}
	if e := os.RemoveAll(a.modelDir(id)); e != nil {
		apphost.Problem(w, 500, "Model could not be removed")
		return
	}
	delete(a.installing, id)
	apphost.JSON(w, 200, map[string]bool{"removed": true})
}
func installed(dir string) bool {
	raw, e := os.ReadFile(filepath.Join(dir, "installed.json"))
	if e != nil {
		return false
	}
	var marker struct {
		Version string `json:"version"`
	}
	if json.Unmarshal(raw, &marker) != nil {
		return false
	}
	for _, m := range modelList() {
		if m.ID != filepath.Base(dir) {
			continue
		}
		if marker.Version != m.Version {
			return false
		}
		for _, asset := range m.Assets {
			info, e := os.Stat(filepath.Join(dir, filepath.FromSlash(asset.Path)))
			if e != nil || !info.Mode().IsRegular() || info.Size() != asset.Size {
				return false
			}
		}
		if m.ID == "moss" {
			info, e := os.Stat(filepath.Join(dir, "preset-state.json"))
			if e != nil || info.Size() == 0 {
				return false
			}
		}
		return true
	}
	return false
}
func canClone(dir string) bool {
	b, _ := os.ReadFile(filepath.Join(dir, "installed.json"))
	var v struct {
		Cloning bool `json:"cloning"`
	}
	_ = json.Unmarshal(b, &v)
	return v.Cloning
}

type PocketEngine struct {
	A  *App
	ID string
}

func (p *PocketEngine) Available() bool { return installed(p.A.modelDir(p.ID)) }
func (p *PocketEngine) Cloning() bool   { return canClone(p.A.modelDir(p.ID)) }
func (p *PocketEngine) Compatible(v Voice, lang string) bool {
	return lang == "en" && (v.Engine == "pocket" || v.Engine == p.ID && p.Cloning())
}
func (p *PocketEngine) Generate(ctx context.Context, text string, v Voice, lang, out string) error {
	dir, e := p.A.runtime()
	if e != nil {
		return e
	}
	voice := filepath.Join(p.A.modelDir(p.ID), "alba.safetensors")
	if !v.Preset {
		voice = p.A.path("projects", "voice-"+v.ID, "state.safetensors")
	}
	return p.A.pocketService(ctx, dir, []string{"generate", "--quiet", "--config", filepath.Join(p.A.modelDir(p.ID), "config.yaml"), "--device", "cpu", "--voice", voice, "--text", text, "--output-path", out})
}
func (p *PocketEngine) Clone(ctx context.Context, src, dst string) error {
	if !p.Cloning() {
		return errors.New("This model was installed without voice cloning")
	}
	dir, e := p.A.runtime()
	if e != nil {
		return e
	}
	return p.A.pocketService(ctx, dir, []string{"export-voice", "--quiet", "--config", filepath.Join(p.A.modelDir(p.ID), "config.yaml"), src, dst + ".safetensors"})
}

type MossEngine struct{ A *App }

func (m *MossEngine) Available() bool { return installed(m.A.modelDir("moss")) }
func (m *MossEngine) Cloning() bool   { return m.Available() }
func (m *MossEngine) Compatible(v Voice, lang string) bool {
	langs := []string{"zh", "en", "de", "es", "fr", "ja", "it", "hu", "ko", "ru", "fa", "ar", "pl", "pt", "cs", "da", "sv", "el", "tr"}
	supported := false
	for _, l := range langs {
		if l == lang {
			supported = true
		}
	}
	return supported && v.Engine == "moss"
}
func (m *MossEngine) Generate(ctx context.Context, text string, v Voice, lang, out string) error {
	dir, e := m.A.runtimeFor(true)
	if e != nil {
		return e
	}
	args := []string{filepath.Join(dir, "moss", "worker.py"), "generate", m.A.modelDir("moss"), text, out, m.A.threads()}
	if !v.Preset {
		args = append(args, m.A.path("projects", "voice-"+v.ID, "state.json"))
	} else {
		args = append(args, filepath.Join(m.A.modelDir("moss"), "preset-state.json"))
	}
	return m.A.service(ctx, filepath.Join(dir, "moss", "python.exe"), args, nil)
}
func (m *MossEngine) Clone(ctx context.Context, src, dst string) error {
	dir, e := m.A.runtimeFor(true)
	if e != nil {
		return e
	}
	return m.A.service(ctx, filepath.Join(dir, "moss", "python.exe"), []string{filepath.Join(dir, "moss", "worker.py"), "clone", m.A.modelDir("moss"), src, dst + ".json", m.A.threads()}, nil)
}
