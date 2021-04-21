package mustache

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

var Debug = true

const (
	Strings               = '*'
	TAG_Variable          = '$'
	TAG_Variable_UnEscape = '&'
	TAG_Section           = '#'
	TAG_Inverted_Section  = '^'
	TAG_Comment           = '!'
	TAG_Partial           = '>'
	TAG_End               = '/'
)

func (tpl *Template) AddSegment(segment Segment) error {
	log.Printf("Add New Segment: Type=%s, Value=%s", string(segment.Type), segment.Value)
	if KeepSegment {
		tpl.Smts = append(tpl.Smts, segment)
	}

	if segment.Type == TAG_Comment {
		return nil
	}

	if segment.Type == TAG_End {
		if tpl.cur.Name() != segment.Value {
			return parseError{segment.LineNumber, fmt.Sprintf("CloseTag expected [%s] but [%s]", tpl.cur.Name(), segment.Value)}
		}
		tpl.cur = tpl.cur.Father()
		log.Println("End Tag --> " + segment.Value)
		return nil
	}

	node := makeRenderNode(segment, tpl.cur)
	tpl.cur.AddChildren(node)
	if segment.Type == Strings || segment.Type == TAG_Variable || segment.Type == TAG_Variable_UnEscape {
		log.Println("Non-Children Node --> " + segment.Value)
		return nil
	}
	log.Println("New Node with Children -->" + segment.Value)
	tpl.cur = node
	return nil
}

func ParseReader(r io.Reader) (*Template, error) {
	root := &TopRenderNode{}
	tpl := &Template{make([]Segment, 0), root, root}

	rd := bufio.NewReader(r)
	lineNumber := -1

	end := false
	for {
		lineNumber++
		line, err := rd.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return nil, parseError{lineNumber, err.Error()}
			} else {
				end = true
			}
		}

		err = parseLine(line, lineNumber, tpl)
		if err != nil {
			return nil, parseError{lineNumber, err.Error()}
		}
		if end {
			break
		}
	}
	return tpl, nil
}

func parseLine(line string, lineNumber int, tpl *Template) error {
	sz := len(line)
	start := 0
	end := 0
	started := false
	for end < sz {
		//log.Println("start=", start, "end=", end, "line=", lineNumber)
		started = false

		//Search Start
		for end < sz {
			//log.Println("Search TAG,  start=", start, "end=", end)
			if line[end] == '{' && (end+3) < sz && line[end+1] == '{' {
				started = true
				if end > start {
					tpl.AddSegment(Segment{Strings, line[start:end], lineNumber})
				}
				start = end + 2
				end = start + 1
				//log.Println("Tag Start Found, beark")
				break
			} else {
				started = false
				if line[end] == '\\' {
					if (end + 1) < sz {
						end++
					}
				}
				end++
			}
		}

		if !started {
			if end > start {
				tpl.AddSegment(Segment{Strings, line[start:end], lineNumber})
			}
			break
		}

		//log.Println("Search Tag end", "start=", start, "end=", end)

		tagValue := ""

		for (end + 1) < sz {
			if line[end] != '}' || line[end+1] != '}' {
				//log.Println("Current=" + line[end:end+2])
				end++
				continue
			}
			escape := true
			// {{{ABC}}} --> start=2, end=
			if line[start] == '{' && (end+2) < sz && line[end+2] == '}' {
				tagValue = line[start+1 : end]
				start = end + 3
				end = start
				escape = false
			} else {
				tagValue = line[start:end]
				start = end + 2
				end = start
			}

			//log.Println("Tag init Value=" + tagValue)

			tagValue = strings.TrimRight(tagValue, " \t")
			if tagValue == "" {
				return parseError{lineNumber, "Blank Tag"}
			}

			switch tagValue[0] {
			case TAG_Section:
				fallthrough
			case TAG_Inverted_Section:
				fallthrough
			case TAG_Comment:
				fallthrough
			case TAG_Partial:
				fallthrough
			case TAG_Variable_UnEscape:
				fallthrough
			case TAG_End:
				if len(tagValue) == 1 {
					return parseError{lineNumber, "Invaild Tag"}
				}

				typeRune := rune(tagValue[0])
				//log.Println("TAG Type" + string(typeRune))
				tagValue = strings.Trim(tagValue[1:], " \t")
				if tagValue == "" {
					return parseError{lineNumber, "Emtry Tag"}
				}
				err := tpl.AddSegment(Segment{typeRune, tagValue, lineNumber})
				if err != nil {
					return err
				}
			default:
				//log.Println("First Char = " + tagValue[0:1])
				tagValue = strings.Trim(tagValue, " \t")
				if !escape {
					tpl.AddSegment(Segment{TAG_Variable_UnEscape, tagValue, lineNumber})
				} else {
					tpl.AddSegment(Segment{TAG_Variable, tagValue, lineNumber})
				}
			}
			break
		}
	}
	return nil
}

type parseError struct {
	LineNumber int
	Message    string
}

func (p parseError) Error() string {
	return "Error at Line" + strconv.Itoa(p.LineNumber) + " " + p.Message
}

func (s *Segment) String() string {
	if s.Type == Strings {
		return s.Value
	}
	if s.Type == TAG_Variable {
		return fmt.Sprintf("{{%s}}", s.Value)
	}
	return fmt.Sprintf("{{%s %s}}", string(s.Type), s.Value)
}

func DumpSegments(w io.Writer, segs []Segment) {
	for _, s := range segs {
		io.WriteString(w, s.String())
	}
}

//--------------------------------------------------------------
