package atom

import (
	"encoding/xml"

	"github.com/JLoup/errors"
	"github.com/JLoup/xml/feed/extension"
	"github.com/JLoup/xml/helper"
)

type Generator struct {
	CommonAttributes
	Uri     helper.Element
	Version helper.Element
	Content string

	Extension extension.VisitorExtension
	Parent    helper.Visitor
	depth     helper.DepthWatcher
}

func NewGenerator() *Generator {
	g := Generator{depth: helper.NewDepthWatcher()}

	g.Uri = helper.NewElement("uri", "", IsValidIRI)
	g.Uri.SetOccurence(helper.NewOccurence("uri", helper.UniqueValidator(AttributeDuplicated)))

	g.Version = helper.NewElement("version", "", helper.Nop)
	g.Version.SetOccurence(helper.NewOccurence("version", helper.UniqueValidator(AttributeDuplicated)))

	g.InitCommonAttributes()

	g.depth.SetMaxDepth(1)
	return &g
}

func NewGeneratorExt(manager extension.Manager) *Generator {
	g := NewGenerator()
	g.Extension = extension.InitExtension("generator", manager)

	return g
}

func (g *Generator) reset() {
	g.Uri.Reset()
	g.Version.Reset()
	g.ResetAttr()
}

func (g *Generator) ProcessStartElement(el helper.StartElement) (helper.Visitor, helper.ParserError) {
	if g.depth.IsRoot() {
		g.reset()
		for _, attr := range el.Attr {
			switch attr.Name.Space {
			case helper.XML_NS:
				g.ProcessAttr(attr)

			case "":
				switch attr.Name.Local {
				case "uri":
					g.Uri.Value = attr.Value
					g.Uri.IncOccurence()

				case "version":
					g.Version.Value = attr.Value
					g.Version.IncOccurence()
				}
			default:
				g.Extension.ProcessAttr(attr, g)
			}
		}
	}

	if g.depth.Down() == helper.MaxDepthReached {
		return g, helper.NewError(LeafElementHasChild, "'generator' shoud not have childs")

	}

	return g, nil
}

func (g *Generator) ProcessEndElement(el xml.EndElement) (helper.Visitor, helper.ParserError) {
	if g.depth.Up() == helper.RootLevel {
		return g.Parent, g.validate()
	}

	return g, nil
}

func (g *Generator) ProcessCharData(el xml.CharData) (helper.Visitor, helper.ParserError) {
	g.Content = string(el)
	return g, nil
}

func (g *Generator) validate() helper.ParserError {
	error := errors.NewErrorAggregator()

	helper.ValidateElements("generator", &error, g.Uri, g.Version)
	g.Extension.Validate(&error)
	g.ValidateCommonAttributes("generator", &error)

	return error.ErrorObject()
}

func (g *Generator) String() string {
	return g.Content
}
