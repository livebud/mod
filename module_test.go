package mod_test

import (
	"io/fs"
	"os"
	"strings"
	"testing"

	"github.com/livebud/mod"
	"github.com/matryer/is"
)

func TestReadFile(t *testing.T) {
	is := is.New(t)
	module, err := mod.Find(".")
	is.NoErr(err)
	gomod, err := fs.ReadFile(module, "go.mod")
	is.NoErr(err)
	is.True(strings.Contains(string(gomod), "module github.com/livebud/mod"))
}

func TestReadDir(t *testing.T) {
	is := is.New(t)
	module, err := mod.Find(".")
	is.NoErr(err)
	des, err := fs.ReadDir(module, ".")
	is.NoErr(err)
	is.True(len(des) > 0)
}

func TestContains(t *testing.T) {
	is := is.New(t)
	module, err := mod.Find(".")
	is.NoErr(err)
	is.True(module.Contains(module.Import()))
	is.True(module.Contains(module.Import("hello", "world")))
	is.True(!module.Contains("net/http"))
	is.True(!module.Contains("github.com/livebud/bud"))
}

func TestResolveDirRel(t *testing.T) {
	is := is.New(t)
	module, err := mod.Find(".")
	is.NoErr(err)
	dir, err := module.ResolveDir(module.Import())
	is.NoErr(err)
	is.Equal(dir, module.Dir())
}

func TestResolveDirStd(t *testing.T) {
	is := is.New(t)
	module, err := mod.Find(".")
	is.NoErr(err)
	dir, err := module.ResolveDir("net/http")
	is.NoErr(err)
	// Check the the dir exists
	_, err = os.Stat(dir)
	is.NoErr(err)
}

func TestResolveDirModCache(t *testing.T) {
	is := is.New(t)
	module, err := mod.Find(".")
	is.NoErr(err)
	dir, err := module.ResolveDir("golang.org/x/mod")
	is.NoErr(err)
	is.True(strings.Contains(dir, "golang.org/x/mod@"))
	// Check the the dir exists
	_, err = os.Stat(dir)
	is.NoErr(err)
}

func TestResolveImport(t *testing.T) {
	is := is.New(t)
	module, err := mod.Find(".")
	is.NoErr(err)
	im, err := module.ResolveImport(module.Dir())
	is.NoErr(err)
	is.Equal(module.Import(), im)
}
