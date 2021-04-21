package mustache

import (
	//"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

var VALUE_NIL = &Value{reflect.Zero(reflect.TypeOf(0))}
var NIL = reflect.Zero(reflect.TypeOf(0))

//---------------------------------------------------------------

type ComboContext struct {
	Ctxs []Context
	dir  string
}

func (me *ComboContext) Get(key string) (*Value, bool) {
	for _, ctx := range me.Ctxs {
		val, found := ctx.Get(key)
		if found {
			return val, true
		}
	}
	return nil, false
}

func (me *ComboContext) Dir() string {
	dir := me.dir
	if dir == "" {
		for _, ctx := range me.Ctxs {
			dir = ctx.Dir()
			if dir != "" {
				break
			}
		}
	}
	if dir == "" {
		log.Println("Not Ctx Dir Found?! Return Emtry", me)
	}
	return dir
}

//-----------------------------------------------------

func MakeContexts(objs ...interface{}) Context {
	ctxs := make([]Context, len(objs))
	for i, obj := range objs {
		//comboCtx, ok := obj.(ComboContext)
		//if ok {
		//	for _, ctx := range comboCtx.Ctxs {
		//		ctxs = append(ctxs, ctx)
		//	}
		//	continue
		//}
		//ctxs = append(ctxs, MakeContext(obj))
		ctxs[i] = MakeContext(obj)
	}
	return &ComboContext{ctxs, ""}
}

func MakeContext(obj interface{}) Context {
	val, ok := obj.(Context)
	if ok {
		return val
	}
	t := reflect.TypeOf(obj)
	//log.Println(t.String())
	if t.String() == "reflect.Value" {
		return &BasicContext{obj.(reflect.Value), ""}
	}
	return &BasicContext{reflect.ValueOf(obj), ""}
}

func MakeContextDir(obj interface{}, dir string) Context {
	ctx := MakeContext(obj)
	ctx.(*BasicContext).dir = dir
	return ctx
}

//-------------------------------------

type BasicContext struct {
	value reflect.Value
	dir   string
}

func (ctx *BasicContext) Get(key string) (*Value, bool) {
	key = strings.Trim(key, "\t\n ")
	if !ctx.value.IsValid() {
		return nil, false
	}
	switch {
	case key == "":
		return nil, false
	case key == ".":
		return &Value{ctx.value}, true
	case !strings.Contains(key, "."):
		val, found := _Get(ctx.value, key)
		if found {
			return &Value{val}, true
		}
		return nil, false
	}

	tmp := ctx.value
	var _found bool
	for _, _key := range strings.Split(key, ".") {
		tmp, _found = _Get(tmp, _key)
		if !_found {
			return nil, false
		}
	}

	return &Value{tmp}, true
}

func _Get(val reflect.Value, key string) (reflect.Value, bool) {
	if !val.IsValid() {
		return val, false
	}
	_val := val.Interface()
	var _rs interface{}
	switch _val.(type) {
	case map[string]interface{}:
		_rs = _val.(map[string]interface{})[key]
	case int:
		return NIL, false
	case int64:
		return NIL, false
	case uint:
		return NIL, false
	default:
		mer, ok := _val.(MapGet)
		if ok {
			_rs = mer.Get(key)
		}
	}
	if _rs != nil {
		return reflect.ValueOf(_rs), true
	}
	//if key == "." {
	//	return val, true
	//}
	//log.Println("ctx value kind=", val.Kind().String(), "key=", key)
	switch val.Kind() {
	case reflect.Int:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Uint:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		fallthrough
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		fallthrough
	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		fallthrough
	case reflect.Bool:
		fallthrough
	case reflect.String:
		return reflect.Zero(val.Type()), false
	case reflect.Map:
		for _, _key := range val.MapKeys() {
			mKey := _key.Interface()
			switch mKey.(type) {
			case string:
				if key == mKey.(string) {
					return val.MapIndex(_key), true
				}
			}
		}
		//log.Println("Key Not found--?>" + key)
		return reflect.Zero(val.Type()), false
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		index, err := strconv.ParseInt(key, 0, 32)
		if err != nil {
			return reflect.Zero(val.Type()), false
		}
		if index < 0 || int(index) >= val.Len() {
			return reflect.Zero(val.Type()), false
		}
		return val.Index(int(index)), true
	case reflect.Ptr:
		if val.Elem().Kind() != reflect.Struct {
			return reflect.Zero(val.Type()), false
		}
		fallthrough
	case reflect.Struct:
		key = strings.Title(key)
		//log.Println(val.Type())
		v := reflect.Indirect(val)
		field := v.FieldByName(key)
		if field.IsValid() {
			//log.Println("Return Field Value for key=" + key)
			return field, true
		}

		t := val.Type()
		method, ok := t.MethodByName(key)
		if !ok {
			//log.Println("Miss field or method for -->" + key)
			return reflect.Zero(val.Type()), false
		}
		if method.Func.Type().NumIn() != 1 || method.Func.Type().NumOut() == 0 {
			log.Println("Named Func Found , but has-args or void return", method, method.Func.Type().NumIn(), method.Func.Type().NumOut())
			return reflect.Zero(val.Type()), false
		}
		return method.Func.Call([]reflect.Value{val})[0], true
	default:
		log.Println("Not Support kind=" + val.Kind().String())
		return reflect.Zero(val.Type()), false
	}
	panic("Impossible")
	return reflect.Zero(reflect.TypeOf("")), false
}

func (ctx *BasicContext) Dir() string {
	return ctx.dir
}

type MapGet interface {
	Get(string) interface{}
}
