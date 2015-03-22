package atom

import (
	"encoding/xml"

	"github.com/JLoup/xml/helper"
)

type InlineTextContent struct {
	Content string

	depth  helper.DepthWatcher
	Parent helper.Visitor
}

func NewInlineTextContent() *InlineTextContent {
	d := helper.NewDepthWatcher()

	return &InlineTextContent{depth: d}
}

func (i *InlineTextContent) ProcessStartElement(el helper.StartElement) (helper.Visitor, helper.ParserError) {
	i.depth.Down()
	return i, helper.NewError(LeafElementHasChild, "inline text shoud not contain XML childs")
}

func (i *InlineTextContent) ProcessEndElement(el xml.EndElement) (helper.Visitor, helper.ParserError) {
	if i.depth.Up() == helper.ParentLevel {
		return i.Parent.ProcessEndElement(el)
	}

	return i, nil
}

func (i *InlineTextContent) ProcessCharData(el xml.CharData) (helper.Visitor, helper.ParserError) {
	i.Content = string(el)
	return i, nil
}

func (i *InlineTextContent) String() string {
	return i.Content
}
