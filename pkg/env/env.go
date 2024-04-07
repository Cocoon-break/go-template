package env

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var (
	Built   = "0"
	Version = "0.0.0"
	GitHash = "no git hash set"
	App     = "go-template"
)

var (
	e *RuntimeEnv
	m sync.Mutex
)

type RuntimeEnv struct {
	AppName   string `json:"app_name"`
	Built     string `json:"built"`
	Version   string `json:"version"`
	GoVersion string `json:"go_version"`
	GitHash   string `json:"git_hash"`
	OSArch    string `json:"os_arch"`
}

func CompileInfo() *RuntimeEnv {
	m.Lock()
	defer m.Unlock()

	if e == nil {
		e = &RuntimeEnv{
			AppName:   App,
			Built:     Built,
			Version:   Version,
			GoVersion: runtime.Version(),
			GitHash:   GitHash,
			OSArch:    runtime.GOOS + "-" + runtime.GOARCH,
		}
	}
	return e
}

func (e *RuntimeEnv) Print() {
	t, _ := strconv.ParseInt(e.Built, 0, 64)
	fmt.Printf(`%s (release)
  Version: %s
  Go version: %s
  Git commit: %s
  OS/Arch: %s
  Built: %s
`, App, e.Version, e.GoVersion, e.GitHash, e.OSArch, time.Unix(t, 0))
}

func AppName() string {
	return App
}
