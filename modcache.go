package mod

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"

	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"
)

// Cache for faster subsequent requests
var modCacheDir string

// getModDir returns the module cache directory
func getModCacheDir() string {
	if modCacheDir != "" {
		return modCacheDir
	}
	env := os.Getenv("GOMODCACHE")
	if env != "" {
		modCacheDir = env
		return env
	}
	modCacheDir = filepath.Join(build.Default.GOPATH, "pkg", "mod")
	return modCacheDir
}

// getModuleDirectory returns an absolute path to the required module.
func getModuleDirectory(modulePath, version string) (string, error) {
	enc, err := module.EscapePath(modulePath)
	if err != nil {
		return "", err
	}
	if !semver.IsValid(version) {
		return "", fmt.Errorf("non-semver module version %q", version)
	}
	if module.CanonicalVersion(version) != version {
		return "", fmt.Errorf("non-canonical module version %q", version)
	}
	encVer, err := module.EscapeVersion(version)
	if err != nil {
		return "", err
	}
	dir := filepath.Join(getModCacheDir(), enc+"@"+encVer)
	return dir, nil
}
