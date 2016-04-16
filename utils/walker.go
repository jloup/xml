// package utils provides tools to help walking, processing and validating XML 1.0 documents
package utils

import (
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"
	"strings"

	"golang.org/x/net/html/charset"
)

type StartElement struct {
	*xml.StartElement
	Ns *Namespaces
}

func Walk(r io.Reader, v Visitor, custom FlagChecker, xmlTokenErrorRetry int) ParserError {

	var err error
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return NewError(IOError, "Cannot read content")
	}

	for {

		b = bytes.TrimSpace(b)
		r = bytes.NewReader(b)
		dec := xml.NewDecoder(r)
		dec.CharsetReader = charset.NewReaderLabel

		var t xml.Token
		var startOffset int64
		var endOffset int64
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

				// XMLTokenError - should we retry ?
				if xmlTokenErrorRetry <= 0 {
					return Error{flag: XMLTokenError, msg: err.Error()} // we must abort
				} else {
					endOffset = dec.InputOffset()
					b = append(b[:startOffset], b[endOffset:]...)
					xmlTokenErrorRetry -= 1
					break
				}
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
					startOffset = dec.InputOffset()
					dec.Skip()
				} else {
					startOffset = dec.InputOffset()
					v = startVisitor
				}

			case xml.EndElement:
				tokenName = tt.Name.Local
				namespaces.Dec(tt.Name.Space)
				v, perr = v.ProcessEndElement(tt)
				startOffset = dec.InputOffset()
			case xml.CharData:
				v, perr = v.ProcessCharData(tt)
				startOffset = dec.InputOffset()
			}

			if perr != nil && custom.CheckFlag(tokenName, perr) {
				return &delegatedError{delegatedError: custom.ErrorWithCode(tokenName, perr), tokenName: tokenName}
			}

		}

	}
}
