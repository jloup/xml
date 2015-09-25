package atom

import (
	"fmt"

	xmlutils "github.com/jloup/xml/utils"
)

var (
	IsAbsoluteIRI  = xmlutils.IsValidAbsoluteIri(IriNotAbsolute)
	IsValidIRI     = xmlutils.IsValidIri(IriNotValid)
	IsValidLength  = xmlutils.IsValidNumber(NotPositiveNumber)
	IsValidMIME    = xmlutils.IsValidMIME(IsNotMIME)
	IsXMLMediaType = xmlutils.IsValidXMLMediaType(NotXMLMediaType)
)

func relIsIANA(name, s string) xmlutils.ParserError {
	if s == "alternate" || s == "related" || s == "self" || s == "enclosure" || s == "via" {
		return nil
	}

	return xmlutils.NewError(RelNotValid, fmt.Sprintf("rel is not valid: %s", s))
}

func contentTypeIsValid(name, s string) xmlutils.ParserError {
	if s == "text" || s == "html" || s == "xhtml" {
		return xmlutils.NewError(ContentTypeIsNotValid, "type not valid")
	}

	return nil
}
