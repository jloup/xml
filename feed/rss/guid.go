package rss

import (
	"encoding/xml"
	"strings"

	"github.com/jloup/errors"
	"github.com/jloup/xml/feed/extension"
	"github.com/jloup/xml/utils"
)

type Guid struct {
	IsPermalink utils.Element
	Content     utils.Element

	Extension extension.VisitorExtension
	Parent    utils.Visitor
	depth     utils.DepthWatcher
}

func NewGuid() *Guid {
	g := Guid{depth: utils.NewDepthWatcher()}

	g.IsPermalink = utils.NewElement("isPermalink", "true", utils.Nop)
	g.IsPermalink.SetOccurence(utils.NewOccurence("isPermalink", utils.UniqueValidator(AttributeDuplicated)))

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

func (g *Guid) ProcessStartElement(el utils.StartElement) (utils.Visitor, utils.ParserError) {
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

func (g *Guid) ProcessEndElement(el xml.EndElement) (utils.Visitor, utils.ParserError) {
	if g.depth.Up() == utils.RootLevel {
		return g.Parent, g.validate()
	}

	return g, nil
}

func (g *Guid) ProcessCharData(el xml.CharData) (utils.Visitor, utils.ParserError) {
	g.Content.Value = strings.TrimSpace(string(el))
	return g, nil
}

func (g *Guid) validate() utils.ParserError {
	error := errors.NewErrorAggregator()

	utils.ValidateElements("guid", &error, g.IsPermalink)
	g.Extension.Validate(&error)

	return error.ErrorObject()
}
