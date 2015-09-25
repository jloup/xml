package rss

import (
	"encoding/xml"
	"strings"

	"github.com/jloup/utils"
	"github.com/jloup/xml/feed/extension"
	xmlutils "github.com/jloup/xml/utils"
)

type Guid struct {
	IsPermalink xmlutils.Element
	Content     xmlutils.Element

	Extension extension.VisitorExtension
	Parent    xmlutils.Visitor
	depth     xmlutils.DepthWatcher
}

func NewGuid() *Guid {
	g := Guid{depth: xmlutils.NewDepthWatcher()}

	g.IsPermalink = xmlutils.NewElement("isPermalink", "true", xmlutils.Nop)
	g.IsPermalink.SetOccurence(xmlutils.NewOccurence("isPermalink", xmlutils.UniqueValidator(AttributeDuplicated)))

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

func (g *Guid) ProcessStartElement(el xmlutils.StartElement) (xmlutils.Visitor, xmlutils.ParserError) {
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

func (g *Guid) ProcessEndElement(el xml.EndElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if g.depth.Up() == xmlutils.RootLevel {
		return g.Parent, g.validate()
	}

	return g, nil
}

func (g *Guid) ProcessCharData(el xml.CharData) (xmlutils.Visitor, xmlutils.ParserError) {
	g.Content.Value = strings.TrimSpace(string(el))
	return g, nil
}

func (g *Guid) validate() xmlutils.ParserError {
	error := utils.NewErrorAggregator()

	xmlutils.ValidateElements("guid", &error, g.IsPermalink)
	g.Extension.Validate(&error)

	return error.ErrorObject()
}
