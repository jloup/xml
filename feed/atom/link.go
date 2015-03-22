package atom

import (
	"encoding/xml"

	"github.com/JLoup/errors"
	"github.com/JLoup/xml/feed/extension"
	"github.com/JLoup/xml/helper"
)

type Link struct {
	CommonAttributes
	Href     helper.Element
	Rel      helper.Element
	Type     helper.Element
	HrefLang helper.Element
	Title    helper.Element
	Length   helper.Element

	Extension extension.VisitorExtension
	Parent    helper.Visitor
	depth     helper.DepthWatcher
}

func NewLink() *Link {
	l := Link{depth: helper.NewDepthWatcher()}

	l.Href = helper.NewElement("href", "", IsValidIRI)
	l.Href.SetOccurence(helper.NewOccurence("href", helper.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	l.Rel = helper.NewElement("rel", "alternate", helper.Nop)
	l.Rel.SetOccurence(helper.NewOccurence("rel", helper.UniqueValidator(AttributeDuplicated)))

	l.Type = helper.NewElement("type", "", IsValidMIME)
	l.Type.SetOccurence(helper.NewOccurence("type", helper.UniqueValidator(AttributeDuplicated)))

	l.HrefLang = helper.NewElement("hreflang", "", helper.Nop)
	l.HrefLang.SetOccurence(helper.NewOccurence("hreflang", helper.UniqueValidator(AttributeDuplicated)))

	l.Title = helper.NewElement("title", "", helper.Nop)
	l.Title.SetOccurence(helper.NewOccurence("title", helper.UniqueValidator(AttributeDuplicated)))

	l.Length = helper.NewElement("length", "", IsValidLength)
	l.Length.SetOccurence(helper.NewOccurence("length", helper.UniqueValidator(AttributeDuplicated)))

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

func (l *Link) ProcessStartElement(el helper.StartElement) (helper.Visitor, helper.ParserError) {
	if l.depth.IsRoot() {
		l.reset()
		for _, attr := range el.Attr {
			switch attr.Name.Space {
			case helper.XML_NS:
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

func (l *Link) ProcessEndElement(el xml.EndElement) (helper.Visitor, helper.ParserError) {
	if l.depth.Up() == helper.RootLevel {
		return l.Parent, l.validate()
	}

	return l, nil
}

func (l *Link) ProcessCharData(el xml.CharData) (helper.Visitor, helper.ParserError) {
	return l, nil
}

func (l *Link) validate() helper.ParserError {
	error := errors.NewErrorAggregator()

	helper.ValidateElements("link", &error, l.Href, l.Rel, l.Type, l.HrefLang, l.Title, l.Length)
	l.Extension.Validate(&error)
	l.ValidateCommonAttributes("link", &error)

	return error.ErrorObject()
}
