package atom

import (
	"encoding/xml"
	"strings"

	"github.com/JLoup/errors"
	"github.com/JLoup/xml/feed/extension"
	"github.com/JLoup/xml/helper"
)

type Id struct {
	CommonAttributes
	Content helper.Element

	Extension extension.VisitorExtension
	Parent    helper.Visitor
	depth     helper.DepthWatcher
}

func NewId() *Id {
	i := Id{depth: helper.NewDepthWatcher()}

	i.Content = helper.NewElement("iri", "", IsAbsoluteIRI)

	i.InitCommonAttributes()
	i.depth.SetMaxDepth(1)

	return &i
}

func NewIdExt(manager extension.Manager) *Id {
	i := NewId()

	i.Extension = extension.InitExtension("id", manager)

	return i
}
func (i *Id) ProcessStartElement(el helper.StartElement) (helper.Visitor, helper.ParserError) {
	if i.depth.IsRoot() {
		i.ResetAttr()
		for _, attr := range el.Attr {
			if !i.ProcessAttr(attr) {
				i.Extension.ProcessAttr(attr, i)
			}
		}
	}

	if i.depth.Down() == helper.MaxDepthReached {
		return i, helper.NewError(LeafElementHasChild, "id element should not have childs")
	}

	return i, nil
}

func (i *Id) ProcessEndElement(el xml.EndElement) (helper.Visitor, helper.ParserError) {
	if i.depth.Up() == helper.RootLevel {
		return i.Parent, i.validate()
	}

	return i, nil
}

func (i *Id) ProcessCharData(el xml.CharData) (helper.Visitor, helper.ParserError) {
	i.Content.Value = strings.TrimSpace(string(el))
	return i, nil
}

func (i *Id) validate() helper.ParserError {
	error := errors.NewErrorAggregator()

	helper.ValidateElement("id", i.Content, &error)
	i.Extension.Validate(&error)
	i.ValidateCommonAttributes("id", &error)

	return error.ErrorObject()
}

func (i *Id) String() string {
	return i.Content.String()
}
