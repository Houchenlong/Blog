package mustache

import (
	"io"
	"log"
	"reflect"
	"runtime/debug"
	"strings"
	_g_tpl "text/template"
)

type ValNode struct {
	name   string
	Escape bool
}

func (node *ValNode) Render(ctx Context, w io.Writer) error {
	val, found := ctx.Get(node.name)
	if !found {
		return nil
		//return errors.New("key=" + node.name + " NOT Found")
	}
	str := val.String()
	//log.Println("ValNode --> " + str)
	if node.Escape {
		str = _g_tpl.HTMLEscapeString(str)
	}
	_, err := w.Write([]byte(str))
	return err
}

func (node *ValNode) String() string {
	if node.Escape {
		return "{{&" + node.name + "}}"
	}
	return "{{" + node.name + "}}"
}

func (node *ValNode) Name() string {
	return node.name
}

//------------------------------------------------------
type SectionNode struct {
	name     string
	Inverted bool
	Clildren []Node
}

func (node *SectionNode) Render(ctx Context, w io.Writer) error {
	//log.Println(">> SectionNode : " + node.Name())
	//defer log.Println(">> End SectionNode : " + node.Name())
	//TODO person?to_name
	ctx_helper_name := ""
	key := node.name
	if strings.Contains(node.name, "?") {
		_tmp := strings.Split(node.name, "?")
		key = _tmp[0]
		ctx_helper_name = _tmp[1]
	}
	var err error

	val, found := ctx.Get(key)

	if found {
		if ctx_helper_name != "" {
			//log.Println("Search Ctx Helper", ctx_helper_name)
			ctxHelper, found := ctx.Get(ctx_helper_name)
			if !found {
				log.Println("NO Ctx Helper", ctx_helper_name)

				return nil
			}
			_helper, ok := ctxHelper.Val.Interface().(func(interface{}) interface{})
			if !ok {
				log.Println("NO GOOD Ctx Helper", ctxHelper)
				return nil
			}
			val = &Value{reflect.ValueOf(_helper(val.Val.Interface()))}
			//log.Println("Done for Ctx Helper")
		}
		f, ok := val.Val.Interface().(SectionRenderFunc)
		if ok {
			//log.Println("Using BaiscHelper", key)
			return f(node.Clildren, node.Inverted, ctx, w)
		} else {
			if val.Val.Type().Kind() == reflect.Func {
				log.Println("What?", val.Val.Interface())
			}
		}
	}

	isTrue := found && val.Bool()

	if isTrue && node.Inverted {
		return nil
	}
	if !isTrue && !node.Inverted {
		return nil
	}

	if node.Inverted && !(isTrue) {
		//if found {
		//	log.Println(">>", val.Val.Interface())
		//}

		//log.Println(">> Inverted SectionNode : ", node.name, key, ctx_helper_name)
		//defer log.Println(">> End Inverted SectionNode : " + node.Name())
		for _, _node := range node.Clildren {
			err = _node.Render(ctx, w)
			if err != nil {
				return err
			}
		}
		return nil
	}

	//log.Println("Section True, render children!", node.name)

	switch val.Val.Type().Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		for i := 0; i < val.Val.Len(); i++ {
			_ctx := MakeContexts(val.Val.Index(i), ctx)
			for _, child := range node.Clildren {
				err := child.Render(_ctx, w)
				if err != nil {
					return err
				}
			}
		}
	default:
		//log.Println("using default Section render", key, ctx_helper_name)
		_ctx := MakeContexts(val.Val, ctx)
		for _, child := range node.Clildren {
			err := child.Render(_ctx, w)
			if err != nil {
				return err
			}
		}
	}

	//log.Println("End Section ", node.name)

	return nil
}

func (node *SectionNode) String() string {
	str := "{{"
	if node.Inverted {
		str += "^"
	} else {
		str += "#"
	}
	str += node.name + "}}"
	for _, c := range node.Clildren {
		str += c.String()
	}
	str += "{{/" + node.name + "}}"
	return str
}

func (node *SectionNode) Name() string {
	return node.name
}

//------------------------------------------------------
type PartialNode struct {
	name string
}

func (node *PartialNode) Render(ctx Context, w io.Writer) error {
	dir := ctx.Dir()
	if dir == "" {
		log.Println(ctx)
		dir = "."
	}
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	p := dir + node.name
	str, err := RenderFile(p, ctx)
	if err != nil {
		log.Println(string(debug.Stack()))
		return err
	}
	_, err = w.Write([]byte(str))

	return err
}

func (node *PartialNode) String() string {
	return "{{>" + node.name + "}}"
}

func (node *PartialNode) Name() string {
	return node.name
}

//------------------------------------------------------

type ConstantNode struct {
	Val string
}

func (node *ConstantNode) Render(ctx Context, w io.Writer) error {
	_, err := w.Write([]byte(node.Val))
	return err
}

func (node *ConstantNode) String() string {
	return node.Val
}

func (node *ConstantNode) Name() string {
	return node.Val
}

type SectionRenderFunc func([]Node, bool, Context, io.Writer) error
