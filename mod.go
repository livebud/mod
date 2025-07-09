package mod

import (
	"errors"
	"fmt"
	"go/build"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
)

// ErrFileNotFound occurs when no go.mod can be found
var ErrFileNotFound = fmt.Errorf(`mod: unable to find "go.mod": %w`, fs.ErrNotExist)

// New module
func New(dir string) *Module {
	modulePath := Infer(dir)
	if modulePath == "" {
		modulePath = "change.me"
	}
	module, err := Parse(filepath.Join(dir, "go.mod"), []byte(`module `+modulePath))
	if err != nil {
		panic("mod: invalid module data: " + err.Error())
	}
	return module
}

// Find the first go.mod file in one of the directories below or return an
// error. Find will also search parent directories for a go.mod file.
func Find(dirs ...string) (*Module, error) {
	if len(dirs) == 0 {
		return find(".")
	}
	for _, dir := range dirs {
		module, err := find(dir)
		if err != nil {
			if !errors.Is(err, ErrFileNotFound) {
				return nil, err
			}
			continue
		}
		return module, nil
	}
	return nil, ErrFileNotFound
}

func find(dir string) (*Module, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	modPath, err := lookup(absDir)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(modPath)
	if err != nil {
		return nil, fmt.Errorf(`mod: unable to read "go.mod": %w`, err)
	}
	return Parse(modPath, data)
}

// MustFind a go.mod file in this directory or any parent directory. If no
// go.mod file is found, this will panic.
func MustFind(dirs ...string) *Module {
	module, err := Find(dirs...)
	if err != nil {
		panic(err)
	}
	return module
}

// Parse a go.mod file
func Parse(path string, data []byte) (*Module, error) {
	modfile, err := modfile.Parse(path, data, nil)
	if err != nil {
		return nil, err
	}
	if modfile.Module == nil {
		modFile, err := modfile.Format()
		if err != nil {
			return nil, fmt.Errorf("mod: missing module statement in %q and got an error while formatting %s", path, err)
		}
		return nil, fmt.Errorf("mod: missing module statement in %q, received %q", path, string(modFile))
	}
	dir := filepath.Dir(path)
	return &Module{dir, modfile, os.DirFS(dir)}, nil
}

// Lookup finds the absolute path of the go.mod file in the given directory
func Lookup(dir string) (path string, err error) {
	path, err = lookup(dir)
	if err != nil {
		return "", err
	}
	return filepath.Abs(path)
}

func lookup(dir string) (path string, err error) {
	path = filepath.Join(dir, "go.mod")
	// Check if this path exists, otherwise recursively traverse towards root
	if _, err = os.Stat(path); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return "", err
		}
		nextDir := filepath.Dir(dir)
		if nextDir == dir {
			return "", ErrFileNotFound
		}
		return lookup(nextDir)
	}
	return filepath.EvalSymlinks(path)
}

// Abs finds the absolute directory that contains the go.mod file
func Abs(dir string) (absDir string, err error) {
	path, err := lookup(dir)
	if err != nil {
		return "", err
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return filepath.Dir(absPath), nil
}

// GOPATH is the default GOPATH for the current build environment. This is
// exposed for testing purposes, but should not be changed in production code.
var GOPATH = build.Default.GOPATH

// Infer the module path from the $GOPATH. This only works if you work inside
// $GOPATH otherwise returns an empty string.
func Infer(dir string) string {
	src := filepath.Join(GOPATH, "src") + "/"
	if !strings.HasPrefix(dir, src) {
		return ""
	}
	modulePath := strings.TrimPrefix(dir, src)
	return modulePath
}
