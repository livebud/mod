package mod_test

import (
	"os"
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
	is.Equal(wd, module.Directory())
	is.Equal(module.Import(), "github.com/livebud/mod")
}

func TestMustFind(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	module := mod.MustFind(wd)
	is.Equal(wd, module.Directory())
	is.Equal(module.Import(), "github.com/livebud/mod")
}

func TestNew(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	module := mod.New(dir)
	is.Equal(dir, module.Directory())
	is.Equal(module.Import(), "change.me")
}
