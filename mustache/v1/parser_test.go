package mustache

import (
	"bytes"
	"encoding/json"
	"log"
	"testing"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	KeepSegment = true
}

func TestSimpleParse(*testing.T) {
	segs := _parseSegments("{{Name}}")
	if len(segs) != 1 || segs[0].Value != "Name" {
		log.Fatal("Error")
	}

	segs = _parseSegments("{{A}}H         H{{BB}}")
	if len(segs) != 3 || segs[0].Value != "A" || segs[1].Value != "H         H" || segs[2].Value != "BB" {
		log.Fatal("Error")
	}

	segs = _parseSegments("{{#CC}}{{/CC}}")
	if len(segs) != 2 || segs[0].Type != TAG_Section || segs[0].Value != "CC" || segs[1].Type != TAG_End || segs[1].Value != "CC" {
		printSegments(segs)
		log.Println(segs[1].Type != TAG_End)
		log.Fatal("Error")
	}
	segs = _parseSegments("{{{EE}}}{{&FF}}")
	if len(segs) != 2 || segs[0].Type != TAG_Variable_UnEscape || segs[1].Type != TAG_Variable_UnEscape {
		printSegments(segs)
		log.Fatal("Error")
	}

	segs = _parseSegments("{{A}}\n\n\n\n\n\n\n\n{{UUU}}\n\n\n\\{{}}\n")
	if segs[len(segs)-1].Type != Strings || segs[len(segs)-1].Value != "\\{{}}\n" || segs[len(segs)-1].LineNumber < 5 {
		log.Fatal("Error")
	}
}

func _parseSegments(str string) []Segment {
	tpl, err := ParseReader(bytes.NewBufferString(str))
	if err != nil {
		log.Fatal(err)
	}
	return tpl.Smts
}

func printSegments(ss []Segment) {
	for _, s := range ss {
		log.Println(s.String())
	}
}

func TestPrintTemplate(*testing.T) {
	KeepSegment = false
	tpl, err := ParseReader(bytes.NewBufferString("\t{{#Nutz}}\nHi,{{Name}}!\n{{/Nutz}}"))
	if err != nil {
		log.Fatal(err)
	}
	d, err := json.Marshal(tpl)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("TPL = \n" + string(d))
}
