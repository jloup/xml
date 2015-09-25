package atom

import (
	"encoding/xml"

	"github.com/jloup/utils"
	"github.com/jloup/xml/feed/extension"
	xmlutils "github.com/jloup/xml/utils"
)

type Icon struct {
	CommonAttributes
	Iri xmlutils.Element

	Extension extension.VisitorExtension
	Parent    xmlutils.Visitor
	depth     xmlutils.DepthWatcher
}

func NewIcon() *Icon {
	i := Icon{depth: xmlutils.NewDepthWatcher()}

	i.Iri = xmlutils.NewElement("iri", "", IsValidIRI)

	i.InitCommonAttributes()
	i.depth.SetMaxDepth(1)

	return &i
}

func NewIconExt(manager extension.Manager) *Icon {
	i := NewIcon()

	i.Extension = extension.InitExtension("icon", manager)
	return i
}

func (i *Icon) ProcessStartElement(el xmlutils.StartElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if i.depth.IsRoot() {
		i.ResetAttr()
		for _, attr := range el.Attr {
			if !i.ProcessAttr(attr) {
				i.Extension.ProcessAttr(attr, i)
			}
		}
	}

	if i.depth.Down() == xmlutils.MaxDepthReached {
		return i, xmlutils.NewError(LeafElementHasChild, "icon element should not have childs")
	}

	return i, nil
}

func (i *Icon) ProcessEndElement(el xml.EndElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if i.depth.Up() == xmlutils.RootLevel {
		return i.Parent, i.validate()
	}

	return i, nil
}

func (i *Icon) ProcessCharData(el xml.CharData) (xmlutils.Visitor, xmlutils.ParserError) {
	i.Iri.Value = string(el)
	return i, nil
}

func (i *Icon) validate() xmlutils.ParserError {
	error := utils.NewErrorAggregator()

	xmlutils.ValidateElement("icon", i.Iri, &error)
	i.Extension.Validate(&error)
	i.ValidateCommonAttributes("icon", &error)

	return error.ErrorObject()
}
