package atom

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/JLoup/errors"
	"github.com/JLoup/xml/feed/extension"
	"github.com/JLoup/xml/helper"
)

type BasicElement struct {
	CommonAttributes
	Content helper.Element
	name    xml.Name

	Extension extension.VisitorExtension
	depth     helper.DepthWatcher
	Parent    helper.Visitor
}

func NewBasicElement(parent helper.Visitor) *BasicElement {
	b := BasicElement{Parent: parent, depth: helper.NewDepthWatcher()}
	b.depth.SetMaxDepth(1)

	b.InitCommonAttributes()

	return &b
}

func NewBasicElementExt(parent helper.Visitor, manager extension.Manager) *BasicElement {
	b := NewBasicElement(parent)

	b.Extension = extension.InitExtension("basicelement", manager)

	return b
}

func (b *BasicElement) SetParent(parent helper.Visitor) {
	b.Parent = parent
}

func (b *BasicElement) Name() xml.Name {
	return b.name
}

func (b *BasicElement) ProcessStartElement(el helper.StartElement) (helper.Visitor, helper.ParserError) {
	if b.depth.IsRoot() {
		b.name = el.Name
		b.Extension = extension.InitExtension(b.name.Local, b.Extension.Manager)

		for _, attr := range el.Attr {
			if !b.ProcessAttr(attr) {
				b.Extension.ProcessAttr(attr, b)
			}
		}
	}

	if b.depth.Down() == helper.MaxDepthReached {
		return b, helper.NewError(LeafElementHasChild, fmt.Sprintf("'%s' shoud not have childs", b.name.Local))
	}

	return b, nil
}

func (b *BasicElement) ProcessEndElement(el xml.EndElement) (helper.Visitor, helper.ParserError) {
	if b.depth.Up() == helper.RootLevel {
		return b.Parent, b.Validate()
	}

	return b, nil
}

func (b *BasicElement) ProcessCharData(el xml.CharData) (helper.Visitor, helper.ParserError) {
	b.Content.Value = strings.TrimSpace(string(el))
	return b, nil
}

func (b *BasicElement) Validate() helper.ParserError {
	error := errors.NewErrorAggregator()

	b.Extension.Validate(&error)
	b.ValidateCommonAttributes(b.name.Local, &error)

	if err := b.Content.Validate(); err != nil {
		error.NewError(helper.NewError(err.Flag(), fmt.Sprintf("%s's %s", b.name.Local, err.Msg())))
	}

	return error.ErrorObject()
}

func (b *BasicElement) String() string {
	return b.Content.Value
}

func (b *BasicElement) Reset() {
	b.depth.Reset()
}
