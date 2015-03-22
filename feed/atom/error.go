package atom

import (
	"github.com/JLoup/xml/helper"
	"github.com/JLoup/flag"
)

var (
	MissingAttribute = flag.InitFlag(&helper.ErrorFlagCounter, "MissingAttribute")
	MissingDate      = flag.InitFlag(&helper.ErrorFlagCounter, "MissingDate")
	MissingAuthor    = flag.InitFlag(&helper.ErrorFlagCounter, "MissingAuthor")
	MissingId        = flag.InitFlag(&helper.ErrorFlagCounter, "MissingId")
	MissingTitle     = flag.InitFlag(&helper.ErrorFlagCounter, "MissingTitle")
	MissingSummary   = flag.InitFlag(&helper.ErrorFlagCounter, "MissingSummary")
	MissingSelfLink  = flag.InitFlag(&helper.ErrorFlagCounter, "MissingSelfLink")

	AttributeDuplicated = flag.InitFlag(&helper.ErrorFlagCounter, "AttributeDuplicated")
	TitleDuplicated     = flag.InitFlag(&helper.ErrorFlagCounter, "TitleDuplicated")
	IdDuplicated        = flag.InitFlag(&helper.ErrorFlagCounter, "IdDuplicated")
	AttributeForbidden  = flag.InitFlag(&helper.ErrorFlagCounter, "AttributeForbidden")

	LeafElementHasChild = flag.InitFlag(&helper.ErrorFlagCounter, "LeafElementHasChild")
	NotUniqueChild      = flag.InitFlag(&helper.ErrorFlagCounter, "NotUniqueChild")
	IsNotMIME           = flag.InitFlag(&helper.ErrorFlagCounter, "IsNotMIME")
	IriNotValid         = flag.InitFlag(&helper.ErrorFlagCounter, "IriNotValid")
	IriNotAbsolute      = flag.InitFlag(&helper.ErrorFlagCounter, "IriNotAbsolute")
	NotXMLMediaType     = flag.InitFlag(&helper.ErrorFlagCounter, "NotXMLMediaType")

	SourcedContentElementNotEmpty = flag.InitFlag(&helper.ErrorFlagCounter, "SourcedContentElementNotEmpty")
	EntryWithIdAndDateDuplicated  = flag.InitFlag(&helper.ErrorFlagCounter, "EntryWithIdAndDateDuplicated")
	ContentTypeIsNotValid         = flag.InitFlag(&helper.ErrorFlagCounter, "ContentTypeIsNotValid")
	XHTMLElementNotNamespaced     = flag.InitFlag(&helper.ErrorFlagCounter, "XHTMLElementNotNamespaced")
	XHTMLEncodeToStringError      = flag.InitFlag(&helper.ErrorFlagCounter, "XHTMLEncodeToStringError")
	XHTMLRootNodeNotDiv           = flag.InitFlag(&helper.ErrorFlagCounter, "XHTMLRootNodeNotDiv")
	DateFormat                    = flag.InitFlag(&helper.ErrorFlagCounter, "DateFormat")
	CannotFlush                   = flag.InitFlag(&helper.ErrorFlagCounter, "CannotFlush")
	RelNotValid                   = flag.InitFlag(&helper.ErrorFlagCounter, "RelNotValid")
	LinkAlternateDuplicated       = flag.InitFlag(&helper.ErrorFlagCounter, "LinkAlternateDuplicated")
	NoContentOrAlternateLink      = flag.InitFlag(&helper.ErrorFlagCounter, "NoContentOrAlternateLink")
	NotPositiveNumber             = flag.InitFlag(&helper.ErrorFlagCounter, "NotPositiveNumber")
)
