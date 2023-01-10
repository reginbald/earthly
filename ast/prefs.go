//go:generate hel

package ast

import (
	"os"

	"github.com/pkg/errors"
)

type prefs struct {
	reader          NamedReader
	done            func()
	enableSourceMap bool
}

// Opt is an option function for customizing the behavior of parsing functions.
type Opt func(prefs) (prefs, error)

// WithSourceMap tells ParseVersion to enable a source map when parsing.
func WithSourceMap() Opt {
	return func(p prefs) (prefs, error) {
		p.enableSourceMap = true
		return p, nil
	}
}

// FromOpt is an option function for customizing the source reader of parsing
// functions.
type FromOpt func(prefs) (prefs, error)

// FromPath tells parsing functions to open and read from a file at path.
func FromPath(path string) FromOpt {
	return func(p prefs) (prefs, error) {
		f, err := os.Open(path)
		if err != nil {
			return p, errors.Wrapf(err, "ast: unable to open file '%v'", path)
		}
		p.reader = f
		p.done = func() { f.Close() }
		return p, nil
	}
}

// NamedReader is simply an io.Reader with a name. It matches os.File, but
// allows non-file types to be passed in.
type NamedReader interface {
	Name() string
	Seek(offset int64, whence int) (int64, error)
	Read(buff []byte) (n int, err error)
}

// FromReader tells ParseVersion to read from r.
func FromReader(r NamedReader) FromOpt {
	return func(p prefs) (prefs, error) {
		p.reader = r
		return p, nil
	}
}
