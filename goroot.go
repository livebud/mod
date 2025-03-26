package mod

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var goRoot string

// This function is a heavily modified version of the following:
// https://github.com/golang/go/blob/89044b6d423a07bea3b6f80210f780e859dd2700/src/cmd/go/internal/cfg/cfg.go#L369
func findGoRoot() (string, error) {
	if goRoot != "" {
		return goRoot, nil
	}
	if env := os.Getenv("GOROOT"); env != "" {
		goRoot = filepath.Clean(env)
		return goRoot, nil
	}
	if r := runtime.GOROOT(); r != "" {
		goRoot = filepath.Clean(r)
		return goRoot, nil
	}
	if runtime.Compiler == "gccgo" {
		// gccgo has no real GOROOT, and it certainly doesn't
		// depend on the executable's location.
		return "", errors.New("gccgo has no GOROOT")
	}
	// Run `go env GOROOT` to compute the GOROOT.
	// TODO: this takes about 8ms on boot, see if we can find a better way to
	// compute the correct GOROOT.
	// Note: The tricky part is making it work with --trimpath.
	cmd := exec.Command("go", "env", "GOROOT")
	stdout := new(bytes.Buffer)
	cmd.Stdout = stdout
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("mod: unable to find GOROOT: %w", err)
	}
	goRoot = strings.TrimSpace(stdout.String())
	return goRoot, nil
}
