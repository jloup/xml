package atom

import (
	"fmt"

	"github.com/JLoup/xml/helper"
)

var (
	IsAbsoluteIRI  = helper.IsValidAbsoluteIri(IriNotAbsolute)
	IsValidIRI     = helper.IsValidIri(IriNotValid)
	IsValidLength  = helper.IsValidNumber(NotPositiveNumber)
	IsValidMIME    = helper.IsValidMIME(IsNotMIME)
	IsXMLMediaType = helper.IsValidXMLMediaType(NotXMLMediaType)
)

func relIsIANA(name, s string) helper.ParserError {
	if s == "alternate" || s == "related" || s == "self" || s == "enclosure" || s == "via" {
		return nil
	}

	return helper.NewError(RelNotValid, fmt.Sprintf("rel is not valid: %s", s))
}

func contentTypeIsValid(name, s string) helper.ParserError {
	if s == "text" || s == "html" || s == "xhtml" {
		return helper.NewError(ContentTypeIsNotValid, "type not valid")
	}

	return nil
}
