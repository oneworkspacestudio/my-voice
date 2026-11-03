package main

import (
	"github.com/oneworkspacestudio/my-voice/internal/app"
	"github.com/oneworkspacestudio/my-voice/internal/webui"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
	"net"
	"os"
	"path/filepath"
	"strings"
)

var version = "0.1.0-preview.2"

func main() {
	// Portable launchers still keep private state outside replaceable programs.
	hasData, hasListen := false, false
	for _, s := range os.Args[1:] {
		if s == "--data" || strings.HasPrefix(s, "--data=") {
			hasData = true
		}
		if s == "--listen" || strings.HasPrefix(s, "--listen=") {
			hasListen = true
		}
	}
	if !hasData {
		os.Args = append(os.Args, "--data", filepath.Join(os.Getenv("LOCALAPPDATA"), "OneWorkspace", "MyVoice"))
	}
	listen := "127.0.0.1:9038"
	if !hasListen {
		if probe, e := net.Listen("tcp4", listen); e == nil {
			probe.Close()
		} else if probe, e := net.Listen("tcp4", "127.0.0.1:0"); e == nil {
			listen = probe.Addr().String()
			probe.Close()
		}
	}
	err := apphost.Main(apphost.Config{ID: "voice", Name: "My Voice", Version: version, Listen: listen, MediaCapture: true, RestoreFiles: []string{"blobs", "models", "projects"}}, webui.Files(), func(h *apphost.Host) (apphost.Application, error) { return app.New(h) })
	if err != nil {
		apphost.Report(err)
		os.Exit(1)
	}
}
