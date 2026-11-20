package app

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
)

const speechKeyPrefix = "speech-api-key:"

type speechAPI struct {
	app       *App
	phonon    *residentVoice
	mu        sync.Mutex
	activeKey string
	cancel    context.CancelFunc
}

type speechKey struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

func (s *speechAPI) keyList(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	rows, err := s.app.H.DB.Query("SELECT substr(key, ?), substr(body, 66) FROM preferences WHERE owner='' AND key LIKE ? ORDER BY key", len(speechKeyPrefix)+1, speechKeyPrefix+"%")
	if err != nil {
		s.app.H.DBError(w, err)
		return
	}
	defer rows.Close()
	keys := []speechKey{}
	for rows.Next() {
		var key speechKey
		if err := rows.Scan(&key.ID, &key.Label); err != nil {
			s.app.H.DBError(w, err)
			return
		}
		keys = append(keys, key)
	}
	if err := rows.Err(); err != nil {
		s.app.H.DBError(w, err)
		return
	}
	rows.Close()
	for _, key := range []speechKey{{"legacy-1", "Existing read-aloud key"}, {"legacy-2", "Existing narration key"}} {
		var exists bool
		if err := s.app.H.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM preferences WHERE owner='' AND key=?)", speechKeyStorage(key.ID)).Scan(&exists); err != nil {
			s.app.H.DBError(w, err)
			return
		}
		if exists {
			keys = append(keys, key)
		}
	}
	apphost.JSON(w, 200, map[string]any{"keys": keys, "origin": "http://" + s.app.H.Config.Listen})
}

func speechKeyStorage(id string) string {
	switch id {
	case "legacy-1":
		return wonderGrantKey
	case "legacy-2":
		return podcastGrantKey
	default:
		return speechKeyPrefix + id
	}
}

func (s *speechAPI) issueKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	var in struct{ Label string }
	if !apphost.Decode(w, r, &in) {
		return
	}
	in.Label = strings.TrimSpace(in.Label)
	if in.Label == "" || len(in.Label) > 100 {
		apphost.Problem(w, 400, "Enter a label of up to 100 bytes")
		return
	}
	var secret [32]byte
	if _, err := rand.Read(secret[:]); err != nil {
		apphost.Problem(w, 500, "Could not issue an API key")
		return
	}
	token := hex.EncodeToString(secret[:])
	id := token[:24]
	digest := sha256.Sum256([]byte(token))
	_, err := s.app.H.DB.Exec("INSERT INTO preferences(owner,key,body) VALUES('',?,?)", speechKeyStorage(id), hex.EncodeToString(digest[:])+" "+in.Label)
	if err != nil {
		s.app.H.DBError(w, err)
		return
	}
	s.app.H.Audit(uid(r), "speech-key-issued", id)
	apphost.JSON(w, 201, map[string]string{"credential": token})
}

func (s *speechAPI) revokeKey(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if !safeID.MatchString(id) {
		apphost.Problem(w, 400, "Invalid API key")
		return
	}
	key := speechKeyStorage(id)
	s.mu.Lock()
	_, err := s.app.H.DB.Exec("DELETE FROM preferences WHERE owner='' AND key=?", key)
	if err == nil && s.activeKey == key {
		s.cancel()
	}
	s.mu.Unlock()
	if err != nil {
		s.app.H.DBError(w, err)
		return
	}
	switch key {
	case wonderGrantKey:
		s.app.wonderVoice.reset()
	case podcastGrantKey:
		s.app.podcastMu.Lock()
		if s.app.podcastCancel != nil {
			s.app.podcastCancel()
		}
		s.app.podcastMu.Unlock()
	}
	s.app.H.Audit(uid(r), "speech-key-revoked", id)
	w.Header().Set("Cache-Control", "no-store")
	apphost.JSON(w, 200, map[string]bool{"revoked": true})
}

func (s *speechAPI) authorizedKey(r *http.Request) string {
	token := directSpeechToken(r)
	if token == "" {
		return ""
	}
	key := speechKeyStorage(token[:24])
	for _, candidate := range []string{key, wonderGrantKey, podcastGrantKey} {
		if s.app.voiceKeyMatches(token, candidate) {
			return candidate
		}
	}
	return ""
}

type speechChoice struct {
	Value     string   `json:"value"`
	Label     string   `json:"label"`
	Installed bool     `json:"installed"`
	Voices    []string `json:"voices"`
}

