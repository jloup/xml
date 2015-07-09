package atom

import (
	"encoding/xml"
	"strings"

	"github.com/jloup/errors"
	"github.com/jloup/xml/feed/extension"
	"github.com/jloup/xml/utils"
)

type Content struct {
	CommonAttributes
	Type utils.Element
	Src  utils.Element

	XHTML            *InlineXHTMLContent
	PlainText        *InlineTextContent
	InlineContent    *InlineOtherContent
	OutOfLineContent *OutOfLineContent

	hasStarted bool
	Extension  extension.VisitorExtension
	Parent     utils.Visitor
}

func NewContent() *Content {
	c := Content{hasStarted: false}

	c.Type = utils.NewElement("type", "text", utils.Nop)
	c.Src = utils.NewElement("src", "", utils.Nop)

	c.XHTML = NewInlineXHTMLContent()
	c.PlainText = NewInlineTextContent()
	c.InlineContent = NewInlineOtherContent()
	c.OutOfLineContent = NewOutOfLineContent()

	c.XHTML.Parent = &c
	c.PlainText.Parent = &c
	c.InlineContent.Parent = &c
	c.OutOfLineContent.Parent = &c

	c.InitCommonAttributes()

	return &c
}

func NewContentExt(manager extension.Manager) *Content {
	c := NewContent()
	c.Extension = extension.InitExtension("content", manager)

	return c
}

func (c *Content) reset() {
	c.ResetAttr()
	c.hasStarted = false
	c.Type.Value = "text"
}

func (c *Content) HasReadableContent() bool {
	return c.hasStarted && (c.Type.Value == "text" || c.Type.Value == "html" || c.Type.Value == "xhtml" || c.InlineContent.HasReadableContent())
}

func (c *Content) ProcessStartElement(el utils.StartElement) (utils.Visitor, utils.ParserError) {
	c.reset()
	for _, attr := range el.Attr {
		switch attr.Name.Space {
		case utils.XML_NS:
			c.ProcessAttr(attr)
		case "":
			switch attr.Name.Local {
			case "type":
				c.Type.Value = attr.Value
			case "src":
				c.Src.Value = attr.Value
			}
		default:
			c.Extension.ProcessAttr(attr, c)
		}
	}

	c.hasStarted = true

	if c.Src.Value != "" {
		return c.OutOfLineContent.ProcessStartElement(el)
	}

	if c.Type.Value == "text" || c.Type.Value == "html" || strings.HasPrefix(c.Type.Value, "text/") {
		return c.PlainText, nil

	} else if c.Type.Value == "xhtml" {
		return c.XHTML, nil
	}

	return c.InlineContent.ProcessStartElement(el)
}

func (c *Content) ProcessEndElement(el xml.EndElement) (utils.Visitor, utils.ParserError) {
	if c.Type.Value != "text" && c.Type.Value != "html" && c.Type.Value != "xhtml" {
		if c.Parent != nil {
			return c.Parent.ProcessEndElement(el)
		}
		return nil, nil
	}
	return c.Parent, nil
}

func (c *Content) ProcessCharData(el xml.CharData) (utils.Visitor, utils.ParserError) {
	return c, nil
}

func (c *Content) validate() utils.ParserError {
	error := errors.NewErrorAggregator()

	c.Extension.Validate(&error)
	c.ValidateCommonAttributes("content", &error)

	return error.ErrorObject()
}

func (c *Content) String() string {

	if c.Src.Value != "" {

		return c.OutOfLineContent.Src.Value
	}

	if c.Type.Value == "text" || c.Type.Value == "html" {
		return c.PlainText.String()

	} else if c.Type.Value == "xhtml" {

		return c.XHTML.String()
	}
	return c.InlineContent.String()
}
