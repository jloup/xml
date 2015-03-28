package atom

import (
	"fmt"

	"github.com/JLoup/xml/utils"
)

var (
	IsAbsoluteIRI  = utils.IsValidAbsoluteIri(IriNotAbsolute)
	IsValidIRI     = utils.IsValidIri(IriNotValid)
	IsValidLength  = utils.IsValidNumber(NotPositiveNumber)
	IsValidMIME    = utils.IsValidMIME(IsNotMIME)
	IsXMLMediaType = utils.IsValidXMLMediaType(NotXMLMediaType)
)

func relIsIANA(name, s string) utils.ParserError {
	if s == "alternate" || s == "related" || s == "self" || s == "enclosure" || s == "via" {
		return nil
	}

	return utils.NewError(RelNotValid, fmt.Sprintf("rel is not valid: %s", s))
}

func contentTypeIsValid(name, s string) utils.ParserError {
	if s == "text" || s == "html" || s == "xhtml" {
		return utils.NewError(ContentTypeIsNotValid, "type not valid")
	}

	return nil
}
