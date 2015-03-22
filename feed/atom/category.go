package atom

import (
	"encoding/xml"

	"github.com/JLoup/errors"
	"github.com/JLoup/xml/feed/extension"
	"github.com/JLoup/xml/helper"
)

type Category struct {
	CommonAttributes
	Term   helper.Element
	Scheme helper.Element
	Label  helper.Element
	Parent helper.Visitor

	Extension extension.VisitorExtension
	depth     helper.DepthWatcher
}

func NewCategory() *Category {
	c := Category{depth: helper.NewDepthWatcher()}

	c.Term = helper.NewElement("term", "", helper.Nop)
	c.Term.SetOccurence(helper.NewOccurence("term", helper.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	c.Scheme = helper.NewElement("scheme", "", IsValidIRI)
	c.Scheme.SetOccurence(helper.NewOccurence("scheme", helper.UniqueValidator(AttributeDuplicated)))

	c.Label = helper.NewElement("label", "", helper.Nop)
	c.Label.SetOccurence(helper.NewOccurence("label", helper.UniqueValidator(AttributeDuplicated)))

	c.InitCommonAttributes()

	return &c
}

func NewCategoryExt(manager extension.Manager) *Category {
	c := NewCategory()
	c.Extension = extension.InitExtension("category", manager)

	return c
}

func (c *Category) reset() {
	c.Term.Reset()
	c.Label.Reset()
	c.Scheme.Reset()
	c.ResetAttr()
}

func (c *Category) ProcessStartElement(el helper.StartElement) (helper.Visitor, helper.ParserError) {
	if c.depth.IsRoot() {
		c.reset()
		for _, attr := range el.Attr {
			switch attr.Name.Space {
			case helper.XML_NS:
				c.ProcessAttr(attr)
			case "":
				switch attr.Name.Local {
				case "term":
					c.Term.Value = attr.Value
					c.Term.IncOccurence()

				case "label":
					c.Label.Value = attr.Value
					c.Label.IncOccurence()

				case "scheme":
					c.Scheme.Value = attr.Value
					c.Scheme.IncOccurence()
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
	return c, nil
}

func (c *Category) validate() helper.ParserError {
	error := errors.NewErrorAggregator()

	helper.ValidateElements("category", &error, c.Term, c.Label, c.Scheme)
	c.ValidateCommonAttributes("category", &error)

	c.Extension.Validate(&error)

	return error.ErrorObject()
}
