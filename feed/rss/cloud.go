package rss

import (
	"encoding/xml"

	"github.com/jloup/errors"
	"github.com/jloup/xml/feed/extension"
	"github.com/jloup/xml/utils"
)

type Cloud struct {
	Domain            utils.Element
	Port              utils.Element
	Path              utils.Element
	RegisterProcedure utils.Element
	Protocol          utils.Element

	Extension extension.VisitorExtension
	Parent    utils.Visitor
	depth     utils.DepthWatcher
}

func NewCloud() *Cloud {
	c := Cloud{depth: utils.NewDepthWatcher()}

	c.Domain = utils.NewElement("domain", "", utils.Nop)
	c.Domain.SetOccurence(utils.NewOccurence("domain", utils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	c.Port = utils.NewElement("port", "", utils.Nop)
	c.Port.SetOccurence(utils.NewOccurence("port", utils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	c.Path = utils.NewElement("path", "", utils.Nop)
	c.Path.SetOccurence(utils.NewOccurence("path", utils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	c.RegisterProcedure = utils.NewElement("registerProcedure", "", utils.Nop)
	c.RegisterProcedure.SetOccurence(utils.NewOccurence("registerProcedure", utils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	c.Protocol = utils.NewElement("protocol", "", utils.Nop)
	c.Protocol.SetOccurence(utils.NewOccurence("protocol", utils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

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

func (c *Cloud) ProcessStartElement(el utils.StartElement) (utils.Visitor, utils.ParserError) {
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

func (c *Cloud) ProcessEndElement(el xml.EndElement) (utils.Visitor, utils.ParserError) {
	if c.depth.Up() == utils.RootLevel {
		return c.Parent, c.validate()
	}

	return c, nil
}

func (c *Cloud) ProcessCharData(el xml.CharData) (utils.Visitor, utils.ParserError) {
	return c, nil
}

func (c *Cloud) validate() utils.ParserError {
	error := errors.NewErrorAggregator()

	utils.ValidateElements("cloud", &error, c.Domain, c.Port, c.Path, c.RegisterProcedure, c.Protocol)
	c.Extension.Validate(&error)

	return error.ErrorObject()
}
