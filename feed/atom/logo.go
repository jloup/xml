package atom

import (
	"encoding/xml"

	"github.com/jloup/errors"
	"github.com/jloup/xml/feed/extension"
	"github.com/jloup/xml/utils"
)

type Logo struct {
	CommonAttributes
	Iri utils.Element

	Extension extension.VisitorExtension
	Parent    utils.Visitor
	depth     utils.DepthWatcher
}

func NewLogo() *Logo {
	l := Logo{depth: utils.NewDepthWatcher()}

	l.Iri = utils.NewElement("iri", "", IsValidIRI)

	l.InitCommonAttributes()
	l.depth.SetMaxDepth(1)

	return &l
}

func NewLogoExt(manager extension.Manager) *Logo {
	l := NewLogo()
	l.Extension = extension.InitExtension("logo", manager)

	return l
}

func (l *Logo) ProcessStartElement(el utils.StartElement) (utils.Visitor, utils.ParserError) {
	if l.depth.IsRoot() {
		l.ResetAttr()
		for _, attr := range el.Attr {
			if !l.ProcessAttr(attr) {
				l.Extension.ProcessAttr(attr, l)
			}
		}
	}

	if l.depth.Down() == utils.MaxDepthReached {
		return l, utils.NewError(LeafElementHasChild, "logo element should not have childs")
	}

	return l, nil
}

func (l *Logo) ProcessEndElement(el xml.EndElement) (utils.Visitor, utils.ParserError) {
	if l.depth.Up() == utils.RootLevel {
		return l.Parent, l.validate()
	}

	return l, nil
}

func (l *Logo) ProcessCharData(el xml.CharData) (utils.Visitor, utils.ParserError) {
	l.Iri.Value = string(el)
	return l, nil
}

func (l *Logo) validate() utils.ParserError {
	error := errors.NewErrorAggregator()

	utils.ValidateElement("logo", l.Iri, &error)
	l.ValidateCommonAttributes("logo", &error)
	l.Extension.Validate(&error)

	return error.ErrorObject()
}