func (s *speechAPI) choices() []speechChoice {
	choices := []speechChoice{{"phonon", "Phonon 31M", true, []string{"default"}}}
	for _, id := range []string{"pocket18", "pocket7", "moss", "qwen06", "qwen17"} {
		choice := speechChoice{id, residentModelName(id), false, []string{"default"}}
		engine := s.app.Engines[id]
		if engine != nil && engine.Available() {
			choice.Installed = true
			for _, voice := range s.app.presetVoices() {
				if engine.Compatible(voice, "en") {
					choice.Voices = append(choice.Voices, voice.ID)
				}
			}
			if s.app.hasMarius(id) {
				choice.Voices = append(choice.Voices, "marius")
			}
		}
		choices = append(choices, choice)
	}
	return choices
}

func (s *speechAPI) models(w http.ResponseWriter, r *http.Request) {
	key := s.authorizedKey(r)
	if key == "" {
		apphost.Problem(w, 401, "A direct local speech API key is required")
		return
	}
	out := map[string]any{"protocol": 2, "choices": s.choices()}
	if key == wonderGrantKey {
		// Existing clients can adopt their previous selection as their starting value.
		out["previousModel"] = s.app.wonderVoice.selected()
	}
	w.Header().Set("Cache-Control", "no-store")
	apphost.JSON(w, 200, out)
}

func (s *speechAPI) speak(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 20<<10)
	var in struct {
		Model  string  `json:"model"`
		Input  string  `json:"input"`
		Voice  string  `json:"voice"`
		Format string  `json:"response_format"`
		Speed  float64 `json:"speed"`
		Pitch  float64 `json:"pitch"`
	}
	s.mu.Lock()
	key := s.authorizedKey(r)
	if key == "" {
		s.mu.Unlock()
		apphost.Problem(w, 401, "A direct local speech API key is required")
		return
	}
	if s.cancel != nil {
		s.mu.Unlock()
		w.Header().Set("Retry-After", "2")
		apphost.Problem(w, 429, "My Voice is processing another speech request")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 4*time.Minute)
	s.activeKey, s.cancel = key, cancel
	s.mu.Unlock()
	defer func() {
		cancel()
		s.mu.Lock()
		s.activeKey, s.cancel = "", nil
		s.mu.Unlock()
	}()
	if !apphost.Decode(w, r, &in) {
		return
	}
	in.Input = strings.TrimSpace(in.Input)
	if in.Speed == 0 {
		in.Speed = 1
	}
	if in.Input == "" || len(in.Input) > 3500 || in.Format != "pcm" || in.Speed < 0.5 || in.Speed > 2 || in.Pitch < -12 || in.Pitch > 12 {
		apphost.Problem(w, 400, "Send up to 3500 bytes of text with PCM output and supported speed and pitch")
		return
	}
	if in.Model == "phonon" {
		if in.Voice != "default" || in.Speed != 1 || in.Pitch != 0 || len(in.Input) > 1600 {
			apphost.Problem(w, 400, "Phonon supports the default preset, normal speed and pitch, and up to 1600 bytes")
			return
		}
		s.phonon.startFixed()
		w.Header().Set("X-Voice-Sample-Rate", "24000")
		s.phonon.stream(w, r.WithContext(ctx), in.Input)
		return
	}
	engine := s.app.Engines[in.Model]
	if engine == nil || !engine.Available() {
		apphost.Problem(w, 409, "Choose an installed model")
		return
	}
	var voice Voice
	for _, preset := range s.app.presetVoices() {
		if engine.Compatible(preset, "en") && (in.Voice == "default" || in.Voice == preset.ID) {
			voice = preset
			break
		}
	}
	if in.Voice == "marius" && s.app.hasMarius(in.Model) {
		voice = Voice{ID: "marius", Engine: "pocket", Language: "en", Preset: true}
	}
	if voice.ID == "" {
		apphost.Problem(w, 400, "Choose a supported stock voice for this model")
		return
	}
	if !s.app.inference.TryLock() {
		w.Header().Set("Retry-After", "2")
		apphost.Problem(w, 429, "My Voice is processing another speech request")
		return
	}
	defer s.app.inference.Unlock()
	data, err := s.app.recordStockSpeech(ctx, in.Model, voice, in.Input, Script{Speed: in.Speed, Pitch: in.Pitch, Volume: 1})
	if err != nil || ctx.Err() != nil {
		apphost.Problem(w, 502, "My Voice could not generate the audio")
		return
	}
	w.Header().Set("Content-Type", "audio/pcm")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Voice-Sample-Rate", "24000")
	_, _ = w.Write(data)
}
