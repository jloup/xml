package atom

import (
	"encoding/xml"

	"github.com/jloup/utils"
	"github.com/jloup/xml/feed/extension"
	xmlutils "github.com/jloup/xml/utils"
)

type Generator struct {
	CommonAttributes
	Uri     xmlutils.Element
	Version xmlutils.Element
	Content string

	Extension extension.VisitorExtension
	Parent    xmlutils.Visitor
	depth     xmlutils.DepthWatcher
}

func NewGenerator() *Generator {
	g := Generator{depth: xmlutils.NewDepthWatcher()}

	g.Uri = xmlutils.NewElement("uri", "", IsValidIRI)
	g.Uri.SetOccurence(xmlutils.NewOccurence("uri", xmlutils.UniqueValidator(AttributeDuplicated)))

	g.Version = xmlutils.NewElement("version", "", xmlutils.Nop)
	g.Version.SetOccurence(xmlutils.NewOccurence("version", xmlutils.UniqueValidator(AttributeDuplicated)))

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

func (g *Generator) ProcessStartElement(el xmlutils.StartElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if g.depth.IsRoot() {
		g.reset()
		for _, attr := range el.Attr {
			switch attr.Name.Space {
			case xmlutils.XML_NS:
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

	if g.depth.Down() == xmlutils.MaxDepthReached {
		return g, xmlutils.NewError(LeafElementHasChild, "'generator' shoud not have childs")

	}

	return g, nil
}

func (g *Generator) ProcessEndElement(el xml.EndElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if g.depth.Up() == xmlutils.RootLevel {
		return g.Parent, g.validate()
	}

	return g, nil
}

func (g *Generator) ProcessCharData(el xml.CharData) (xmlutils.Visitor, xmlutils.ParserError) {
	g.Content = string(el)
	return g, nil
}

func (g *Generator) validate() xmlutils.ParserError {
	error := utils.NewErrorAggregator()

	xmlutils.ValidateElements("generator", &error, g.Uri, g.Version)
	g.Extension.Validate(&error)
	g.ValidateCommonAttributes("generator", &error)

	return error.ErrorObject()
}

func (g *Generator) String() string {
	return g.Content
}
