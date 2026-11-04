package app

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type testEngine struct{}

func (testEngine) Available() bool                                               { return true }
func (testEngine) Cloning() bool                                                 { return true }
func (testEngine) Compatible(Voice, string) bool                                 { return true }
func (testEngine) Generate(context.Context, string, Voice, string, string) error { return nil }
func (testEngine) Clone(context.Context, string, string) error                   { return nil }
func TestChunkingAndAudioBoundaries(t *testing.T) {
	if got, e := transform("OWS", mergeRules([]Rule{{From: "OWS", To: "default"}}, []Rule{{From: "OWS", To: "project"}})); e != nil || got != "project" {
		t.Fatal("project dictionary did not override account rules", got)
	}
	if _, e := transform(strings.Repeat("a", 500), []Rule{{From: "a", To: strings.Repeat("a", 500)}}); e == nil {
		t.Fatal("unbounded replacement expansion")
	}
	text := "First sentence.\nSecond sentence! " + strings.Repeat("Longword ", 180)
	parts := splitScript(Script{Text: text, Voice: "alba"})
	if len(parts) < 4 {
		t.Fatal("long text not chunked", parts)
	}
	for _, s := range parts {
		if len(s.Text) > 500 || s.Voice != "alba" || !safeID.MatchString(s.ID) {
			t.Fatal(s)
		}
	}
	root := t.TempDir()
	a := filepath.Join(root, "a.wav")
	b := filepath.Join(root, "b.wav")
	out := filepath.Join(root, "out.wav")
	if e := writeWAV(a, make([]byte, 4800), 24000); e != nil {
		t.Fatal(e)
	}
	if e := writeWAV(b, make([]byte, 4800), 24000); e != nil {
		t.Fatal(e)
	}
	if e := joinWAV([]string{a, b}, out, .2); e != nil {
		t.Fatal(e)
	}
	info, e := inspectWAV(out)
	if e != nil || info.Duration != .4 {
		t.Fatal(info, e)
	}
}
func TestModelRegistryIsPinnedAndOpen(t *testing.T) {
	models := modelList()
	if len(models) != 5 {
		t.Fatal("expected five explicit CPU models")
	}
	for _, m := range models {
		if len(m.Languages) == 0 || m.Size <= 0 || len(m.Assets) == 0 {
			t.Fatal("incomplete model", m.ID)
		}
		for _, a := range m.Assets {
			if a.Gated || len(a.SHA) != 64 || strings.Contains(a.SHA, "*") || !strings.HasPrefix(a.URL, "https://huggingface.co/") || a.Size <= 0 {
				t.Fatal("unverified or gated asset", m.ID, a.Path)
			}
		}
	}
	if len(models[2].Languages) != 19 {
		t.Fatal("only explicitly listed upstream language codes may be exposed")
	}
}
func TestOwnerScopeCSRFQueueAndRevocation(t *testing.T) {
	root := os.Getenv("MYVOICE_TEST_ROOT")
	if root == "" {
		t.Fatal("Set MYVOICE_TEST_ROOT to an isolated synthetic NTFS test root")
	}
	data, e := os.MkdirTemp(root, "api-")
	if e != nil {
		t.Fatal(e)
	}
	host, e := apphost.New(apphost.Config{ID: "voice", Name: "My Voice", Version: "test", Data: data, Listen: "127.0.0.1:29389", Origin: "http://localhost:29389"}, nil)
	if e != nil {
		t.Fatal(e)
	}
	defer host.Close()
	a, e := New(host)
	if e != nil {
		t.Fatal(e)
	}
	a.Engines["pocket18"] = testEngine{}
	if e = host.Bootstrap("owner@example.test", "Isolated-fixture-4721"); e != nil {
		t.Fatal(e)
	}
	server := httptest.NewServer(host.Handler())
	defer server.Close()
	call := func(method, path, body, cookie, csrf string) (int, []byte, *http.Response) {
		r, _ := http.NewRequest(method, server.URL+path, strings.NewReader(body))
		r.Host = "localhost:29389"
		r.Header.Set("Origin", "http://localhost:29389")
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Cookie", cookie)
		r.Header.Set("X-CSRF-Token", csrf)
		res, e := server.Client().Do(r)
		if e != nil {
			t.Fatal(e)
		}
		defer res.Body.Close()
		var out bytes.Buffer
		out.ReadFrom(res.Body)
		return res.StatusCode, out.Bytes(), res
	}
	status, _, login := call("POST", "/api/login", `{"email":"owner@example.test","password":"Isolated-fixture-4721"}`, "", "")
	if status != 200 {
		t.Fatal(status)
	}
	cookie := login.Cookies()[0].String()
	_, raw, _ := call("GET", "/api/me", "", cookie, "")
	var session apphost.Session
	_ = json.Unmarshal(raw, &session)
	_, raw, _ = call("GET", "/api/models", "", cookie, "")
	var models []Model
	if e := json.Unmarshal(raw, &models); e != nil {
		t.Fatal(e)
	}
	for _, model := range models {
		if model.ID == "moss" && (!model.Cloning || model.CloningReady || model.Installed) {
			t.Fatal("uninstalled MOSS must advertise capability separately from readiness", model)
		}
	}
	// Preparation can take longer than the transfer. Neither another installer
	// nor removal may race the first real model/reference preparation.
	a.installing["moss"] = "Preparing cloning model…"
	if c, _, _ := call("POST", "/api/models/moss/install", `{"license":true}`, cookie, session.CSRF); c != 409 {
		t.Fatal("concurrent preparation install", c)
	}
	if c, _, _ := call("DELETE", "/api/models/moss", "", cookie, session.CSRF); c != 409 {
		t.Fatal("concurrent preparation removal", c)
	}
	delete(a.installing, "moss")
	if c, _, _ := call("GET", "/api/voices", "", "", ""); c != 401 {
		t.Fatal("anonymous voice library", c)
	}
	if c, _, _ := call("POST", "/api/tts/jobs", `{"text":"Hello"}`, cookie, ""); c != 403 {
		t.Fatal("CSRF bypass", c)
	}
	a.inference.Lock()
	busy, _, _ := call("POST", "/api/voices", "", cookie, session.CSRF)
	a.inference.Unlock()
	if busy != 409 {
		t.Fatal("reference preparation must reject concurrent costly work before reading audio", busy)
	}
	v := Voice{ID: "privatevoice", Name: "Private", Engine: "pocket18", Owner: session.User.ID, Consent: now()}
	_, e = host.DB.Exec(`INSERT INTO voices(id,owner,body)VALUES(?,?,?)`, v.ID, v.Owner, jsonText(v))
	if e != nil {
		t.Fatal(e)
	}
	if _, e = a.findVoice(v.ID, "other"); e == nil {
		t.Fatal("ungranted private voice")
	}
	_, e = host.DB.Exec(`INSERT INTO grants(voice,user)VALUES(?,?)`, v.ID, "other")
	if e != nil {
		t.Fatal(e)
	}
	if _, e = a.findVoice(v.ID, "other"); e != nil {
		t.Fatal(e)
	}
	_, _ = host.DB.Exec(`DELETE FROM grants WHERE voice=?`, v.ID)
	if _, e = a.findVoice(v.ID, "other"); e == nil {
		t.Fatal("revoked voice usable")
	}
	first, e := a.enqueue(Script{RequestID: "repeated-request", Text: "Hello", Voice: "alba"}, session.User.ID)
	if e != nil {
		t.Fatal(e)
	}
	again, e := a.enqueue(Script{RequestID: "repeated-request", Text: "Hello", Voice: "alba"}, session.User.ID)
	if e != nil || first.ID != again.ID {
		t.Fatal("lost-response retry queued duplicate work")
	}
	var j Job
	for i := 0; i < 19; i++ {
		j, e = a.enqueue(Script{Text: "Hello", Voice: "alba"}, session.User.ID)
		if e != nil {
			t.Fatal(e)
		}
	}
	if _, e = a.enqueue(Script{Text: "Overflow"}, session.User.ID); e == nil {
		t.Fatal("unbounded generation queue")
	}
	if c, _, _ := call("GET", "/api/outputs/"+j.ID, "", "", ""); c != 401 {
		t.Fatal("private audio exposed", c)
	}
	_, _ = host.DB.Exec(`UPDATE jobs SET owner='another' WHERE id=?`, j.ID)
	if c, _, _ := call("GET", "/api/tts/jobs/"+j.ID, "", cookie, session.CSRF); c != 404 {
		t.Fatal("cross-owner job read", c)
	}
	if c, _, _ := call("POST", "/api/projects", `{"id":"../escape","text":"hi"}`, cookie, session.CSRF); c != 400 {
		t.Fatal("project path traversal", c)
	}
}

