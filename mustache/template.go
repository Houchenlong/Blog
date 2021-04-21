package mustache

import (
	"bytes"
	"fmt"
	"io"
	//"log"
	"reflect"
)

type Template struct {
	Tree []Node
}

func (tpl *Template) Render(ctx Context, w io.Writer) error {
	for _, node := range tpl.Tree {
		err := node.Render(ctx, w)
		if err != nil {
			return err
		}
	}
	return nil
}

func (tpl *Template) String() string {
	w := bytes.NewBuffer(nil)
	for _, node := range tpl.Tree {
		//log.Println("!>" + node.Name())
		w.WriteString(node.String())
	}
	return w.String()
}

type Node interface {
	Name() string
	Render(Context, io.Writer) error
	String() string
	//Clone() *Node
}

type Context interface {
	Get(string) (*Value, bool)
	Dir() string
}

type Value struct {
	Val reflect.Value
}

func (v *Value) String() string {
	return fmt.Sprintf("%v", v.Val.Interface())
}

func (v *Value) Bool() bool {
	if !v.Val.IsValid() {
		return false
	}
	switch v.Val.Type().Kind() {
	case reflect.Array:
		return v.Val.Len() != 0
	case reflect.Map:
		return v.Val.Len() != 0
	case reflect.Slice:
		return v.Val.Len() != 0
	case reflect.String:
		return len(v.Val.String()) != 0
	case reflect.Bool:
		return v.Val.Bool()
	case reflect.Struct:
		return true
	case reflect.Ptr:
		return (&Value{v.Val.Elem()}).Bool()
	case reflect.Func:
		return true
	}
	return false
}
