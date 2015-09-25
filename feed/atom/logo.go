package atom

import (
	"encoding/xml"

	"github.com/jloup/utils"
	"github.com/jloup/xml/feed/extension"
	xmlutils "github.com/jloup/xml/utils"
)

type Logo struct {
	CommonAttributes
	Iri xmlutils.Element

	Extension extension.VisitorExtension
	Parent    xmlutils.Visitor
	depth     xmlutils.DepthWatcher
}

func NewLogo() *Logo {
	l := Logo{depth: xmlutils.NewDepthWatcher()}

	l.Iri = xmlutils.NewElement("iri", "", IsValidIRI)

	l.InitCommonAttributes()
	l.depth.SetMaxDepth(1)

	return &l
}

func NewLogoExt(manager extension.Manager) *Logo {
	l := NewLogo()
	l.Extension = extension.InitExtension("logo", manager)

	return l
}

func (l *Logo) ProcessStartElement(el xmlutils.StartElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if l.depth.IsRoot() {
		l.ResetAttr()
		for _, attr := range el.Attr {
			if !l.ProcessAttr(attr) {
				l.Extension.ProcessAttr(attr, l)
			}
		}
	}

	if l.depth.Down() == xmlutils.MaxDepthReached {
		return l, xmlutils.NewError(LeafElementHasChild, "logo element should not have childs")
	}

	return l, nil
}

func (l *Logo) ProcessEndElement(el xml.EndElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if l.depth.Up() == xmlutils.RootLevel {
		return l.Parent, l.validate()
	}

	return l, nil
}

func (l *Logo) ProcessCharData(el xml.CharData) (xmlutils.Visitor, xmlutils.ParserError) {
	l.Iri.Value = string(el)
	return l, nil
}

func (l *Logo) validate() xmlutils.ParserError {
	error := utils.NewErrorAggregator()

	xmlutils.ValidateElement("logo", l.Iri, &error)
	l.ValidateCommonAttributes("logo", &error)
	l.Extension.Validate(&error)

	return error.ErrorObject()
}
