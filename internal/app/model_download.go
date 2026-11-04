package app

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func verifiedModelFile(path string, a Asset) bool {
	f, e := os.Open(path)
	if e != nil {
		return false
	}
	defer f.Close()
	info, e := f.Stat()
	if e != nil || info.Size() != a.Size {
		return false
	}
	h := sha256.New()
	n, e := io.Copy(h, io.LimitReader(f, a.Size+1))
	return e == nil && n == a.Size && hex.EncodeToString(h.Sum(nil)) == a.SHA
}

type modelProgress struct {
	w    io.Writer
	done int64
	last time.Time
	show func(int64)
}

func (p *modelProgress) Write(b []byte) (int, error) {
	n, e := p.w.Write(b)
	p.done += int64(n)
	if time.Since(p.last) >= time.Second {
		p.show(p.done)
		p.last = time.Now()
	}
	return n, e
}

// Pinned URLs are immutable. Incomplete bytes can be resumed, but no file is
// installed until its entire length and SHA-256 match the registry.
func downloadModelAsset(ctx context.Context, client *http.Client, a Asset, dest string, show func(int64)) error {
	part := dest + ".part"
	var offset int64
	if f, e := os.Stat(part); e == nil {
		offset = f.Size()
	}
	if offset == a.Size && verifiedModelFile(part, a) {
		show(a.Size)
		return os.Rename(part, dest)
	}
	if offset >= a.Size {
		offset = 0
	}
	req, e := http.NewRequestWithContext(ctx, "GET", a.URL, nil)
	if e != nil {
		return e
	}
	if offset > 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", offset))
	}
	resp, e := client.Do(req)
	if e != nil {
		return errors.New("Download interrupted. Retry to continue the saved download")
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		offset = 0 // This host does not resume; safely restart the partial file.
	} else if resp.StatusCode == http.StatusPartialContent && offset > 0 {
		want := fmt.Sprintf("bytes %d-%d/%d", offset, a.Size-1, a.Size)
		if strings.TrimSpace(resp.Header.Get("Content-Range")) != want {
			return errors.New("The download server returned an invalid range. Retry later")
		}
	} else {
		return fmt.Errorf("Model download returned HTTP %d. Retry later", resp.StatusCode)
	}
	if s := resp.Header.Get("Content-Length"); s != "" {
		if n, e := strconv.ParseInt(s, 10, 64); e != nil || n != a.Size-offset {
			return errors.New("Model download length does not match the pinned file")
		}
	}
	flags := os.O_CREATE | os.O_WRONLY
	if offset == 0 {
		flags |= os.O_TRUNC
	} else {
		flags |= os.O_APPEND
	}
	f, e := os.OpenFile(part, flags, 0600)
	if e != nil {
		return e
	}
	p := &modelProgress{w: f, done: offset, show: show}
	show(offset)
	n, copyErr := io.Copy(p, io.LimitReader(resp.Body, a.Size-offset+1))
	closeErr := f.Close()
	if copyErr != nil || n != a.Size-offset {
		return errors.New("Download interrupted. Retry to continue the saved download")
	}
	if closeErr != nil {
		return closeErr
	}
	if !verifiedModelFile(part, a) {
		_ = os.Remove(part)
		return errors.New("Download checksum failed. Retry to download a clean copy")
	}
	show(a.Size)
	return os.Rename(part, dest)
}
