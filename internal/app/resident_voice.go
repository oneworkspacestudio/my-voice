package app

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
)

const wonderGrantKey = "wonderpal-speech-grant"
const wonderModelKey = "wonderpal-speech-model"
const phononModelURL = "https://huggingface.co/gradium/phonon-31m-en/resolve/92741c5b07d085ff840e26fea8ce43110b87852f/"

var phononAssets = []Asset{
	{Path: "config.json", Size: 4634, SHA: "b5d0f88353638370d94035b7d73d7d6577c0554d78f8d7ecfd40490aae39ae85"},
	{Path: "tokenizer.json", Size: 277967, SHA: "691a00068fa7b989f2f3e5ae5daa2d1a5bb319f2d5360deafcfb19a8bb6e7936"},
	{Path: "default-voice.safetensors", Size: 256136, SHA: "dc537478eaadc1a5257fcda2e3501ab81e14c111a5d95a5838f350db2afd0f1a"},
	{Path: "model.q8.gguf", Size: 64512352, SHA: "a7d9aad2cd6dc66473f24a08bc6d2939030ef4d751ef82349361fbef3272dbae"},
}

var phononHTTP = &http.Client{
	Transport:     &http.Transport{DisableKeepAlives: true, ResponseHeaderTimeout: 15 * time.Second},
	Timeout:       45 * time.Second,
	CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
}

type residentVoice struct {
	app        *App
	changed    chan struct{}
	mu         sync.Mutex
	render     sync.Mutex
	worker     *residentWorker
	state      string
	fixedModel string
	requested  bool
}

type residentWorker struct {
	ctx    context.Context
	stop   context.CancelFunc
	url    string
	secret string
	model  string
}

func (v *residentVoice) selected() string {
	if v.fixedModel != "" {
		return v.fixedModel
	}
	var model string
	_ = v.app.H.DB.QueryRow("SELECT body FROM preferences WHERE owner='' AND key=?", wonderModelKey).Scan(&model)
	if model == "" {
		return "phonon"
	}
	return model
}

func residentModelName(id string) string {
	switch id {
	case "phonon":
		return "Phonon 31M"
	case "pocket18":
		return "Pocket TTS 18M"
	case "pocket7":
		return "Pocket TTS 7M"
	case "moss":
		return "MOSS TTS Nano"
	case "qwen06":
		return "Qwen3 TTS 0.6B"
	case "qwen17":
		return "Qwen3 TTS 1.7B"
	}
	return ""
}

func (v *residentVoice) available(id string) bool {
	return id == "phonon" || (residentModelName(id) != "" && v.app.Engines[id] != nil && v.app.Engines[id].Available())
}

// The scoped local grant can select only this service's stock model allowlist.
// It does not expose owner sessions, private voices, downloads or configuration.
func (v *residentVoice) modelControl(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPatch {
		v.selectModel(w, r)
		return
	}
	v.mu.Lock()
	defer v.mu.Unlock()
	if !v.app.localVoiceGrant(r, wonderGrantKey) {
		apphost.Problem(w, 401, "A direct local WonderPal voice credential is required")
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	state := v.modelState()
	state["protocol"] = 1
	apphost.JSON(w, 200, state)
}

func (v *residentVoice) selectModel(w http.ResponseWriter, r *http.Request) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if !v.app.localVoiceGrant(r, wonderGrantKey) {
		apphost.Problem(w, 401, "A direct local WonderPal voice credential is required")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	var in struct {
		Model string `json:"model"`
	}
	if !apphost.Decode(w, r, &in) {
		return
	}
	if residentModelName(in.Model) == "" || in.Model == "qwen17" {
		apphost.Problem(w, 400, "Choose a supported WonderPal speech model")
		return
	}
	if !v.available(in.Model) {
		apphost.Problem(w, 409, "Download this model in Settings → Models first")
		return
	}
	if in.Model != v.selected() {
		if _, err := v.app.H.DB.Exec("INSERT INTO preferences(owner,key,body) VALUES('',?,?) ON CONFLICT(owner,key) DO UPDATE SET body=excluded.body", wonderModelKey, in.Model); err != nil {
			v.app.H.DBError(w, err)
			return
		}
		if v.worker != nil {
			v.worker.stop()
			v.worker = nil
		}
		v.state = ""
		select {
		case v.changed <- struct{}{}:
		default:
		}
		v.app.H.Audit("legacy-speech", "speech-model", in.Model)
	}
	w.Header().Set("Cache-Control", "no-store")
	apphost.JSON(w, 200, map[string]string{"model": in.Model})
}

func (v *residentVoice) enabled() bool {
	if v.fixedModel != "" {
		return v.requested
	}
	var enabled bool
	_ = v.app.H.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM preferences WHERE owner='' AND key=?)", wonderGrantKey).Scan(&enabled)
	return enabled
}

