package mod

import (
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
)

type Module struct {
	dir  string
	file *modfile.File
	fs.FS
}

// Dir returns the absolute directory to the module.
func (m *Module) Dir(subpaths ...string) string {
	return filepath.Join(append([]string{m.dir}, subpaths...)...)
}

// Import returns the base import path of the module.
func (m *Module) Import(subpaths ...string) string {
	modulePath := m.file.Module.Mod.Path
	subPath := path.Join(subpaths...)
	if modulePath == "std" {
		return subPath
	}
	return path.Join(modulePath, subPath)
}

// Contains checks if the module contains the import path.
func (m *Module) Contains(importPath string) bool {
	return contains(m.Import(), importPath)
}

// ResolveImport returns an import path from a local directory.
func (m *Module) ResolveImport(dir string) (importPath string, err error) {
	return m.resolveImport(dir, true)
}

func (m *Module) resolveImport(dir string, evalSymlinks bool) (string, error) {
	relPath, err := filepath.Rel(m.dir, dir)
	if err != nil {
		return "", err
	} else if strings.HasPrefix(relPath, "..") {
		if !evalSymlinks {
			return "", fmt.Errorf("module: unable to resolve import. %q can't be outside the module directory %q", dir, m.dir)
		}
		// Maybe the directory is a symlink, resolve that symlink and try again.
		if dir, err = filepath.EvalSymlinks(dir); err != nil {
			return "", fmt.Errorf("module: unable to resolve import for %q. %w", dir, err)
		}
		return m.resolveImport(dir, false)
	}
	return m.Import(relPath), nil
}

// ResolveDir resolves an import path to an absolute path.
func (m *Module) ResolveDir(importPath string) (dir string, err error) {
	// Handle standard library
	if InStdlib(importPath) {
		return m.resolveStdlib(importPath)
	}

	// Handle local packages
	modulePath := m.Import()
	if contains(modulePath, importPath) {
		// Ensure the resolved relative dir exists
		rel, err := filepath.Rel(modulePath, importPath)
		if err != nil {
			return "", err
		}
		return filepath.Join(m.dir, rel), nil
	}

	// Handle replaces
	for _, rep := range m.file.Replace {
		if contains(rep.Old.Path, importPath) {
			relPath := strings.TrimPrefix(importPath, rep.Old.Path)
			newPath := filepath.Join(rep.New.Path, relPath)
			absdir, err := resolvePath(m.dir, newPath)
			if err != nil {
				return "", err
			}
			return absdir, nil
		}
	}

	// Handle requires
	for _, req := range m.file.Require {
		if contains(req.Mod.Path, importPath) {
			relPath := strings.TrimPrefix(importPath, req.Mod.Path)
			dir, err := getModuleDir(req.Mod.Path, req.Mod.Version)
			if err != nil {
				return "", err
			}
			return filepath.Join(dir, relPath), nil
		}
	}

	// Lastly, check if the import path is a vendored package within stdlib
	vendoredPath := path.Join("vendor", importPath)
	if InStdlib(vendoredPath) {
		return m.resolveStdlib(vendoredPath)
	}

	return "", fmt.Errorf("mod: unable to resolve directory for import path %q. %w", importPath, fs.ErrNotExist)
}

func (m *Module) resolveStdlib(importPath string) (dir string, err error) {
	goRoot, err := findGoRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(goRoot, "src", importPath), nil
}

// Contains checks if `importPath` is a subpath of `basePath`.
func contains(basePath, importPath string) bool {
	return basePath == importPath || strings.HasPrefix(importPath, basePath+"/")
}

// Resolve allows `path` to be replaced by an absolute path in `rest`
func resolvePath(path string, rest ...string) (string, error) {
	result := path
	for _, p := range rest {
		if filepath.IsAbs(p) {
			result = p
			continue
		}
		result = filepath.Join(result, p)
	}
	return filepath.Abs(result)
}
