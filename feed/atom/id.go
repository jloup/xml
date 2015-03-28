package atom

import (
	"encoding/xml"
	"strings"

	"github.com/JLoup/errors"
	"github.com/JLoup/xml/feed/extension"
	"github.com/JLoup/xml/utils"
)

type Id struct {
	CommonAttributes
	Content utils.Element

	Extension extension.VisitorExtension
	Parent    utils.Visitor
	depth     utils.DepthWatcher
}

func NewId() *Id {
	i := Id{depth: utils.NewDepthWatcher()}

	i.Content = utils.NewElement("iri", "", IsAbsoluteIRI)

	i.InitCommonAttributes()
	i.depth.SetMaxDepth(1)

	return &i
}

func NewIdExt(manager extension.Manager) *Id {
	i := NewId()

	i.Extension = extension.InitExtension("id", manager)

	return i
}
func (i *Id) ProcessStartElement(el utils.StartElement) (utils.Visitor, utils.ParserError) {
	if i.depth.IsRoot() {
		i.ResetAttr()
		for _, attr := range el.Attr {
			if !i.ProcessAttr(attr) {
				i.Extension.ProcessAttr(attr, i)
			}
		}
	}

	if i.depth.Down() == utils.MaxDepthReached {
		return i, utils.NewError(LeafElementHasChild, "id element should not have childs")
	}

	return i, nil
}

func (i *Id) ProcessEndElement(el xml.EndElement) (utils.Visitor, utils.ParserError) {
	if i.depth.Up() == utils.RootLevel {
		return i.Parent, i.validate()
	}

	return i, nil
}

func (i *Id) ProcessCharData(el xml.CharData) (utils.Visitor, utils.ParserError) {
	i.Content.Value = strings.TrimSpace(string(el))
	return i, nil
}

func (i *Id) validate() utils.ParserError {
	error := errors.NewErrorAggregator()

	utils.ValidateElement("id", i.Content, &error)
	i.Extension.Validate(&error)
	i.ValidateCommonAttributes("id", &error)

	return error.ErrorObject()
}

func (i *Id) String() string {
	return i.Content.String()
}
