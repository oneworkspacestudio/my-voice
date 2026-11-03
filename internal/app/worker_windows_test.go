//go:build windows

package app

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
)

// Native child processes exercise the real suspended-start Job Object and
// Windows exception exit status, without executing unsupported instructions.
func TestMain(m *testing.M) {
	name := filepath.Base(os.Args[0])
	if name == "pocket-tts.exe" && len(os.Args) == 2 && os.Args[1] == "test-exit1" {
		_, _ = os.Stderr.WriteString("private-speech-must-not-escape")
		os.Exit(1)
	}
	if len(os.Args) == 2 && os.Args[1] == "test-fallback" {
		switch name {
		case "pocket-tts.exe":
			os.Exit(int(uint32(0xc000001d)))
		case "pocket-tts-portable.exe":
			if os.Getenv("MYVOICE_PRIVATE_SENTINEL") != "" {
				os.Exit(9)
			}
			_ = os.WriteFile("portable-ran.json", []byte(`{"portable":true}`), 0600)
			os.Exit(0)
		}
	}
	os.Exit(m.Run())
}

func TestWorkerFallbackNative(t *testing.T) {
	data := t.TempDir()
	host, err := apphost.New(apphost.Config{ID: "voice", Name: "My Voice", Version: "test", Data: filepath.Join(data, "private"), Listen: "127.0.0.1:29389", Origin: "http://localhost:29389"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer host.Close()
	a := &App{H: host}
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(exe)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"pocket-tts.exe", "pocket-tts-portable.exe"} {
		if err := os.WriteFile(filepath.Join(data, name), raw, 0600); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("MYVOICE_PRIVATE_SENTINEL", "must-not-reach-worker")
	err = a.service(context.Background(), filepath.Join(data, "pocket-tts.exe"), []string{"test-fallback"}, nil)
	var failure *workerFailure
	if !errors.As(err, &failure) || failure.code != 0xc000001d || !failure.permanent() {
		t.Fatalf("Windows exception status lost: %v", err)
	}
	if err := a.pocketService(context.Background(), data, []string{"test-fallback"}); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(a.path("tmp", "portable-ran.json"))
	var ran map[string]bool
	if err != nil || json.Unmarshal(b, &ran) != nil || !ran["portable"] {
		t.Fatal("portable worker did not run", err)
	}
	if err := os.Remove(a.path("tmp", "portable-ran.json")); err != nil {
		t.Fatal(err)
	}
	err = a.pocketService(context.Background(), data, []string{"test-exit1"})
	if !errors.As(err, &failure) || failure.code != 1 || strings.Contains(err.Error(), "private-speech") {
		t.Fatal("ordinary worker failure lost its safe status", err)
	}
	if _, err := os.Stat(a.path("tmp", "portable-ran.json")); !os.IsNotExist(err) {
		t.Fatal("ordinary model errors must not trigger the CPU fallback")
	}
	err = a.service(context.Background(), filepath.Join(data, "missing-private-path.exe"), nil, nil)
	if !errors.As(err, &failure) || failure.kind != "missing" || strings.Contains(err.Error(), data) {
		t.Fatal("missing-worker response must be actionable and omit private paths", err)
	}
}

func TestPortablePocketNativeSpeech(t *testing.T) {
	exe, models := os.Getenv("MYVOICE_TEST_PORTABLE_POCKET"), os.Getenv("MYVOICE_TEST_POCKET_MODEL")
	if exe == "" || models == "" {
		t.Skip("Set portable worker and isolated pinned model paths for real CPU acceptance")
	}
	host, err := apphost.New(apphost.Config{ID: "voice", Name: "My Voice", Version: "test", Data: filepath.Join(t.TempDir(), "private"), Listen: "127.0.0.1:29389", Origin: "http://localhost:29389"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer host.Close()
	a := &App{H: host}
	out := a.path("tmp", "portable.wav")
	if err := a.service(context.Background(), exe, []string{"generate", "--quiet", "--config", filepath.Join(models, "config.yaml"), "--device", "cpu", "--voice", filepath.Join(models, "alba.safetensors"), "--text", "hello world", "--output-path", out}, nil); err != nil {
		t.Fatal(err)
	}
	info, err := inspectWAV(out)
	if err != nil || info.Duration < .1 || info.RMS <= .0001 || info.Rate != 24000 {
		t.Fatalf("portable worker did not produce audible 24 kHz PCM: %+v, %v", info, err)
	}
	t.Logf("portable native CPU: %.3fs audio, RMS %.4f", info.Duration, info.RMS)
}
