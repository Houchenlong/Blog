package mustache

import (
	"fmt"
	"io"
	"strings"
)

type RenderError struct {
	lineNumber int
	msg        string
}

func (r *RenderError) Error() string {
	return fmt.Sprintf("Render FAIL lineNumber=%d, msg=%s", r.lineNumber, r.msg)
}

type RenderNode interface {
	Render(io.Writer, Context) error
	Name() string
	AddChildren(RenderNode)
	Father() RenderNode
	LineNumber() int
}

type TopRenderNode struct {
	Clildren []RenderNode
}

func (t *TopRenderNode) Name() string {
	return ""
}

func (t *TopRenderNode) Father() RenderNode {
	return nil
}

func (t *TopRenderNode) Render(w io.Writer, ctx Context) error {
	for _, node := range t.Clildren {
		err := node.Render(w, ctx)
		if err != nil {
			return err
		}
		//ctx.Reset()
	}
	return nil
}

func (t *TopRenderNode) AddChildren(node RenderNode) {
	if t.Clildren == nil {
		t.Clildren = make([]RenderNode, 0)
	}
	t.Clildren = append(t.Clildren, node)
}

func (t *TopRenderNode) LineNumber() int {
	return 0
}

func makeRenderNode(segment Segment, father RenderNode) RenderNode {
	switch segment.Type {
	case Strings:
		return &StringsRenderNode{strings.Replace(segment.Value, "\\{{", "{{", -1), segment.LineNumber}
	case TAG_Variable:
		return &VariableRenderNode{segment.Value, true, segment.LineNumber}
	case TAG_Variable_UnEscape:
		return &VariableRenderNode{segment.Value, false, segment.LineNumber}
	case TAG_Section:
		return &SectionRenderNode{segment.Value, false, segment.LineNumber, father, make([]RenderNode, 0)}
	case TAG_Inverted_Section:
		return &SectionRenderNode{segment.Value, true, segment.LineNumber, father, make([]RenderNode, 0)}
	case TAG_Partial:
		return &PartialRenderNode{segment.Value, segment.LineNumber, father, make([]RenderNode, 0)}
	}

	panic("Impossible")
}

type StringsRenderNode struct {
	Value      string
	lineNumber int
}

func (s *StringsRenderNode) Name() string {
	return s.Value
}

func (s *StringsRenderNode) Father() RenderNode {
	return nil
}

func (s *StringsRenderNode) Render(w io.Writer, ctx Context) error {
	_, err := io.WriteString(w, s.Value)
	return err
}

func (s *StringsRenderNode) AddChildren(child RenderNode) {
}

func (s *StringsRenderNode) LineNumber() int {
	return s.lineNumber
}