func (v *residentVoice) reset() {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.worker != nil {
		v.worker.stop()
		v.worker = nil
	}
	v.state = ""
	v.wake()
}

func (v *residentVoice) wake() {
	select {
	case v.changed <- struct{}{}:
	default:
	}
}

func (v *residentVoice) startFixed() {
	v.mu.Lock()
	defer v.mu.Unlock()
	if !v.requested {
		v.requested = true
		v.wake()
	}
}

// Caller holds mu. This is the complete safe stock-model control response.
func (v *residentVoice) modelState() map[string]any {
	choices := []map[string]any{}
	for _, id := range []string{"phonon", "pocket18", "pocket7", "moss", "qwen06"} {
		choices = append(choices, map[string]any{"value": id, "label": residentModelName(id), "installed": v.available(id)})
	}
	return map[string]any{"status": v.state, "model": v.selected(), "choices": choices, "ready": v.worker != nil && v.worker.url != ""}
}

func (v *residentVoice) run(ctx context.Context) {
	for ctx.Err() == nil {
		v.mu.Lock()
		var worker *residentWorker
		if v.enabled() {
			model := v.selected()
			child, stop := context.WithCancel(ctx)
			worker = &residentWorker{ctx: child, stop: stop, secret: rand.Text(), model: model}
			v.worker = worker
			v.state = "Loading " + residentModelName(model)
		}
		v.mu.Unlock()
		if worker != nil {
			v.serve(worker)
			worker.stop()
			v.mu.Lock()
			if v.worker == worker {
				v.worker = nil
				v.state = residentModelName(worker.model) + " is unavailable. Check its download in Models; retrying"
			}
			v.mu.Unlock()
		}
		select {
		case <-ctx.Done():
			return
		case <-v.changed:
		case <-time.After(10 * time.Second):
		}
	}
}

func (v *residentVoice) serve(worker *residentWorker) {
	a := v.app
	model := a.modelDir(worker.model)
	if worker.model == "phonon" {
		model = a.modelDir("phonon-wonderpal")
	}
	if !v.available(worker.model) {
		return
	}
	if os.MkdirAll(model, 0700) != nil || os.MkdirAll(a.path("tmp"), 0700) != nil {
		return
	}
	client := &http.Client{Timeout: 4 * time.Minute}
	if worker.model == "phonon" {
		for _, asset := range phononAssets {
			file := filepath.Join(model, asset.Path)
			if !verifiedModelFile(file, asset) {
				asset.URL = phononModelURL + asset.Path
				if downloadModelAsset(worker.ctx, client, asset, file, func(int64) {}) != nil {
					return
				}
			}
		}
	}
	runtime, err := a.runtimeFor(true)
	if err != nil || worker.ctx.Err() != nil {
		return
	}
	tmp, err := os.MkdirTemp(a.path("tmp"), "phonon-worker-")
	if err != nil {
		return
	}
	defer os.RemoveAll(tmp)
	marker := filepath.Join(tmp, "ready")
	python := filepath.Join(runtime, "moss", "python.exe")
	script := filepath.Join(runtime, "moss", "phonon_service.py")
	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = apphost.Service(worker.ctx, python, tmp, []string{script, worker.model, model, marker, a.threads()}, []string{
			"HF_HUB_OFFLINE=1", "HF_HOME=" + a.path("models", "hf"), "USERPROFILE=" + a.H.Config.Data,
			"RAYON_NUM_THREADS=" + a.threads(), "OMP_NUM_THREADS=" + a.threads(), "MYVOICE_WORKER_TOKEN=" + worker.secret,
		}, io.Discard)
	}()
	defer func() { worker.stop(); <-done }()
	start := time.NewTimer(90 * time.Second)
	defer start.Stop()
	poll := time.NewTicker(50 * time.Millisecond)
	defer poll.Stop()
	for {
		select {
		case <-worker.ctx.Done():
			return
		case <-done:
			return
		case <-start.C:
			return
		case <-poll.C:
			ready, err := os.ReadFile(marker)
			if err != nil {
				continue
			}
			port, err := strconv.ParseUint(string(ready), 10, 16)
			if err != nil || port == 0 {
				return
			}
			v.mu.Lock()
			if v.worker == worker {
				worker.url = "http://127.0.0.1:" + strconv.FormatUint(port, 10) + "/speech"
				v.state = residentModelName(worker.model) + " is loaded and ready"
			}
			v.mu.Unlock()
			select {
			case <-worker.ctx.Done():
			case <-done:
			}
			return
		}
	}
}

