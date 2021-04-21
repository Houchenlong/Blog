package mustache

import (
	"bytes"
	"log"
	"testing"
)

func TestSimpleRender(*testing.T) {
	m := map[string]map[string]string{}
	m2 := map[string]string{}
	m2["media"] = "wendal.net"
	m["urls"] = m2
	str, err := RenderString(">>{{urls.media}}<<", m)
	if err != nil {
		log.Fatal(err)
	}
	if str != ">>wendal.net<<" {
		log.Fatal(str)
	}
	log.Println(str)
}

func TestSimpleRender2(*testing.T) {
	m := map[string]map[string]string{}
	m2 := map[string]string{}
	m2["media"] = "<a></a>"
	m["urls"] = m2
	str, err := RenderString(">>{{{urls.media}}}<<", m)
	if err != nil {
		log.Fatal(err)
	}
	if str != ">><a></a><<" {
		log.Fatal(str)
	}
	log.Println(str)
}

func TestSimpleRender3(*testing.T) {
	m := map[string]string{"name": "wendal"}
	str, err := RenderString("{{name}}A\nB{{   name}}C\tD{{name}}E", m)
	if err != nil {
		log.Fatal(err)
	}
	if str != "wendalA\nBwendalC\tDwendalE" {
		log.Fatal(str)
	}
	log.Println(str)
}

func TestSection(*testing.T) {
	m := map[string]string{"name": "wendal"}
	str, err := RenderString("{{# name}}{{.}}{{/name}}{{^name}}ABC{{/name}}", m)
	if err != nil {
		log.Fatal(err)
	}
	if str != "wendal" {
		tpl, _ := Parse(bytes.NewBufferString("{{# name}}{{.}}{{/name}}{{^name}}ABC{{/name}}"))
		log.Println(tpl)
		log.Println(len(tpl.Tree))
		log.Fatal(str)
	}
	log.Println(str)
}
