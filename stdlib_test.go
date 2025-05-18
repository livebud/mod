package mod_test

import (
	"testing"

	"github.com/livebud/mod"
	"github.com/matryer/is"
)

func TestInStdlib(t *testing.T) {
	is := is.New(t)
	is.True(mod.InStdlib("net/http"))
	is.True(!mod.InStdlib("github.com/livebud/mod"))
	is.True(!mod.InStdlib("github.com/livebud/mod/hello"))
	is.True(mod.InStdlib("internal/reflectlite"))
}
