package rss

import (
	"github.com/jloup/errors"
	"github.com/jloup/xml/feed/extension"
	"github.com/jloup/xml/utils"
	"bytes"
	"encoding/xml"
	"strings"
)

type UnescapedContent struct {
	Content *bytes.Buffer
	name    xml.Name

	Encoder   *xml.Encoder
	Extension extension.VisitorExtension
	depth     utils.DepthWatcher
	Parent    utils.Visitor
}

func NewUnescapedContent() *UnescapedContent {
	u := UnescapedContent{depth: utils.NewDepthWatcher()}
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
	case utils.StartElement:
		err = u.Encoder.EncodeToken(*t.StartElement)
	case xml.EndElement:
		err = u.Encoder.EncodeToken(t)

	}
	return err
}

func (u *UnescapedContent) ProcessStartElement(el utils.StartElement) (utils.Visitor, utils.ParserError) {
	if u.depth.IsRoot() {
		u.name = el.Name
		u.Extension = extension.InitExtension(u.name.Local, u.Extension.Manager)

		for _, attr := range el.Attr {
			u.Extension.ProcessAttr(attr, u)
		}
	}

	err := errors.NewErrorAggregator()

	u.depth.Down()

	if error := u.EncodeXHTMLToken(el); error != nil {
		err.NewError(utils.NewError(XHTMLEncodeToStringError, "cannot encode XHTML"))
	}

	return u, err.ErrorObject()
}

func (u *UnescapedContent) ProcessEndElement(el xml.EndElement) (utils.Visitor, utils.ParserError) {
	level := u.depth.Up()
	if level == utils.ParentLevel {
		ferr := u.flush()
		if ferr != nil {
			return u.Parent, ferr
		}
		return u.Parent, u.Validate()
	}

	if err := u.EncodeXHTMLToken(el); err != nil {
		return u, utils.NewError(XHTMLEncodeToStringError, "cannot encode XHTML")
	}

	return u, nil
}

func (u *UnescapedContent) ProcessCharData(el xml.CharData) (utils.Visitor, utils.ParserError) {
	if len(strings.Fields(string(el))) > 0 {
		if ferr := u.flush(); ferr != nil {
			return u, ferr
		}

		if _, err := u.Content.Write(el); err != nil {
			return u, utils.NewError(CannotFlush, "cannot flush content")
		}
	}

	return u, nil
}

func (u *UnescapedContent) Validate() utils.ParserError {
	error := errors.NewErrorAggregator()

	u.Extension.Validate(&error)

	return error.ErrorObject()
}

func (u *UnescapedContent) flush() utils.ParserError {
	if err := u.Encoder.Flush(); err != nil {
		return utils.NewError(CannotFlush, "cannot flush content")
	}
	return nil
}

func (u *UnescapedContent) String() string {
	return string(u.Content.String())
}
