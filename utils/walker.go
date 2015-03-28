// package utils provides tools to help walking, processing and validating XML 1.0 documents
package utils

import (
	"bytes"
	"code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"
	"encoding/xml"
	"io"
	"io/ioutil"
	"strings"
)

type StartElement struct {
	*xml.StartElement
	Ns *Namespaces
}

func Walk(r io.Reader, v Visitor, custom FlagChecker) ParserError {

	var err error
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return NewError(IOError, "Cannot read content")
	}

	b = bytes.TrimSpace(b)
	r = bytes.NewReader(b)
	dec := xml.NewDecoder(r)
	dec.CharsetReader = charset.NewReader

	var t xml.Token
	var startVisitor Visitor
	var perr ParserError
	var tokenName string
	var element StartElement
	namespaces := Namespaces{}

	for {

		if t, err = dec.Token(); err != nil {
			if err == io.EOF {
				return nil
			}

			return Error{flag: XMLTokenError, msg: err.Error()} // we must abort
		}

		switch tt := t.(type) {
		case xml.SyntaxError:
			return Error{flag: XMLSyntaxError, msg: tt.Error()} // we must abort

		case xml.Comment:
		case xml.Directive:
		case xml.ProcInst:

		case xml.StartElement:
			tokenName = tt.Name.Local
			namespaces.Inc(tt.Name.Space)

			element = StartElement{&tt, &namespaces}
			element.Name.Space = strings.ToLower(tt.Name.Space)
			element.Name.Local = strings.ToLower(tt.Name.Local)
			for i, _ := range element.Attr {
				element.Attr[i].Name.Space = strings.ToLower(element.Attr[i].Name.Space)
				element.Attr[i].Name.Local = strings.ToLower(element.Attr[i].Name.Local)
			}
			startVisitor, perr = v.ProcessStartElement(element)

			if startVisitor == nil {
				dec.Skip()
			} else {
				v = startVisitor
			}

		case xml.EndElement:
			tokenName = tt.Name.Local
			namespaces.Dec(tt.Name.Space)
			v, perr = v.ProcessEndElement(tt)
		case xml.CharData:
			v, perr = v.ProcessCharData(tt)
		}

		if perr != nil && custom.CheckFlag(tokenName, perr) {
			return &delegatedError{delegatedError: custom.ErrorWithCode(tokenName, perr), tokenName: tokenName}
		}

	}

}
