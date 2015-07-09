package rss

import (
	"encoding/xml"
	"strings"

	"github.com/jloup/errors"
	"github.com/jloup/xml/feed/extension"
	"github.com/jloup/xml/utils"
)

type Category struct {
	Domain  utils.Element
	Content utils.Element

	Extension extension.VisitorExtension
	Parent    utils.Visitor
	depth     utils.DepthWatcher
}

func NewCategory() *Category {
	c := Category{depth: utils.NewDepthWatcher()}

	c.Domain = utils.NewElement("domain", "", utils.Nop)
	c.Domain.SetOccurence(utils.NewOccurence("domain", utils.UniqueValidator(AttributeDuplicated)))

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

func (c *Category) ProcessStartElement(el utils.StartElement) (utils.Visitor, utils.ParserError) {
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

func (c *Category) ProcessEndElement(el xml.EndElement) (utils.Visitor, utils.ParserError) {
	if c.depth.Up() == utils.RootLevel {

		return c.Parent, c.validate()
	}

	return c, nil
}

func (c *Category) ProcessCharData(el xml.CharData) (utils.Visitor, utils.ParserError) {
	c.Content.Value = strings.TrimSpace(string(el))
	return c, nil
}

func (c *Category) validate() utils.ParserError {
	error := errors.NewErrorAggregator()

	utils.ValidateElements("category", &error, c.Domain)
	c.Extension.Validate(&error)

	return error.ErrorObject()
}
