package rss

import (
	"encoding/xml"
	"strings"

	"github.com/JLoup/errors"
	"github.com/JLoup/xml/feed/extension"
	"github.com/JLoup/xml/helper"
)

type Guid struct {
	IsPermalink helper.Element
	Content     helper.Element

	Extension extension.VisitorExtension
	Parent    helper.Visitor
	depth     helper.DepthWatcher
}

func NewGuid() *Guid {
	g := Guid{depth: helper.NewDepthWatcher()}

	g.IsPermalink = helper.NewElement("isPermalink", "true", helper.Nop)
	g.IsPermalink.SetOccurence(helper.NewOccurence("isPermalink", helper.UniqueValidator(AttributeDuplicated)))

	return &g
}

func NewGuidExt(manager extension.Manager) *Guid {
	g := NewGuid()
	g.Extension = extension.InitExtension("guid", manager)

	return g
}

func (g *Guid) reset() {
	g.IsPermalink.Reset()
	g.IsPermalink.Value = "true"
}

func (g *Guid) ProcessStartElement(el helper.StartElement) (helper.Visitor, helper.ParserError) {
	if g.depth.IsRoot() {
		g.reset()
		for _, attr := range el.Attr {
			switch attr.Name.Space {
			case "":
				switch attr.Name.Local {
				case "ispermalink":
					g.IsPermalink.Value = attr.Value
					g.IsPermalink.IncOccurence()
				}
			default:
				g.Extension.ProcessAttr(attr, g)
			}
		}
	}

	g.depth.Down()

	return g, nil
}

func (g *Guid) ProcessEndElement(el xml.EndElement) (helper.Visitor, helper.ParserError) {
	if g.depth.Up() == helper.RootLevel {
		return g.Parent, g.validate()
	}

	return g, nil
}

func (g *Guid) ProcessCharData(el xml.CharData) (helper.Visitor, helper.ParserError) {
	g.Content.Value = strings.TrimSpace(string(el))
	return g, nil
}

func (g *Guid) validate() helper.ParserError {
	error := errors.NewErrorAggregator()

	helper.ValidateElements("guid", &error, g.IsPermalink)
	g.Extension.Validate(&error)

	return error.ErrorObject()
}
