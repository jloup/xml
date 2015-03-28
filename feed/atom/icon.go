package atom

import (
	"encoding/xml"

	"github.com/JLoup/errors"
	"github.com/JLoup/xml/feed/extension"
	"github.com/JLoup/xml/utils"
)

type Icon struct {
	CommonAttributes
	Iri utils.Element

	Extension extension.VisitorExtension
	Parent    utils.Visitor
	depth     utils.DepthWatcher
}

func NewIcon() *Icon {
	i := Icon{depth: utils.NewDepthWatcher()}

	i.Iri = utils.NewElement("iri", "", IsValidIRI)

	i.InitCommonAttributes()
	i.depth.SetMaxDepth(1)

	return &i
}

func NewIconExt(manager extension.Manager) *Icon {
	i := NewIcon()

	i.Extension = extension.InitExtension("icon", manager)
	return i
}

func (i *Icon) ProcessStartElement(el utils.StartElement) (utils.Visitor, utils.ParserError) {
	if i.depth.IsRoot() {
		i.ResetAttr()
		for _, attr := range el.Attr {
			if !i.ProcessAttr(attr) {
				i.Extension.ProcessAttr(attr, i)
			}
		}
	}

	if i.depth.Down() == utils.MaxDepthReached {
		return i, utils.NewError(LeafElementHasChild, "icon element should not have childs")
	}

	return i, nil
}

func (i *Icon) ProcessEndElement(el xml.EndElement) (utils.Visitor, utils.ParserError) {
	if i.depth.Up() == utils.RootLevel {
		return i.Parent, i.validate()
	}

	return i, nil
}

func (i *Icon) ProcessCharData(el xml.CharData) (utils.Visitor, utils.ParserError) {
	i.Iri.Value = string(el)
	return i, nil
}

func (i *Icon) validate() utils.ParserError {
	error := errors.NewErrorAggregator()

	utils.ValidateElement("icon", i.Iri, &error)
	i.Extension.Validate(&error)
	i.ValidateCommonAttributes("icon", &error)

	return error.ErrorObject()
}
