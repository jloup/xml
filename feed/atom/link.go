package atom

import (
	"encoding/xml"

	"github.com/JLoup/errors"
	"github.com/JLoup/xml/feed/extension"
	"github.com/JLoup/xml/utils"
)

type Link struct {
	CommonAttributes
	Href     utils.Element
	Rel      utils.Element
	Type     utils.Element
	HrefLang utils.Element
	Title    utils.Element
	Length   utils.Element

	Extension extension.VisitorExtension
	Parent    utils.Visitor
	depth     utils.DepthWatcher
}

func NewLink() *Link {
	l := Link{depth: utils.NewDepthWatcher()}

	l.Href = utils.NewElement("href", "", IsValidIRI)
	l.Href.SetOccurence(utils.NewOccurence("href", utils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	l.Rel = utils.NewElement("rel", "alternate", utils.Nop)
	l.Rel.SetOccurence(utils.NewOccurence("rel", utils.UniqueValidator(AttributeDuplicated)))

	l.Type = utils.NewElement("type", "", IsValidMIME)
	l.Type.SetOccurence(utils.NewOccurence("type", utils.UniqueValidator(AttributeDuplicated)))

	l.HrefLang = utils.NewElement("hreflang", "", utils.Nop)
	l.HrefLang.SetOccurence(utils.NewOccurence("hreflang", utils.UniqueValidator(AttributeDuplicated)))

	l.Title = utils.NewElement("title", "", utils.Nop)
	l.Title.SetOccurence(utils.NewOccurence("title", utils.UniqueValidator(AttributeDuplicated)))

	l.Length = utils.NewElement("length", "", IsValidLength)
	l.Length.SetOccurence(utils.NewOccurence("length", utils.UniqueValidator(AttributeDuplicated)))

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

func (l *Link) ProcessStartElement(el utils.StartElement) (utils.Visitor, utils.ParserError) {
	if l.depth.IsRoot() {
		l.reset()
		for _, attr := range el.Attr {
			switch attr.Name.Space {
			case utils.XML_NS:
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

func (l *Link) ProcessEndElement(el xml.EndElement) (utils.Visitor, utils.ParserError) {
	if l.depth.Up() == utils.RootLevel {
		return l.Parent, l.validate()
	}

	return l, nil
}

func (l *Link) ProcessCharData(el xml.CharData) (utils.Visitor, utils.ParserError) {
	return l, nil
}

func (l *Link) validate() utils.ParserError {
	error := errors.NewErrorAggregator()

	utils.ValidateElements("link", &error, l.Href, l.Rel, l.Type, l.HrefLang, l.Title, l.Length)
	l.Extension.Validate(&error)
	l.ValidateCommonAttributes("link", &error)

	return error.ErrorObject()
}
