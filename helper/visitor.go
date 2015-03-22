package helper

import (
	"encoding/xml"
)

type Visitor interface {
	ProcessStartElement(el StartElement) (Visitor, ParserError)
	ProcessEndElement(el xml.EndElement) (Visitor, ParserError)
	ProcessCharData(el xml.CharData) (Visitor, ParserError)
}