func TestNativeSchemaUpgradeAndRestart(t *testing.T) {
	cfg := apphost.Config{ID: "voice", Name: "My Voice", Version: "test", Data: t.TempDir(), Listen: "127.0.0.1:29388", Origin: "http://localhost:29388"}
	h, e := apphost.New(cfg, nil)
	if e != nil {
		t.Fatal(e)
	}
	if e = h.Migrate(1, appMigrations[0]); e != nil {
		t.Fatal(e)
	}
	if _, e = New(h); e != nil {
		t.Fatal("upgrade from schema one", e)
	}
	if _, e = h.DB.Exec(`INSERT INTO projects(id,owner,body,updated)VALUES('preserved','synthetic','{}','now')`); e != nil {
		t.Fatal(e)
	}
	h.Close()
	h, e = apphost.New(cfg, nil)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = New(h); e != nil {
		h.Close()
		t.Fatal("restart at latest schema", e)
	}
	var n int
	if e = h.DB.QueryRow(`SELECT count(*) FROM projects WHERE id='preserved'`).Scan(&n); e != nil || n != 1 {
		t.Fatal("project not preserved", e)
	}
	_, e = h.DB.Exec(`INSERT INTO host_migrations VALUES(3,'future')`)
	if e != nil {
		t.Fatal(e)
	}
	h.Close()
	h, e = apphost.New(cfg, nil)
	if e != nil {
		t.Fatal(e)
	}
	defer h.Close()
	if _, e = New(h); e == nil {
		t.Fatal("newer schema downgrade accepted")
	}
}
