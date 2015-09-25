package atom

import (
	"encoding/xml"

	"github.com/jloup/utils"
	"github.com/jloup/xml/feed/extension"
	xmlutils "github.com/jloup/xml/utils"
)

type Link struct {
	CommonAttributes
	Href     xmlutils.Element
	Rel      xmlutils.Element
	Type     xmlutils.Element
	HrefLang xmlutils.Element
	Title    xmlutils.Element
	Length   xmlutils.Element

	Extension extension.VisitorExtension
	Parent    xmlutils.Visitor
	depth     xmlutils.DepthWatcher
}

func NewLink() *Link {
	l := Link{depth: xmlutils.NewDepthWatcher()}

	l.Href = xmlutils.NewElement("href", "", IsValidIRI)
	l.Href.SetOccurence(xmlutils.NewOccurence("href", xmlutils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	l.Rel = xmlutils.NewElement("rel", "alternate", xmlutils.Nop)
	l.Rel.SetOccurence(xmlutils.NewOccurence("rel", xmlutils.UniqueValidator(AttributeDuplicated)))

	l.Type = xmlutils.NewElement("type", "", IsValidMIME)
	l.Type.SetOccurence(xmlutils.NewOccurence("type", xmlutils.UniqueValidator(AttributeDuplicated)))

	l.HrefLang = xmlutils.NewElement("hreflang", "", xmlutils.Nop)
	l.HrefLang.SetOccurence(xmlutils.NewOccurence("hreflang", xmlutils.UniqueValidator(AttributeDuplicated)))

	l.Title = xmlutils.NewElement("title", "", xmlutils.Nop)
	l.Title.SetOccurence(xmlutils.NewOccurence("title", xmlutils.UniqueValidator(AttributeDuplicated)))

	l.Length = xmlutils.NewElement("length", "", IsValidLength)
	l.Length.SetOccurence(xmlutils.NewOccurence("length", xmlutils.UniqueValidator(AttributeDuplicated)))

	l.InitCommonAttributes()

	return &l
}

func NewLinkExt(manager extension.Manager) *Link {
	l := NewLink()
	l.Extension = extension.InitExtension("link", manager)

	return l
}

func (l *Link) reset() {
	l.Href.Reset()

	l.Rel.Reset()
	l.Rel.Value = "alternate"

	l.Type.Reset()
	l.HrefLang.Reset()
	l.Title.Reset()
	l.Length.Reset()

	l.ResetAttr()
}

func (l *Link) ProcessStartElement(el xmlutils.StartElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if l.depth.IsRoot() {
		l.reset()
		for _, attr := range el.Attr {
			switch attr.Name.Space {
			case xmlutils.XML_NS:
				l.ProcessAttr(attr)

			case "":
				switch attr.Name.Local {
				case "href":
					l.Href.IncOccurence()
					l.Href.Value = attr.Value

				case "rel":
					l.Rel.IncOccurence()
					l.Rel.Value = attr.Value

				case "type":
					l.Type.IncOccurence()
					l.Type.Value = attr.Value

				case "hreflang":
					l.HrefLang.IncOccurence()
					l.HrefLang.Value = attr.Value

				case "title":
					l.Title.IncOccurence()
					l.Title.Value = attr.Value

				case "length":
					l.Length.IncOccurence()
					l.Length.Value = attr.Value
				}
			default:
				l.Extension.ProcessAttr(attr, l)
			}
		}
	}

	l.depth.Down()
	return l, nil
}

func (l *Link) ProcessEndElement(el xml.EndElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if l.depth.Up() == xmlutils.RootLevel {
		return l.Parent, l.validate()
	}

	return l, nil
}

func (l *Link) ProcessCharData(el xml.CharData) (xmlutils.Visitor, xmlutils.ParserError) {
	return l, nil
}

func (l *Link) validate() xmlutils.ParserError {
	error := utils.NewErrorAggregator()

	xmlutils.ValidateElements("link", &error, l.Href, l.Rel, l.Type, l.HrefLang, l.Title, l.Length)
	l.Extension.Validate(&error)
	l.ValidateCommonAttributes("link", &error)

	return error.ErrorObject()
}
