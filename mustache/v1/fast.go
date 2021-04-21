package mustache

import (
	"bytes"
)

var Fast = &fast{}

type fast struct {
}

func (f fast) RenderString(str string, objs ...interface{}) (string, error) {
	w := &bytes.Buffer{}
	err := RenderString(str, w, objs...)
	return string(w.String()), err
}

func (f fast) RenderFile(path string, objs ...interface{}) (string, error) {
	w := &bytes.Buffer{}
	err := RenderFile(path, w, objs...)
	return string(w.String()), err
}
