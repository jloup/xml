package atom

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/jloup/utils"
	"github.com/jloup/xml/feed/extension"
	xmlutils "github.com/jloup/xml/utils"
)

type BasicElement struct {
	CommonAttributes
	Content xmlutils.Element
	name    xml.Name

	Extension extension.VisitorExtension
	depth     xmlutils.DepthWatcher
	Parent    xmlutils.Visitor
}

func NewBasicElement(parent xmlutils.Visitor) *BasicElement {
	b := BasicElement{Parent: parent, depth: xmlutils.NewDepthWatcher()}
	b.depth.SetMaxDepth(1)

	b.InitCommonAttributes()

	return &b
}

func NewBasicElementExt(parent xmlutils.Visitor, manager extension.Manager) *BasicElement {
	b := NewBasicElement(parent)

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
			if !b.ProcessAttr(attr) {
				b.Extension.ProcessAttr(attr, b)
			}
		}
	}

	if b.depth.Down() == xmlutils.MaxDepthReached {
		return b, xmlutils.NewError(LeafElementHasChild, fmt.Sprintf("'%s' shoud not have childs", b.name.Local))
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
	b.ValidateCommonAttributes(b.name.Local, &error)

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
