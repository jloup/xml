package atom

import (
	"encoding/xml"

	"github.com/jloup/utils"
	"github.com/jloup/xml/feed/extension"
	xmlutils "github.com/jloup/xml/utils"
)

type Category struct {
	CommonAttributes
	Term   xmlutils.Element
	Scheme xmlutils.Element
	Label  xmlutils.Element
	Parent xmlutils.Visitor

	Extension extension.VisitorExtension
	depth     xmlutils.DepthWatcher
}

func NewCategory() *Category {
	c := Category{depth: xmlutils.NewDepthWatcher()}

	c.Term = xmlutils.NewElement("term", "", xmlutils.Nop)
	c.Term.SetOccurence(xmlutils.NewOccurence("term", xmlutils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	c.Scheme = xmlutils.NewElement("scheme", "", IsValidIRI)
	c.Scheme.SetOccurence(xmlutils.NewOccurence("scheme", xmlutils.UniqueValidator(AttributeDuplicated)))

	c.Label = xmlutils.NewElement("label", "", xmlutils.Nop)
	c.Label.SetOccurence(xmlutils.NewOccurence("label", xmlutils.UniqueValidator(AttributeDuplicated)))

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

func (c *Category) ProcessStartElement(el xmlutils.StartElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if c.depth.IsRoot() {
		c.reset()
		for _, attr := range el.Attr {
			switch attr.Name.Space {
			case xmlutils.XML_NS:
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

func (c *Category) ProcessEndElement(el xml.EndElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if c.depth.Up() == xmlutils.RootLevel {
		return c.Parent, c.validate()
	}

	return c, nil
}

func (c *Category) ProcessCharData(el xml.CharData) (xmlutils.Visitor, xmlutils.ParserError) {
	return c, nil
}

func (c *Category) validate() xmlutils.ParserError {
	error := utils.NewErrorAggregator()

	xmlutils.ValidateElements("category", &error, c.Term, c.Label, c.Scheme)
	c.ValidateCommonAttributes("category", &error)

	c.Extension.Validate(&error)

	return error.ErrorObject()
}
