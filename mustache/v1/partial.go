package mustache

import (
	"io"
)

type PartialRenderNode struct {
	Key        string
	lineNumber int
	father     RenderNode
	Clildren   []RenderNode
}

func (node *PartialRenderNode) Name() string {
	return node.Key
}

func (node *PartialRenderNode) Father() RenderNode {
	return node.father
}

func (node *PartialRenderNode) Render(w io.Writer, ctx Context) error {
	return &RenderError{node.lineNumber, "Not support yet"}
}

func (node *PartialRenderNode) AddChildren(child RenderNode) {
	node.Clildren = append(node.Clildren, child)
}

func (node *PartialRenderNode) LineNumber() int {
	return node.lineNumber
}
