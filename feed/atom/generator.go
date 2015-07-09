package atom

import (
	"encoding/xml"

	"github.com/jloup/errors"
	"github.com/jloup/xml/feed/extension"
	"github.com/jloup/xml/utils"
)

type Generator struct {
	CommonAttributes
	Uri     utils.Element
	Version utils.Element
	Content string

	Extension extension.VisitorExtension
	Parent    utils.Visitor
	depth     utils.DepthWatcher
}

func NewGenerator() *Generator {
	g := Generator{depth: utils.NewDepthWatcher()}

	g.Uri = utils.NewElement("uri", "", IsValidIRI)
	g.Uri.SetOccurence(utils.NewOccurence("uri", utils.UniqueValidator(AttributeDuplicated)))

	g.Version = utils.NewElement("version", "", utils.Nop)
	g.Version.SetOccurence(utils.NewOccurence("version", utils.UniqueValidator(AttributeDuplicated)))

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

func (g *Generator) ProcessStartElement(el utils.StartElement) (utils.Visitor, utils.ParserError) {
	if g.depth.IsRoot() {
		g.reset()
		for _, attr := range el.Attr {
			switch attr.Name.Space {
			case utils.XML_NS:
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

	if g.depth.Down() == utils.MaxDepthReached {
		return g, utils.NewError(LeafElementHasChild, "'generator' shoud not have childs")

	}

	return g, nil
}

func (g *Generator) ProcessEndElement(el xml.EndElement) (utils.Visitor, utils.ParserError) {
	if g.depth.Up() == utils.RootLevel {
		return g.Parent, g.validate()
	}

	return g, nil
}

func (g *Generator) ProcessCharData(el xml.CharData) (utils.Visitor, utils.ParserError) {
	g.Content = string(el)
	return g, nil
}

func (g *Generator) validate() utils.ParserError {
	error := errors.NewErrorAggregator()

	utils.ValidateElements("generator", &error, g.Uri, g.Version)
	g.Extension.Validate(&error)
	g.ValidateCommonAttributes("generator", &error)

	return error.ErrorObject()
}

func (g *Generator) String() string {
	return g.Content
}
