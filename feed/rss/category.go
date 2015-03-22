package rss

import (
	"encoding/xml"
	"strings"

	"github.com/JLoup/errors"
	"github.com/JLoup/xml/feed/extension"
	"github.com/JLoup/xml/helper"
)

type Category struct {
	Domain  helper.Element
	Content helper.Element

	Extension extension.VisitorExtension
	Parent    helper.Visitor
	depth     helper.DepthWatcher
}

func NewCategory() *Category {
	c := Category{depth: helper.NewDepthWatcher()}

	c.Domain = helper.NewElement("domain", "", helper.Nop)
	c.Domain.SetOccurence(helper.NewOccurence("domain", helper.UniqueValidator(AttributeDuplicated)))

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

func (c *Category) ProcessStartElement(el helper.StartElement) (helper.Visitor, helper.ParserError) {
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

func (c *Category) ProcessEndElement(el xml.EndElement) (helper.Visitor, helper.ParserError) {
	if c.depth.Up() == helper.RootLevel {

		return c.Parent, c.validate()
	}

	return c, nil
}

func (c *Category) ProcessCharData(el xml.CharData) (helper.Visitor, helper.ParserError) {
	c.Content.Value = strings.TrimSpace(string(el))
	return c, nil
}

func (c *Category) validate() helper.ParserError {
	error := errors.NewErrorAggregator()

	helper.ValidateElements("category", &error, c.Domain)
	c.Extension.Validate(&error)

	return error.ErrorObject()
}
