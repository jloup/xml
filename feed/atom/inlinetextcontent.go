package atom

import (
	"encoding/xml"

	"github.com/JLoup/xml/utils"
)

type InlineTextContent struct {
	Content string

	depth  utils.DepthWatcher
	Parent utils.Visitor
}

func NewInlineTextContent() *InlineTextContent {
	d := utils.NewDepthWatcher()

	return &InlineTextContent{depth: d}
}

func (i *InlineTextContent) ProcessStartElement(el utils.StartElement) (utils.Visitor, utils.ParserError) {
	i.depth.Down()
	return i, utils.NewError(LeafElementHasChild, "inline text shoud not contain XML childs")
}

func (i *InlineTextContent) ProcessEndElement(el xml.EndElement) (utils.Visitor, utils.ParserError) {
	if i.depth.Up() == utils.ParentLevel {
		return i.Parent.ProcessEndElement(el)
	}

	return i, nil
}

func (i *InlineTextContent) ProcessCharData(el xml.CharData) (utils.Visitor, utils.ParserError) {
	i.Content = string(el)
	return i, nil
}

func (i *InlineTextContent) String() string {
	return i.Content
}
