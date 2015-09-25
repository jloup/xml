package atom

import (
	"encoding/xml"

	xmlutils "github.com/jloup/xml/utils"
)

type InlineTextContent struct {
	Content string

	depth  xmlutils.DepthWatcher
	Parent xmlutils.Visitor
}

func NewInlineTextContent() *InlineTextContent {
	d := xmlutils.NewDepthWatcher()

	return &InlineTextContent{depth: d}
}

func (i *InlineTextContent) ProcessStartElement(el xmlutils.StartElement) (xmlutils.Visitor, xmlutils.ParserError) {
	i.depth.Down()
	return i, xmlutils.NewError(LeafElementHasChild, "inline text shoud not contain XML childs")
}

func (i *InlineTextContent) ProcessEndElement(el xml.EndElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if i.depth.Up() == xmlutils.ParentLevel {
		return i.Parent.ProcessEndElement(el)
	}

	return i, nil
}

func (i *InlineTextContent) ProcessCharData(el xml.CharData) (xmlutils.Visitor, xmlutils.ParserError) {
	i.Content = string(el)
	return i, nil
}

func (i *InlineTextContent) String() string {
	return i.Content
}
