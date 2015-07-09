package atom

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/jloup/errors"
	"github.com/jloup/xml/utils"
)

type InlineXHTMLContent struct {
	Content *bytes.Buffer

	Encoder   *xml.Encoder
	depth     utils.DepthWatcher
	completed bool
	Parent    utils.Visitor
}

func NewInlineXHTMLContent() *InlineXHTMLContent {
	i := InlineXHTMLContent{depth: utils.NewDepthWatcher(), completed: false}
	i.Content = &bytes.Buffer{}
	i.Encoder = xml.NewEncoder(i.Content)
	return &i
}

func stripXHTMLNamespace(t utils.StartElement) utils.StartElement {
	for i, attr := range t.Attr {
		if attr.Name.Local == "xmlns" && attr.Value == "http://www.w3.org/1999/xhtml" {
			t.Attr = append(t.Attr[:i], t.Attr[i+1:]...)
			break
		}
	}
	return t
}

func (i *InlineXHTMLContent) EncodeXHTMLToken(t xml.Token) error {
	var err error
	switch t := t.(type) {
	case utils.StartElement:
		t = stripXHTMLNamespace(t)
		t.Name.Space = ""
		err = i.Encoder.EncodeToken(*t.StartElement)
	case xml.EndElement:
		t.Name.Space = ""
		err = i.Encoder.EncodeToken(t)

	}
	return err
}

func (i *InlineXHTMLContent) CheckXHTMLSpace(el utils.StartElement) utils.ParserError {
	if !el.Ns.Has("http://www.w3.org/1999/xhtml") {
		return utils.NewError(XHTMLElementNotNamespaced, fmt.Sprintf("'%s' element is not in XHTML namespace (ns => '%s')", el.Name.Local, el.Name.Space))
	}
	return nil
}

func (i *InlineXHTMLContent) ProcessStartElement(el utils.StartElement) (utils.Visitor, utils.ParserError) {
	err := errors.NewErrorAggregator()
	if i.depth.Level == 0 && el.Name.Local != "div" {
		err.NewError(utils.NewError(XHTMLRootNodeNotDiv, "Inline XHTML root node must be a div"))
	}

	if i.completed == true {
		err.NewError(utils.NewError(NotUniqueChild, "Inline XHTML should be contained in a unique node"))
	}

	i.depth.Down()
	if error := i.CheckXHTMLSpace(el); error != nil {
		err.NewError(error)
	}
	if error := i.EncodeXHTMLToken(el); error != nil {
		err.NewError(utils.NewError(XHTMLEncodeToStringError, "cannot encode XHTML"))
	}

	return i, err.ErrorObject()
}

func (i *InlineXHTMLContent) ProcessEndElement(el xml.EndElement) (utils.Visitor, utils.ParserError) {
	level := i.depth.Up()
	if level == utils.RootLevel {

		if err := i.EncodeXHTMLToken(el); err != nil {
			return i, utils.NewError(XHTMLEncodeToStringError, "cannot encode XHTML")
		}

		if err := i.flush(); err != nil {
			return i, err
		}

		i.completed = true

		return i, nil
	}

	if level == utils.ParentLevel {
		if err := i.flush(); err != nil {
			return i.Parent, err
		}

		if i.Parent != nil {
			return i.Parent.ProcessEndElement(el)
		}
		return nil, nil
	}

	if err := i.EncodeXHTMLToken(el); err != nil {
		return i, utils.NewError(XHTMLEncodeToStringError, "cannot encode XHTML")
	}

	return i, nil
}

func (i *InlineXHTMLContent) ProcessCharData(el xml.CharData) (utils.Visitor, utils.ParserError) {
	if len(strings.Fields(string(el))) > 0 {
		if err := i.flush(); err != nil {
			return i, err
		}

		if _, err := i.Content.Write(el); err != nil {
			return i, utils.NewError(CannotFlush, "cannot flush content")
		}

		if i.depth.Level == 0 {
			return i, utils.NewError(XHTMLRootNodeNotDiv, "XHTML element should have a root")
		}
	}

	return i, nil
}

func (i *InlineXHTMLContent) flush() utils.ParserError {
	err := i.Encoder.Flush()

	if err != nil {
		return utils.NewError(CannotFlush, "cannot flush XHTML")
	}
	return nil
}

func (i *InlineXHTMLContent) String() string {
	return string(i.Content.String())
}
