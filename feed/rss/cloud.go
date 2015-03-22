package rss

import (
	"encoding/xml"

	"github.com/JLoup/errors"
	"github.com/JLoup/xml/feed/extension"
	"github.com/JLoup/xml/helper"
)

type Cloud struct {
	Domain            helper.Element
	Port              helper.Element
	Path              helper.Element
	RegisterProcedure helper.Element
	Protocol          helper.Element

	Extension extension.VisitorExtension
	Parent    helper.Visitor
	depth     helper.DepthWatcher
}

func NewCloud() *Cloud {
	c := Cloud{depth: helper.NewDepthWatcher()}

	c.Domain = helper.NewElement("domain", "", helper.Nop)
	c.Domain.SetOccurence(helper.NewOccurence("domain", helper.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	c.Port = helper.NewElement("port", "", helper.Nop)
	c.Port.SetOccurence(helper.NewOccurence("port", helper.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	c.Path = helper.NewElement("path", "", helper.Nop)
	c.Path.SetOccurence(helper.NewOccurence("path", helper.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	c.RegisterProcedure = helper.NewElement("registerProcedure", "", helper.Nop)
	c.RegisterProcedure.SetOccurence(helper.NewOccurence("registerProcedure", helper.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	c.Protocol = helper.NewElement("protocol", "", helper.Nop)
	c.Protocol.SetOccurence(helper.NewOccurence("protocol", helper.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

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

func (c *Cloud) ProcessStartElement(el helper.StartElement) (helper.Visitor, helper.ParserError) {
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

func (c *Cloud) ProcessEndElement(el xml.EndElement) (helper.Visitor, helper.ParserError) {
	if c.depth.Up() == helper.RootLevel {
		return c.Parent, c.validate()
	}

	return c, nil
}

func (c *Cloud) ProcessCharData(el xml.CharData) (helper.Visitor, helper.ParserError) {
	return c, nil
}

func (c *Cloud) validate() helper.ParserError {
	error := errors.NewErrorAggregator()

	helper.ValidateElements("cloud", &error, c.Domain, c.Port, c.Path, c.RegisterProcedure, c.Protocol)
	c.Extension.Validate(&error)

	return error.ErrorObject()
}
