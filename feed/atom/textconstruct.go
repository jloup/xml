package atom

import (
	"encoding/xml"

	"github.com/jloup/utils"
	"github.com/jloup/xml/feed/extension"
	xmlutils "github.com/jloup/xml/utils"
)

type TextConstruct struct {
	CommonAttributes
	Type      string
	name      string
	XHTML     *InlineXHTMLContent
	PlainText *InlineTextContent

	Extension extension.VisitorExtension
	Parent    xmlutils.Visitor
}

func NewTextConstruct() *TextConstruct {
	t := TextConstruct{Type: "text"}
	t.XHTML = NewInlineXHTMLContent()
	t.PlainText = NewInlineTextContent()

	t.InitCommonAttributes()

	return &t
}

func NewTextConstructExt(manager extension.Manager) *TextConstruct {
	t := NewTextConstruct()

	t.Extension = extension.InitExtension("textconstruct", manager)
	return t
}

func (t *TextConstruct) reset() {
	t.ResetAttr()
	t.Type = "text"
}

func (t *TextConstruct) ProcessStartElement(el xmlutils.StartElement) (xmlutils.Visitor, xmlutils.ParserError) {
	t.name = el.Name.Local
	t.Extension = extension.InitExtension(t.name, t.Extension.Manager)
	t.reset()

	for _, attr := range el.Attr {
		switch attr.Name.Space {
		case xmlutils.XML_NS:
			t.ProcessAttr(attr)

		case "":
			if attr.Name.Local == "type" {
				t.Type = attr.Value
				break
			}
		default:
			t.Extension.ProcessAttr(attr, t)
		}
	}

	switch t.Type {
	case "text", "html":
		t.PlainText.Parent = t
		return t.PlainText, nil

	case "xhtml":
		t.XHTML.Parent = t
		return t.XHTML, nil
	}

	return t.PlainText, nil
}

func (t *TextConstruct) ProcessEndElement(el xml.EndElement) (xmlutils.Visitor, xmlutils.ParserError) {
	return t.Parent, t.Validate()
}

func (t *TextConstruct) Validate() xmlutils.ParserError {
	err := utils.NewErrorAggregator()

	t.ValidateCommonAttributes(t.name, &err)

	return err.ErrorObject()
}

func (t *TextConstruct) ProcessCharData(el xml.CharData) (xmlutils.Visitor, xmlutils.ParserError) {
	return t, nil
}

func (t *TextConstruct) String() string {
	if t.Type == "xhtml" {
		return t.XHTML.String()
	}
	return t.PlainText.String()
}
