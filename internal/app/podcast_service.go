package app

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
)

const podcastGrantKey = "podcasts-speech-grant"
const pocketPodcastCast = "marius-alba-v1"

func (a *App) hasMarius(model string) bool {
	engine, ok := a.Engines[model].(*PocketEngine)
	if !ok || !engine.Available() {
		return false
	}
	for _, m := range modelList() {
		if m.ID != model {
			continue
		}
		for _, asset := range m.Assets {
			if asset.PodcastPreset {
				info, err := os.Stat(filepath.Join(a.modelDir(model), asset.Path))
				if err == nil && info.Mode().IsRegular() && info.Size() == asset.Size {
					return true
				}
			}
		}
	}
	return false
}

func (a *App) podcastModel() (model, cast string, voice Voice) {
	for _, m := range modelList() {
		if a.hasMarius(m.ID) {
			return m.ID, pocketPodcastCast, Voice{ID: "alba", Engine: "pocket", Language: "en", Preset: true}
		}
	}
	for _, model := range []string{"pocket18", "pocket7", "moss", "qwen06", "qwen17"} {
		engine := a.Engines[model]
		if engine == nil || !engine.Available() {
			continue
		}
		for _, voice := range a.presetVoices() {
			if engine.Compatible(voice, "en") {
				return model, "single-v1", voice
			}
		}
	}
	return "", "", Voice{}
}

func loopbackAddress(address string) bool {
	host, _, err := net.SplitHostPort(address)
	return err == nil && net.ParseIP(host).IsLoopback()
}

// This is a direct local service, separate from the authenticated web app.
func (a *App) podcastAuthorized(r *http.Request) bool {
	return a.localVoiceGrant(r, podcastGrantKey)
}

func (a *App) localVoiceGrant(r *http.Request, key string) bool {
	token := directSpeechToken(r)
	return token != "" && a.voiceKeyMatches(token, key)
}

func directSpeechToken(r *http.Request) string {
	if !loopbackAddress(r.Host) || !loopbackAddress(r.RemoteAddr) || r.URL.RawQuery != "" {
		return ""
	}
	for k := range r.Header {
		name := strings.ToLower(k)
		if name == "origin" || name == "forwarded" || name == "cookie" || strings.HasPrefix(name, "x-forwarded-") || strings.HasPrefix(name, "sec-fetch-") {
			return ""
		}
	}
	token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !ok || len(token) != 64 {
		return ""
	}
	return token
}

func (a *App) voiceKeyMatches(token, key string) bool {
	var want string
	if a.H.DB.QueryRow("SELECT body FROM preferences WHERE owner='' AND key=?", key).Scan(&want) != nil {
		return false
	}
	want, _, _ = strings.Cut(want, " ")
	got := sha256.Sum256([]byte(token))
	return subtle.ConstantTimeCompare([]byte(want), []byte(hex.EncodeToString(got[:]))) == 1
}

func (a *App) podcastServiceStatus(w http.ResponseWriter, r *http.Request) {
	if !a.podcastAuthorized(r) {
		apphost.Problem(w, 401, "A direct local Podcasts credential is required")
		return
	}
	model, cast, _ := a.podcastModel()
	w.Header().Set("Cache-Control", "no-store")
	apphost.JSON(w, 200, map[string]any{"protocol": 1, "ready": model != "", "model": model, "cast": cast})
}

func (a *App) podcastSpeech(w http.ResponseWriter, r *http.Request) {
	a.podcastMu.Lock()
	if !a.podcastAuthorized(r) {
		a.podcastMu.Unlock()
		apphost.Problem(w, 401, "A direct local Podcasts credential is required")
		return
	}
	if !a.inference.TryLock() {
		a.podcastMu.Unlock()
		w.Header().Set("Retry-After", "5")
		apphost.Problem(w, 429, "My Voice is recording another request")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 4*time.Minute)
	a.podcastCancel = cancel
	a.podcastMu.Unlock()
	defer func() {
		cancel()
		a.podcastMu.Lock()
		a.podcastCancel = nil
		a.inference.Unlock()
		a.podcastMu.Unlock()
	}()
	r.Body = http.MaxBytesReader(w, r.Body, 20<<10)
	var in struct {
		Text    string `json:"text"`
		Speaker string `json:"speaker"`
		Model   string `json:"model"`
		Cast    string `json:"cast,omitempty"`
	}
	if !apphost.Decode(w, r, &in) {
		return
	}
	in.Text = strings.TrimSpace(in.Text)
	if len(in.Text) == 0 || len(in.Text) > 3500 || (in.Speaker != "Alex" && in.Speaker != "Jamie") {
		apphost.Problem(w, 400, "Send one Alex or Jamie narration turn up to 3500 bytes")
		return
	}
	model, cast, voice := a.podcastModel()
	if model == "" || model != in.Model || (in.Cast != "" && in.Cast != cast) {
		apphost.Problem(w, 409, "Install a speech model in My Voice or refresh the selected model")
		return
	}
	if cast == pocketPodcastCast && in.Speaker == "Alex" {
		voice.ID = "marius"
	}
	style := Script{Speed: 1.02, Pitch: 1, Volume: 1}
	if in.Speaker == "Jamie" {
		style.Speed, style.Pitch = 0.96, -1.5
	}
	if cast == pocketPodcastCast {
		style.Speed, style.Pitch = 1, 0
	}
	b, err := a.recordStockSpeech(ctx, model, voice, in.Text, style)
	if err != nil {
		apphost.Problem(w, 502, "My Voice could not render this turn. Check the installed model and retry")
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Voice-Sample-Rate", "24000")
	_, _ = w.Write(b)
}

func (a *App) recordStockSpeech(ctx context.Context, model string, voice Voice, text string, style Script) ([]byte, error) {
	root := a.path("tmp")
	if err := os.MkdirAll(root, 0700); err != nil {
		return nil, err
	}
	dir, err := os.MkdirTemp(root, "stock-speech-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)
	raw, wav := filepath.Join(dir, "raw.wav"), filepath.Join(dir, "pcm.wav")
	if err := a.Engines[model].Generate(ctx, text, voice, "en", raw); err != nil {
		return nil, err
	}
	if err := a.processAudio(ctx, raw, wav, style); err != nil {
		return nil, err
	}
	info, err := os.Stat(wav)
	if err != nil {
		return nil, err
	}
	if info.Size() > 16<<20 {
		return nil, errors.New("Speech exceeded the PCM size limit")
	}
	b, rate, err := pcm(wav)
	if err != nil || rate != 24000 || len(b) < 1000 || len(b)%2 != 0 {
		return nil, errors.New("Speech was not mono 24 kHz PCM")
	}
	return b, ctx.Err()
}
