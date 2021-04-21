package mustache

import (
	"bytes"
	"io"
	"log"
	"reflect"
	"testing"
)

func TestReflect(*testing.T) {
	log.Println("---------------------------------------")
	log.Println(reflect.TypeOf(uint64(1)).Kind().String())
	log.Println(reflect.TypeOf(uint64(1)).Name())

	log.Println(reflect.TypeOf(map[string]interface{}{}).Kind().String())

	log.Println(reflect.TypeOf(make([]Segment, 0)).Kind().String())

	log.Println(reflect.TypeOf(Segment{}).Kind().String())
	log.Println(reflect.TypeOf(&Segment{}).Kind().String())

	log.Println(">>>>>>>>>>>", reflect.ValueOf(Segment{}).CanAddr())

}

func TestSectionRender(*testing.T) {
	log.Println("---------------------------------------")
	section := &SectionRenderNode{"admin", false, 1, nil, []RenderNode{&VariableRenderNode{"name", true, 2}}}
	m := map[string]interface{}{}
	m2 := map[string]string{"name": "wendal"}
	m["admin"] = m2
	ctx := MakeContext(m)
	w := &bytes.Buffer{}
	err := section.Render(w, ctx)
	if err != nil {
		log.Fatal(err)
	}
	str := w.String()
	if str != "wendal" {
		log.Fatal("Render ERROR --> " + str)
	}

	log.Println("Simple Session Test Success")
}

func TestSectionArraySlice(*testing.T) {
	m := map[string]interface{}{}
	m["list_str"] = []string{"abc", "EFG"}
	m["list_int"] = []int{1, 2, 3, 4, 5}
	m["is_ok"] = true
	m["is_fail"] = false

	ctx := MakeContext(m)
	w := &bytes.Buffer{}

	section := &SectionRenderNode{"list_str", false, 1, nil, []RenderNode{&VariableRenderNode{".", true, 2}}}
	err := section.Render(w, ctx)
	if err != nil {
		log.Fatal(err)
	}
	str := string(w.Bytes())
	if str != "abcEFG" {
		log.Fatal(">>" + str)
	}
	w.Reset()

	section = &SectionRenderNode{"list_int", false, 1, nil, []RenderNode{&VariableRenderNode{".", true, 2}}}
	err = section.Render(w, ctx)
	if err != nil {
		log.Fatal(err)
	}
	str = string(w.Bytes())
	if str != "12345" {
		log.Fatal(">>" + str)
	}

	w.Reset()
	tpl, err := ParseString("{{#list_int}}A{{.}}B{{/list_int}}")
	log.Println(m)
	tpl.Render(w, m)
	str = string(w.Bytes())
	if str != "A1BA2BA3BA4BA5B" {
		log.Println(">>>>" + str)
	}

	w.Reset()
	str, err = Fast.RenderString("{{#is_ok}}AB{{/is_ok}}", m)
	if str != "AB" {
		log.Fatal("E>> " + str)
	}

	str, err = Fast.RenderString("V{{^is_ok}}AB{{/is_ok}}Z", m)
	if str != "VZ" || err != nil {
		log.Fatal("E>> " + str)
	}
}

func _test_section_node_render(ctx Context, nodes []RenderNode, w io.Writer) error {
	io.WriteString(w, "Wendal")
	return nil
}

func TestSectionRenderFunc(*testing.T) {
	m := map[string]interface{}{"func": _test_section_node_render}
	str, err := Fast.RenderString("{{#func}}AB{{/func}}", m)
	if err != nil {
		log.Fatal(err)
	}
	if str != "Wendal" {
		log.Fatal("Not expect : " + str)
	}
}
