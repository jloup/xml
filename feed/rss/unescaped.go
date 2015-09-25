package rss

import (
	"github.com/jloup/utils"
	"github.com/jloup/xml/feed/extension"
	xmlutils "github.com/jloup/xml/utils"
	"bytes"
	"encoding/xml"
	"strings"
)

type UnescapedContent struct {
	Content *bytes.Buffer
	name    xml.Name

	Encoder   *xml.Encoder
	Extension extension.VisitorExtension
	depth     xmlutils.DepthWatcher
	Parent    xmlutils.Visitor
}

func NewUnescapedContent() *UnescapedContent {
	u := UnescapedContent{depth: xmlutils.NewDepthWatcher()}
	u.Content = &bytes.Buffer{}
	u.Encoder = xml.NewEncoder(u.Content)
	return &u
}

func NewUnescapedContentExt(manager extension.Manager) *UnescapedContent {
	u := NewUnescapedContent()
	u.Extension = extension.InitExtension("unescaped", manager)

	return u
}

func (u *UnescapedContent) EncodeXHTMLToken(t xml.Token) error {
	var err error
	switch t := t.(type) {
	case xmlutils.StartElement:
		err = u.Encoder.EncodeToken(*t.StartElement)
	case xml.EndElement:
		err = u.Encoder.EncodeToken(t)

	}
	return err
}

func (u *UnescapedContent) ProcessStartElement(el xmlutils.StartElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if u.depth.IsRoot() {
		u.name = el.Name
		u.Extension = extension.InitExtension(u.name.Local, u.Extension.Manager)

		for _, attr := range el.Attr {
			u.Extension.ProcessAttr(attr, u)
		}
	}

	err := utils.NewErrorAggregator()

	u.depth.Down()

	if error := u.EncodeXHTMLToken(el); error != nil {
		err.NewError(xmlutils.NewError(XHTMLEncodeToStringError, "cannot encode XHTML"))
	}

	return u, err.ErrorObject()
}

func (u *UnescapedContent) ProcessEndElement(el xml.EndElement) (xmlutils.Visitor, xmlutils.ParserError) {
	level := u.depth.Up()
	if level == xmlutils.ParentLevel {
		ferr := u.flush()
		if ferr != nil {
			return u.Parent, ferr
		}
		return u.Parent, u.Validate()
	}

	if err := u.EncodeXHTMLToken(el); err != nil {
		return u, xmlutils.NewError(XHTMLEncodeToStringError, "cannot encode XHTML")
	}

	return u, nil
}

func (u *UnescapedContent) ProcessCharData(el xml.CharData) (xmlutils.Visitor, xmlutils.ParserError) {
	if len(strings.Fields(string(el))) > 0 {
		if ferr := u.flush(); ferr != nil {
			return u, ferr
		}

		if _, err := u.Content.Write(el); err != nil {
			return u, xmlutils.NewError(CannotFlush, "cannot flush content")
		}
	}

	return u, nil
}

func (u *UnescapedContent) Validate() xmlutils.ParserError {
	error := utils.NewErrorAggregator()

	u.Extension.Validate(&error)

	return error.ErrorObject()
}

func (u *UnescapedContent) flush() xmlutils.ParserError {
	if err := u.Encoder.Flush(); err != nil {
		return xmlutils.NewError(CannotFlush, "cannot flush content")
	}
	return nil
}

func (u *UnescapedContent) String() string {
	return string(u.Content.String())
}
