package atom

import (
	"encoding/xml"

	"github.com/jloup/errors"
	"github.com/jloup/xml/feed/extension"
	"github.com/jloup/xml/utils"
)

type Category struct {
	CommonAttributes
	Term   utils.Element
	Scheme utils.Element
	Label  utils.Element
	Parent utils.Visitor

	Extension extension.VisitorExtension
	depth     utils.DepthWatcher
}

func NewCategory() *Category {
	c := Category{depth: utils.NewDepthWatcher()}

	c.Term = utils.NewElement("term", "", utils.Nop)
	c.Term.SetOccurence(utils.NewOccurence("term", utils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	c.Scheme = utils.NewElement("scheme", "", IsValidIRI)
	c.Scheme.SetOccurence(utils.NewOccurence("scheme", utils.UniqueValidator(AttributeDuplicated)))

	c.Label = utils.NewElement("label", "", utils.Nop)
	c.Label.SetOccurence(utils.NewOccurence("label", utils.UniqueValidator(AttributeDuplicated)))

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

func (c *Category) ProcessStartElement(el utils.StartElement) (utils.Visitor, utils.ParserError) {
	if c.depth.IsRoot() {
		c.reset()
		for _, attr := range el.Attr {
			switch attr.Name.Space {
			case utils.XML_NS:
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

func (c *Category) ProcessEndElement(el xml.EndElement) (utils.Visitor, utils.ParserError) {
	if c.depth.Up() == utils.RootLevel {
		return c.Parent, c.validate()
	}

	return c, nil
}

func (c *Category) ProcessCharData(el xml.CharData) (utils.Visitor, utils.ParserError) {
	return c, nil
}

func (c *Category) validate() utils.ParserError {
	error := errors.NewErrorAggregator()

	utils.ValidateElements("category", &error, c.Term, c.Label, c.Scheme)
	c.ValidateCommonAttributes("category", &error)

	c.Extension.Validate(&error)

	return error.ErrorObject()
}
