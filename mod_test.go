package mod_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/livebud/mod"
	"github.com/matryer/is"
)

func TestFind(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	module, err := mod.Find()
	is.NoErr(err)
	is.Equal(wd, module.Dir())
	is.Equal(module.Import(), "github.com/livebud/mod")
}

func TestMustFind(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	module := mod.MustFind(wd)
	is.Equal(wd, module.Dir())
	is.Equal(module.Import(), "github.com/livebud/mod")
}

func TestNew(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	module := mod.New(dir)
	is.Equal(dir, module.Dir())
	is.Equal(module.Import(), "change.me")
}

func TestInferOutsideGoPath(t *testing.T) {
	dir := t.TempDir()
	is := is.New(t)
	is.Equal(mod.Infer(dir), "")
}

func TestInferInsideGoPath(t *testing.T) {
	gopath := mod.GOPATH
	mod.GOPATH = t.TempDir()
	t.Cleanup(func() {
		mod.GOPATH = gopath
	})
	is := is.New(t)
	is.Equal(mod.Infer(filepath.Join(mod.GOPATH, "src", "app.com")), "app.com")
}
