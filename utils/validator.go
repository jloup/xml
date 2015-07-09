package utils

import (
	"fmt"
	"mime"
	"net/url"
	"strconv"
	"strings"

	"github.com/jloup/errors"
	"github.com/jloup/flag"
)

type Valider interface {
	Validate() ParserError
}

func ValidateOccurenceCollection(parentName string, agg *errors.ErrorAggregator, collection OccurenceCollection) {
	for _, occ := range collection.Occurences {
		ValidateElement(parentName, occ, agg)
	}

}

func ValidateOccurences(parentName string, agg *errors.ErrorAggregator, collection ...*Occurence) {
	for _, occ := range collection {
		ValidateElement(parentName, occ, agg)
	}

}

func ValidateElement(parentName string, el Valider, agg *errors.ErrorAggregator) {
	if err := el.Validate(); err != nil {
		agg.NewError(NewError(err.Flag(), fmt.Sprintf("%s's %s", parentName, err.Msg())))
	}
}

func ValidateElements(parentName string, agg *errors.ErrorAggregator, els ...Valider) {
	for _, el := range els {
		ValidateElement(parentName, el, agg)
	}
}

func ExistsValidator(f flag.Flag) func(*Occurence) ParserError {
	return func(o *Occurence) ParserError {
		if o.NbOccurences == 0 {
			return NewError(f, fmt.Sprintf("%s should exist", o.Name))

		}

		return nil
	}

}

func UniqueValidator(f flag.Flag) func(*Occurence) ParserError {
	return func(o *Occurence) ParserError {
		if o.NbOccurences > 1 {
			return NewError(f, fmt.Sprintf("%s should be unique", o.Name))
		}
		return nil
	}
}

func ExistsAndUniqueValidator(fe, fu flag.Flag) func(*Occurence) ParserError {
	existsV := ExistsValidator(fe)
	uniqueV := UniqueValidator(fu)
	return func(o *Occurence) ParserError {
		if err := existsV(o); err != nil {
			return err
		}
		if err := uniqueV(o); err != nil {
			return err
		}
		return nil
	}
}

func IsValidIri(f flag.Flag) func(string, string) ParserError {
	return func(name, s string) ParserError {
		_, err := url.Parse(s)

		if err != nil {
			return NewError(f, fmt.Sprintf("%s is not IRI compliant: %s", name, s))
		}

		return nil
	}
}

func IsValidAbsoluteIri(f flag.Flag) func(string, string) ParserError {
	return func(name, s string) ParserError {
		u, err := url.Parse(s)

		if err != nil {
			return NewError(f, fmt.Sprintf("%s is not IRI compliant: %s", name, s))
		}

		if !u.IsAbs() {
			return NewError(f, fmt.Sprintf("%s is not an absolute IRI: %s", name, s))
		}
		return nil
	}
}

func Nop(name, s string) ParserError {
	return nil
}

func IsValidNumber(f flag.Flag) func(string, string) ParserError {
	return func(name, s string) ParserError {
		n, err := strconv.Atoi(s)

		if err != nil {
			return NewError(f, fmt.Sprintf("Length '%s' is not valid: %v", s, err))
		}

		if n < 0 {
			return NewError(f, fmt.Sprintf("Length '%s' is not valid: it shloud be positive", s))
		}
		return nil
	}
}

func IsValidMIME(f flag.Flag) func(string, string) ParserError {
	return func(name, s string) ParserError {

		sL := strings.ToLower(s)

		if sL == "xml" {
			return NewError(f, "XML type not valid")
		}

		if _, _, err := mime.ParseMediaType(sL); err != nil {
			return NewError(f, "Not a MIME type")
		}

		if strings.HasPrefix(sL, "multipart/") || strings.HasPrefix(sL, "message/") {
			return NewError(f, "Not a MIME type")
		}
		return nil
	}
}

func IsValidXMLMediaType(f flag.Flag) func(string, string) ParserError {
	return func(name, s string) ParserError {
		sL := strings.ToLower(s)
		if strings.HasSuffix(sL, "+xml") || strings.HasSuffix(sL, "/xml") {
			return nil
		}
		return NewError(f, fmt.Sprintf("%s value '%s' an XML media type", name, s))
	}
}
