package mustache

import (
	"bytes"
	"log"
	"reflect"
	"testing"
)

func TestVariableRender(*testing.T) {
	v := &VariableRenderNode{"Name", false, 0}
	ctx := MakeContext(map[string]interface{}{"Name": "wendal", "Url": "<a>http://wendal.net</a>"})

	w := &bytes.Buffer{}

	err := v.Render(w, ctx)
	if err != nil {
		log.Fatal(err)
	}

	str := string(w.Bytes())

	if str != "wendal" {
		log.Fatal("Not match --> " + str)
	} else {
		log.Println(str)
	}

	v = &VariableRenderNode{"Url", false, 0}
	w = &bytes.Buffer{}
	v.Render(w, ctx)
	str = string(w.Bytes())
	if str != "<a>http://wendal.net</a>" {
		log.Fatal("Not Match --> " + str)
	}

	v = &VariableRenderNode{"Url", true, 0}
	w = &bytes.Buffer{}
	v.Render(w, ctx)
	str = string(w.Bytes())

	if str == "<a>http://wendal.net</a>" {
		log.Fatal("Not OK --> " + str)
	}
}

type TestStruct struct {
	Name  string
	Age   int
	url   string
	Count int
}

func (t *TestStruct) Url() string {
	t.Count++
	return t.url
}

func TestSimpleStructRender(*testing.T) {
	t := TestStruct{"wendal", 27, "http://wendal.net", 0}
	ctx := MakeContext(t)
	w := &bytes.Buffer{}

	v := &VariableRenderNode{"Name", false, 0}
	v2 := &VariableRenderNode{"Age", false, 0}
	v3 := &VariableRenderNode{"Url", false, 0}

	err := v.Render(w, ctx)
	if err != nil {
		log.Fatal(err)
	}
	str := string(w.Bytes())
	if str != "wendal" {
		log.Fatal("NOT Match -->        " + str)
	}

	w = &bytes.Buffer{}
	err = v2.Render(w, ctx)
	if err != nil {
		log.Fatal(err)
	}
	str = string(w.Bytes())
	if str != "27" {
		log.Fatal("NOT Match -->" + str)
	}

	//-----------------------------------

	w = &bytes.Buffer{}
	ctx = MakeContext(&t)
	err = v3.Render(w, ctx)
	if err != nil {
		log.Fatal(err)
	}
	str = string(w.Bytes())
	if str != "http://wendal.net" {
		log.Fatal("NOT Match -->" + str)
	}

	if t.Count != 1 {
		log.Fatal("Count != 1 ?!!")
	}
}

func TestPtrRender(*testing.T) {
	log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	t := &TestStruct{"wendal", 27, "http://wendal.net", 0}
	ctx := MakeContext(t)
	w := &bytes.Buffer{}

	log.Println(t, ctx, w)

	log.Println(reflect.ValueOf(t).Elem().Interface().(TestStruct).Name)
}
