package atom

import (
	"encoding/xml"
	"strings"

	"github.com/jloup/utils"
	"github.com/jloup/xml/feed/extension"
	xmlutils "github.com/jloup/xml/utils"
)

type Id struct {
	CommonAttributes
	Content xmlutils.Element

	Extension extension.VisitorExtension
	Parent    xmlutils.Visitor
	depth     xmlutils.DepthWatcher
}

func NewId() *Id {
	i := Id{depth: xmlutils.NewDepthWatcher()}

	i.Content = xmlutils.NewElement("iri", "", IsAbsoluteIRI)

	i.InitCommonAttributes()
	i.depth.SetMaxDepth(1)

	return &i
}

func NewIdExt(manager extension.Manager) *Id {
	i := NewId()

	i.Extension = extension.InitExtension("id", manager)

	return i
}
func (i *Id) ProcessStartElement(el xmlutils.StartElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if i.depth.IsRoot() {
		i.ResetAttr()
		for _, attr := range el.Attr {
			if !i.ProcessAttr(attr) {
				i.Extension.ProcessAttr(attr, i)
			}
		}
	}

	if i.depth.Down() == xmlutils.MaxDepthReached {
		return i, xmlutils.NewError(LeafElementHasChild, "id element should not have childs")
	}

	return i, nil
}

func (i *Id) ProcessEndElement(el xml.EndElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if i.depth.Up() == xmlutils.RootLevel {
		return i.Parent, i.validate()
	}

	return i, nil
}

func (i *Id) ProcessCharData(el xml.CharData) (xmlutils.Visitor, xmlutils.ParserError) {
	i.Content.Value = strings.TrimSpace(string(el))
	return i, nil
}

func (i *Id) validate() xmlutils.ParserError {
	error := utils.NewErrorAggregator()

	xmlutils.ValidateElement("id", i.Content, &error)
	i.Extension.Validate(&error)
	i.ValidateCommonAttributes("id", &error)

	return error.ErrorObject()
}

func (i *Id) String() string {
	return i.Content.String()
}
