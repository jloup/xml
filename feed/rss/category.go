package rss

import (
	"encoding/xml"
	"strings"

	"github.com/jloup/utils"
	"github.com/jloup/xml/feed/extension"
	xmlutils "github.com/jloup/xml/utils"
)

type Category struct {
	Domain  xmlutils.Element
	Content xmlutils.Element

	Extension extension.VisitorExtension
	Parent    xmlutils.Visitor
	depth     xmlutils.DepthWatcher
}

func NewCategory() *Category {
	c := Category{depth: xmlutils.NewDepthWatcher()}

	c.Domain = xmlutils.NewElement("domain", "", xmlutils.Nop)
	c.Domain.SetOccurence(xmlutils.NewOccurence("domain", xmlutils.UniqueValidator(AttributeDuplicated)))

	return &c
}

func NewCategoryExt(manager extension.Manager) *Category {
	c := NewCategory()
	c.Extension = extension.InitExtension("category", manager)

	return c
}

func (c *Category) reset() {
	c.Domain.Reset()
}

func (c *Category) ProcessStartElement(el xmlutils.StartElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if c.depth.IsRoot() {
		c.reset()
		for _, attr := range el.Attr {
			switch attr.Name.Space {
			case "":
				switch attr.Name.Local {
				case "domain":
					c.Domain.Value = attr.Value
					c.Domain.IncOccurence()
				}
			default:
				c.Extension.ProcessAttr(attr, c)
			}
		}
	}

	c.depth.Down()

	return c, nil
}

func (c *Category) ProcessEndElement(el xml.EndElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if c.depth.Up() == xmlutils.RootLevel {

		return c.Parent, c.validate()
	}

	return c, nil
}

func (c *Category) ProcessCharData(el xml.CharData) (xmlutils.Visitor, xmlutils.ParserError) {
	c.Content.Value = strings.TrimSpace(string(el))
	return c, nil
}

func (c *Category) validate() xmlutils.ParserError {
	error := utils.NewErrorAggregator()

	xmlutils.ValidateElements("category", &error, c.Domain)
	c.Extension.Validate(&error)

	return error.ErrorObject()
}
