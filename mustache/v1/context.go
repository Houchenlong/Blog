package mustache

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	//"strings"
)

type context struct {
	dir string
}

type ComboContext struct {
	*context
	Ctxs []Context
}

func (me *ComboContext) Get(key string) (reflect.Value, bool) {
	for _, ctx := range me.Ctxs {
		val, found := ctx.Get(key)
		if found {
			return val, true
		}
	}
	return reflect.Zero(reflect.TypeOf(0)), false
}

func (c *ComboContext) Root() Context {
	return c.root
}

func MakeContext(objs ...interface{}) Context {
	ctxs := make([]Context, len(objs))
	for i, obj := range objs {
		ctxs[i] = _makeContext(obj)
	}
	return &ComboContext{ctxs, nil}
}

func _makeContext(obj interface{}) Context {
	t := reflect.TypeOf(obj)
	//log.Println(t.String())
	if t.String() == "reflect.Value" {
		return &BasicContext{obj.(reflect.Value), nil}
	}
	return &BasicContext{reflect.ValueOf(obj), nil}
}

type BasicContext struct {
	value reflect.Value
	root  Context
}

func (ctx *BasicContext) Get(key string) (reflect.Value, bool) {
	val, found := ctx._get(key)
	if !found {
		return val, found
	}
	if val.Kind() == reflect.Interface {
		log.Println("XXX", val.Type(), val.Interface())
		return reflect.ValueOf(val.Interface()), true
	}
	return val, found
}

func (ctx *BasicContext) _get(key string) (reflect.Value, bool) {
	if !ctx.value.IsValid() {
		return ctx.value, false
	}
	if key == "." {
		return ctx.value, true
	}
	log.Println("ctx value kind=", ctx.value.Kind().String(), "key=", key)
	switch ctx.value.Kind() {
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
		return reflect.Zero(ctx.value.Type()), false
	case reflect.Map:
		for _, _key := range ctx.value.MapKeys() {
			_v := ctx.value.MapIndex(_key)
			if key == fmt.Sprintf("%v", _key.Interface()) {
				log.Println("Map Found key=", key, _v.Kind())
				return _v, true
			}
		}
		log.Println("Key Not found--?>" + key)
		return reflect.Zero(ctx.value.Type()), false
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		index, err := strconv.ParseInt(key, 0, 32)
		if err != nil {
			return reflect.Zero(ctx.value.Type()), false
		}
		if index < 0 || int(index) >= ctx.value.Len() {
			return reflect.Zero(ctx.value.Type()), false
		}
		return ctx.value.Index(int(index)), true
	case reflect.Ptr:
		if ctx.value.Elem().Kind() != reflect.Struct {
			return reflect.Zero(ctx.value.Type()), false
		}
		fallthrough
	case reflect.Struct:
		log.Println(ctx.value.Type())
		v := reflect.Indirect(ctx.value)
		field := v.FieldByName(key)
		if field.IsValid() {
			log.Println("Return Field Value for key=" + key)
			return field, true
		}

		t := ctx.value.Type()
		method, ok := t.MethodByName(key)
		if !ok {
			log.Println("Miss field or method for -->" + key)
			return reflect.Zero(ctx.value.Type()), false
		}
		if method.Func.Type().NumIn() != 1 || method.Func.Type().NumOut() == 0 {
			log.Println("Named Func Found , but has-args or void return", method, method.Func.Type().NumIn(), method.Func.Type().NumOut())
			return reflect.Zero(ctx.value.Type()), false
		}
		return method.Func.Call([]reflect.Value{ctx.value})[0], true
	default:
		log.Println("Not Support kind=" + ctx.value.Kind().String())
		return reflect.Zero(ctx.value.Type()), false
	}
	panic("Impossible")
	return reflect.Zero(reflect.TypeOf("")), false
}

func (ctx *BasicContext) Root() Context {
	return ctx.root
}

func AsBool(value reflect.Value) bool {
	if !value.IsValid() {
		return false
	}
	switch value.Kind() {
	case reflect.Int:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		return value.Int() != 0
	case reflect.Uint:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		return value.Uint() != 0
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		return value.Float() != 0.0
	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		return value.Complex() != complex128(0)
	case reflect.Bool:
		return value.Bool()
	case reflect.String:
		fallthrough
	case reflect.Map:
		fallthrough
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		return value.Len() != 0
	case reflect.Ptr:
		if value.Elem().Kind() != reflect.Struct {
			return false
		}
		fallthrough
	case reflect.Struct:
		return true
	case reflect.Func:
		return true
	}
	log.Println("Not Support kind=" + value.Kind().String())
	return false
}
