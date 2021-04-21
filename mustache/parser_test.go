package mustache

import (
	"bytes"
	_ "encoding/json"
	"log"
	"reflect"
	"testing"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	//KeepSegment = true
}

func TestSimpleParse(t *testing.T) {
	segs := _parseSegments("{{Name}}")
	if len(segs) != 1 {
		t.Log(segs)
		t.FailNow()
	}
	segs = _parseSegments("{{A}}H         H{{BB}}")
	segs = _parseSegments("{{#CC}}{{/CC}}")
	segs = _parseSegments("{{{EE}}}{{&FF}}")
	segs = _parseSegments("{{A}}\n\n\n\n\n\n\n\n{{UUU}}\n\n\n\\{{}}\n")
	_ = segs
}

func _parseSegments(str string) []Node {
	//log.Println("-------------------> " + str)
	tpl, err := Parse(bytes.NewBufferString(str))
	if err != nil {
		log.Fatal(err)
	}
	return tpl.Tree
}

func printSegments(ss []Node) {
	for _, s := range ss {
		log.Println(s.String())
	}
}

func TestPrintTemplate(t *testing.T) {
	//KeepSegment = false
	tpl, err := Parse(bytes.NewBufferString("This {{#Nutz}}\nHi,{{Name}}!\n{{/Nutz}}"))
	if err != nil {
		log.Fatal(err)
	}
	if len(tpl.Tree) != 2 {
		t.FailNow()
	}
	log.Println(tpl.Tree[1].(*SectionNode).Name())
	obj := tpl.Tree[1].(*SectionNode).Clildren[2]
	log.Println(reflect.TypeOf(obj))
	log.Println("TPL = \n" + tpl.String())
}

func TestMuiltLineTemplate(t *testing.T) {
	//KeepSegment = false
	tpl, err := Parse(bytes.NewBufferString("{{name}}\n{{age}}\n"))
	if err != nil {
		log.Fatal(err)
	}
	if len(tpl.Tree) != 4 {
		t.FailNow()
	}
	log.Println("TPL = \n" + tpl.String())
	str, err := RenderString("{{name}}\n{{age}} ABC\n", map[string]interface{}{"name": "wendal", "age": 27})
	if err != nil {
		log.Fatal(err)
	}
	if str != "wendal\n27 ABC\n" {
		log.Fatal(str)
	}
}
