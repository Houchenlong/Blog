package mustache

import (
	"bytes"
	"io"
	"os"
)

func RenderString(str string, objs ...interface{}) (string, error) {
	w := bytes.NewBuffer(nil)
	err := RenderReader(bytes.NewBufferString(str), w, objs...)
	if err != nil {
		return "", err
	}
	return w.String(), nil
}

func RenderFile(path string, objs ...interface{}) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	w := bytes.NewBuffer(nil)
	err = RenderReader(f, w, objs...)
	if err != nil {
		return "", err
	}
	return w.String(), nil
}

func RenderReader(r io.Reader, w io.Writer, objs ...interface{}) error {
	tpl, err := Parse(r)
	if err != nil {
		return err
	}
	ctx := MakeContexts(objs...)
	err = tpl.Render(ctx, w)
	if err != nil {
		return err
	}
	return nil
}
