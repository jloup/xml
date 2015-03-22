package atom

import (
	"encoding/xml"

	"github.com/JLoup/errors"
	"github.com/JLoup/xml/feed/extension"
	"github.com/JLoup/xml/helper"
)

type Icon struct {
	CommonAttributes
	Iri helper.Element

	Extension extension.VisitorExtension
	Parent    helper.Visitor
	depth     helper.DepthWatcher
}

func NewIcon() *Icon {
	i := Icon{depth: helper.NewDepthWatcher()}

	i.Iri = helper.NewElement("iri", "", IsValidIRI)

	i.InitCommonAttributes()
	i.depth.SetMaxDepth(1)

	return &i
}

func NewIconExt(manager extension.Manager) *Icon {
	i := NewIcon()

	i.Extension = extension.InitExtension("icon", manager)
	return i
}

func (i *Icon) ProcessStartElement(el helper.StartElement) (helper.Visitor, helper.ParserError) {
	if i.depth.IsRoot() {
		i.ResetAttr()
		for _, attr := range el.Attr {
			if !i.ProcessAttr(attr) {
				i.Extension.ProcessAttr(attr, i)
			}
		}
	}

	if i.depth.Down() == helper.MaxDepthReached {
		return i, helper.NewError(LeafElementHasChild, "icon element should not have childs")
	}

	return i, nil
}

func (i *Icon) ProcessEndElement(el xml.EndElement) (helper.Visitor, helper.ParserError) {
	if i.depth.Up() == helper.RootLevel {
		return i.Parent, i.validate()
	}

	return i, nil
}

func (i *Icon) ProcessCharData(el xml.CharData) (helper.Visitor, helper.ParserError) {
	i.Iri.Value = string(el)
	return i, nil
}

func (i *Icon) validate() helper.ParserError {
	error := errors.NewErrorAggregator()

	helper.ValidateElement("icon", i.Iri, &error)
	i.Extension.Validate(&error)
	i.ValidateCommonAttributes("icon", &error)

	return error.ErrorObject()
}
