package rss

import (
	"encoding/xml"

	"github.com/jloup/utils"
	"github.com/jloup/xml/feed/extension"
	xmlutils "github.com/jloup/xml/utils"
)

type Cloud struct {
	Domain            xmlutils.Element
	Port              xmlutils.Element
	Path              xmlutils.Element
	RegisterProcedure xmlutils.Element
	Protocol          xmlutils.Element

	Extension extension.VisitorExtension
	Parent    xmlutils.Visitor
	depth     xmlutils.DepthWatcher
}

func NewCloud() *Cloud {
	c := Cloud{depth: xmlutils.NewDepthWatcher()}

	c.Domain = xmlutils.NewElement("domain", "", xmlutils.Nop)
	c.Domain.SetOccurence(xmlutils.NewOccurence("domain", xmlutils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	c.Port = xmlutils.NewElement("port", "", xmlutils.Nop)
	c.Port.SetOccurence(xmlutils.NewOccurence("port", xmlutils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	c.Path = xmlutils.NewElement("path", "", xmlutils.Nop)
	c.Path.SetOccurence(xmlutils.NewOccurence("path", xmlutils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	c.RegisterProcedure = xmlutils.NewElement("registerProcedure", "", xmlutils.Nop)
	c.RegisterProcedure.SetOccurence(xmlutils.NewOccurence("registerProcedure", xmlutils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	c.Protocol = xmlutils.NewElement("protocol", "", xmlutils.Nop)
	c.Protocol.SetOccurence(xmlutils.NewOccurence("protocol", xmlutils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	return &c
}

func NewCloudExt(manager extension.Manager) *Cloud {
	c := NewCloud()

	c.Extension = extension.InitExtension("cloud", manager)

	return c
}

func (c *Cloud) reset() {
	c.Domain.Reset()
	c.Port.Reset()
	c.Path.Reset()
	c.RegisterProcedure.Reset()
	c.Protocol.Reset()
}

func (c *Cloud) ProcessStartElement(el xmlutils.StartElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if c.depth.IsRoot() {
		c.reset()
		for _, attr := range el.Attr {
			switch attr.Name.Space {
			case "":
				switch attr.Name.Local {
				case "domain":
					c.Domain.Value = attr.Value
					c.Domain.IncOccurence()

				case "port":
					c.Port.Value = attr.Value
					c.Port.IncOccurence()

				case "path":
					c.Path.Value = attr.Value
					c.Path.IncOccurence()

				case "registerprocedure":
					c.RegisterProcedure.Value = attr.Value
					c.RegisterProcedure.IncOccurence()

				case "protocol":
					c.Protocol.Value = attr.Value
					c.Protocol.IncOccurence()
				}
			default:
				c.Extension.ProcessAttr(attr, c)
			}
		}
	}

	c.depth.Down()

	return c, nil
}

func (c *Cloud) ProcessEndElement(el xml.EndElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if c.depth.Up() == xmlutils.RootLevel {
		return c.Parent, c.validate()
	}

	return c, nil
}

func (c *Cloud) ProcessCharData(el xml.CharData) (xmlutils.Visitor, xmlutils.ParserError) {
	return c, nil
}

func (c *Cloud) validate() xmlutils.ParserError {
	error := utils.NewErrorAggregator()

	xmlutils.ValidateElements("cloud", &error, c.Domain, c.Port, c.Path, c.RegisterProcedure, c.Protocol)
	c.Extension.Validate(&error)

	return error.ErrorObject()
}
