package atom

import (
	"encoding/xml"

	"github.com/JLoup/errors"
	"github.com/JLoup/xml/feed/extension"
	"github.com/JLoup/xml/helper"
)

type Logo struct {
	CommonAttributes
	Iri helper.Element

	Extension extension.VisitorExtension
	Parent    helper.Visitor
	depth     helper.DepthWatcher
}

func NewLogo() *Logo {
	l := Logo{depth: helper.NewDepthWatcher()}

	l.Iri = helper.NewElement("iri", "", IsValidIRI)

	l.InitCommonAttributes()
	l.depth.SetMaxDepth(1)

	return &l
}

func NewLogoExt(manager extension.Manager) *Logo {
	l := NewLogo()
	l.Extension = extension.InitExtension("logo", manager)

	return l
}

func (l *Logo) ProcessStartElement(el helper.StartElement) (helper.Visitor, helper.ParserError) {
	if l.depth.IsRoot() {
		l.ResetAttr()
		for _, attr := range el.Attr {
			if !l.ProcessAttr(attr) {
				l.Extension.ProcessAttr(attr, l)
			}
		}
	}

	if l.depth.Down() == helper.MaxDepthReached {
		return l, helper.NewError(LeafElementHasChild, "logo element should not have childs")
	}

	return l, nil
}

func (l *Logo) ProcessEndElement(el xml.EndElement) (helper.Visitor, helper.ParserError) {
	if l.depth.Up() == helper.RootLevel {
		return l.Parent, l.validate()
	}

	return l, nil
}

func (l *Logo) ProcessCharData(el xml.CharData) (helper.Visitor, helper.ParserError) {
	l.Iri.Value = string(el)
	return l, nil
}

func (l *Logo) validate() helper.ParserError {
	error := errors.NewErrorAggregator()

	helper.ValidateElement("logo", l.Iri, &error)
	l.ValidateCommonAttributes("logo", &error)
	l.Extension.Validate(&error)

	return error.ErrorObject()
}