func (v *residentVoice) speak(w http.ResponseWriter, r *http.Request) {
	if !v.app.localVoiceGrant(r, wonderGrantKey) {
		apphost.Problem(w, 401, "A direct local WonderPal voice credential is required")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 8192)
	var in struct {
		Model  string `json:"model"`
		Input  string `json:"input"`
		Voice  string `json:"voice"`
		Format string `json:"response_format"`
	}
	if !apphost.Decode(w, r, &in) {
		return
	}
	if in.Input = strings.TrimSpace(in.Input); in.Input == "" || len(in.Input) > 1600 || in.Voice != "default" || in.Format != "pcm" {
		apphost.Problem(w, 400, "Send text of up to 1600 bytes with default voice and pcm output")
		return
	}
	v.stream(w, r, in.Input)
}

func (v *residentVoice) stream(w http.ResponseWriter, r *http.Request, input string) {
	v.mu.Lock()
	if v.fixedModel == "" && !v.app.localVoiceGrant(r, wonderGrantKey) {
		v.mu.Unlock()
		apphost.Problem(w, 401, "A direct local speech credential is required")
		return
	}
	worker := v.worker
	if worker == nil || worker.url == "" {
		v.mu.Unlock()
		w.Header().Set("Retry-After", "2")
		apphost.Problem(w, 503, "The speech model is warming up")
		return
	}
	url := worker.url
	v.mu.Unlock()
	if !v.render.TryLock() {
		apphost.Problem(w, 429, "Workspace Voice is speaking another reply")
		return
	}
	defer v.render.Unlock()
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	stop := context.AfterFunc(worker.ctx, cancel)
	defer stop()
	if worker.ctx.Err() != nil {
		return
	}
	body, _ := json.Marshal(map[string]string{"input": input})
	request, _ := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(body)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+worker.secret)
	response, err := phononHTTP.Do(request)
	if err != nil {
		apphost.Problem(w, 502, "The Workspace Voice worker is unavailable")
		return
	}
	defer response.Body.Close()
	if response.StatusCode != 200 || response.Header.Get("Content-Type") != "audio/pcm" {
		apphost.Problem(w, 502, "The Workspace Voice worker could not render audio")
		return
	}
	w.Header().Set("Content-Type", "audio/pcm")
	w.Header().Set("Cache-Control", "no-store, no-transform")
	w.Header().Set("X-Accel-Buffering", "no")
	stream := http.NewResponseController(w)
	buf := make([]byte, 8192)
	left := 8 << 20
	for {
		n, err := response.Body.Read(buf)
		if n > 0 {
			left -= n
			if left < 0 {
				panic(http.ErrAbortHandler)
			}
			if _, e := w.Write(buf[:n]); e != nil || stream.Flush() != nil {
				return
			}
		}
		if err == io.EOF {
			return
		}
		if err != nil {
			panic(http.ErrAbortHandler)
		}
	}
}
