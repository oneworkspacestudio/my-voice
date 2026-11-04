package app

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestPinnedDownloadResumeAndVerification(t *testing.T) {
	data := []byte("a complete pinned model, not a partial installation")
	for _, mode := range []string{"range", "restart", "bad-range", "checksum", "interrupted"} {
		t.Run(mode, func(t *testing.T) {
			dest := filepath.Join(t.TempDir(), "model.gguf")
			if err := os.WriteFile(dest+".part", data[:9], 0600); err != nil {
				t.Fatal(err)
			}
			requests := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests++
				if requests == 1 && r.Header.Get("Range") != "bytes=9-" {
					t.Error("saved partial was not resumed")
				}
				if mode == "bad-range" {
					w.Header().Set("Content-Range", "bytes 0-9/10")
					w.WriteHeader(206)
					return
				}
				if mode == "restart" {
					w.Write(data)
					return
				}
				start := 9
				if requests > 1 {
					start = 14
				}
				w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, len(data)-1, len(data)))
				w.Header().Set("Content-Length", fmt.Sprint(len(data)-start))
				w.WriteHeader(206)
				if mode == "interrupted" && requests == 1 {
					w.Write(data[start:14])
					return
				}
				body := append([]byte(nil), data[start:]...)
				if mode == "checksum" {
					body[0] ^= 1
				}
				w.Write(body)
			}))
			defer server.Close()
			a := Asset{URL: server.URL, Size: int64(len(data)), SHA: fmt.Sprintf("%x", sha256.Sum256(data))}
			progress := int64(0)
			show := func(n int64) { progress = n }
			err := downloadModelAsset(context.Background(), server.Client(), a, dest, show)
			if mode == "bad-range" || mode == "checksum" || mode == "interrupted" {
				if err == nil {
					t.Fatal("invalid download was accepted")
				}
				if _, e := os.Stat(dest); !os.IsNotExist(e) {
					t.Fatal("failed download promoted")
				}
				if mode == "checksum" {
					if _, e := os.Stat(dest + ".part"); !os.IsNotExist(e) {
						t.Fatal("corrupted bytes retained")
					}
				}
				if mode != "interrupted" {
					return
				}
				if info, e := os.Stat(dest + ".part"); e != nil || info.Size() != 14 {
					t.Fatal("interrupted bytes not retained", e)
				}
				err = downloadModelAsset(context.Background(), server.Client(), a, dest, show)
			}
			if err != nil || !verifiedModelFile(dest, a) || progress != a.Size {
				t.Fatal("verified download did not complete", err, progress)
			}
		})
	}
}

func TestCloningModelsKeepReferencesAndLanguagesSeparate(t *testing.T) {
	for _, id := range []string{"qwen06", "qwen17"} {
		q := &QwenEngine{ID: id}
		m := q.definition()
		if !m.Cloning || len(m.Languages) != 10 || m.ShortName == "" || m.RAMGB < 4 {
			t.Fatal("incomplete model", id)
		}
		if !q.Compatible(Voice{Engine: id}, "en") || !q.Compatible(Voice{Engine: id}, "zh") {
			t.Fatal("supported language rejected", id)
		}
		if q.Compatible(Voice{Engine: "moss"}, "en") || q.Compatible(Voice{Engine: id}, "ar") {
			t.Fatal("incompatible state or language accepted", id)
		}
	}
}
