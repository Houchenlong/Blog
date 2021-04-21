package mustache

import (
	"bytes"
	//"errors"
	"io"
	"os"
)

var (
	KeepSegment      = true
	IgnoreInvaildKey = true
)

type Template struct {
	Smts []Segment
	Root RenderNode
	cur  RenderNode
}

type InvokedVaule struct {
}

type Segment struct {
	Type       rune
	Value      string
	LineNumber int
}

func RenderReader(r io.Reader, w io.Writer, ctx ...interface{}) error {
	t, err := ParseReader(r)
	if err != nil {
		return err
	}
	return t.Render(w, ctx...)
}

//func ParseReader(r io.Reader) (*Template, error)

func RenderString(str string, w io.Writer, ctx ...interface{}) error {
	t, err := ParseString(str)
	if err != nil {
		return err
	}
	return t.Render(w, ctx...)
}

func RenderFile(path string, w io.Writer, ctx ...interface{}) error {
	t, err := ParseFile(path)
	if err != nil {
		return err
	}
	return t.Render(w, ctx...)
}

func ParseString(str string) (*Template, error) {
	return ParseReader(bytes.NewBufferString(str))
}

func ParseFile(path string) (*Template, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseReader(f)
}

func (t *Template) Render(w io.Writer, objs ...interface{}) error {
	context := MakeContext(objs...)
	err := t.Root.Render(w, context)
	if err != nil {
		return err
	}
	return nil
}
