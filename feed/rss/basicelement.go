package rss

import (
	"github.com/jloup/utils"
	"github.com/jloup/xml/feed/extension"
	xmlutils "github.com/jloup/xml/utils"
	"encoding/xml"
	"fmt"
	"strings"
)

type BasicElement struct {
	Content xmlutils.Element
	name    xml.Name

	Extension extension.VisitorExtension
	depth     xmlutils.DepthWatcher
	Parent    xmlutils.Visitor
}

func NewBasicElement() *BasicElement {
	d := xmlutils.NewDepthWatcher()
	d.SetMaxDepth(1)

	return &BasicElement{depth: d, Content: xmlutils.NewElement("", "", xmlutils.Nop)}
}

func NewBasicElementExt(manager extension.Manager) *BasicElement {
	b := NewBasicElement()

	b.Extension = extension.InitExtension("basicelement", manager)

	return b
}

func (b *BasicElement) SetParent(parent xmlutils.Visitor) {
	b.Parent = parent
}

func (b *BasicElement) Name() xml.Name {
	return b.name
}

func (b *BasicElement) ProcessStartElement(el xmlutils.StartElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if b.depth.IsRoot() {
		b.name = el.Name
		b.Extension = extension.InitExtension(b.name.Local, b.Extension.Manager)

		for _, attr := range el.Attr {
			b.Extension.ProcessAttr(attr, b)
		}
	}

	if b.depth.Down() == xmlutils.MaxDepthReached {
		return b, xmlutils.NewError(LeafElementHasChild, fmt.Sprintf("'%s' shoud not have childs", b.Name))
	}

	return b, nil
}

func (b *BasicElement) ProcessEndElement(el xml.EndElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if b.depth.Up() == xmlutils.RootLevel {
		return b.Parent, b.Validate()
	}

	return b, nil
}

func (b *BasicElement) ProcessCharData(el xml.CharData) (xmlutils.Visitor, xmlutils.ParserError) {
	b.Content.Value = strings.TrimSpace(string(el))
	return b, nil
}

func (b *BasicElement) Validate() xmlutils.ParserError {
	error := utils.NewErrorAggregator()

	b.Extension.Validate(&error)

	if err := b.Content.Validate(); err != nil {
		error.NewError(xmlutils.NewError(err.Flag(), fmt.Sprintf("%s's %s", b.name.Local, err.Msg())))
	}

	return error.ErrorObject()
}

func (b *BasicElement) String() string {
	return b.Content.Value
}

func (b *BasicElement) Reset() {
	b.depth.Reset()
}
