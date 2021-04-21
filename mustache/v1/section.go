package mustache

import (
	"io"
	"log"
	"reflect"
)

type SectionRenderNode struct {
	Key        string
	Inverted   bool
	lineNumber int
	father     RenderNode
	Clildren   []RenderNode
}

func (node *SectionRenderNode) Name() string {
	return node.Key
}

func (node *SectionRenderNode) Father() RenderNode {
	return node.father
}

func (node *SectionRenderNode) Render(w io.Writer, ctx Context) error {
	v, found := ctx.Get(node.Key)
	kind := v.Type().Kind()
	log.Println("Section Type Kind = " + kind.String())
	ok := false
	var _ctx Context
	if found {
		ok = AsBool(v)
	}

	if ok && node.Inverted {
		log.Println("Section=true and Inverted, beark")
		return nil
	}
	if !ok && !node.Inverted {
		log.Println("Section=false, beark")
		return nil
	}

	if kind == reflect.Array || kind == reflect.Slice {
		log.Println("Section Arry/Slice")
		for i := 0; i < v.Len(); i++ {
			_ctx = _makeContext(v.Index(i))
			for _, child := range node.Clildren {
				err := child.Render(w, _ctx)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
	if !v.IsValid() {
		_ctx = _makeContext(reflect.ValueOf(""))
	} else {
		if v.Kind() == reflect.Func {
			return v.Interface().(func(Context, []RenderNode, io.Writer) error)(ctx, node.Clildren, w)
		}
		_ctx = _makeContext(v)
	}
	for _, child := range node.Clildren {
		err := child.Render(w, _ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (node *SectionRenderNode) AddChildren(clild RenderNode) {
	node.Clildren = append(node.Clildren, clild)
}

func (node *SectionRenderNode) LineNumber() int {
	return node.lineNumber
}

//type SectionSegmentRender func(Context, []RenderNode, io.Writer) error
