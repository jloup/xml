package atom

import (
	"github.com/jloup/xml/utils"
	"github.com/jloup/flag"
)

var (
	MissingAttribute = flag.InitFlag(&utils.ErrorFlagCounter, "MissingAttribute")
	MissingDate      = flag.InitFlag(&utils.ErrorFlagCounter, "MissingDate")
	MissingAuthor    = flag.InitFlag(&utils.ErrorFlagCounter, "MissingAuthor")
	MissingId        = flag.InitFlag(&utils.ErrorFlagCounter, "MissingId")
	MissingTitle     = flag.InitFlag(&utils.ErrorFlagCounter, "MissingTitle")
	MissingSummary   = flag.InitFlag(&utils.ErrorFlagCounter, "MissingSummary")
	MissingSelfLink  = flag.InitFlag(&utils.ErrorFlagCounter, "MissingSelfLink")

	AttributeDuplicated = flag.InitFlag(&utils.ErrorFlagCounter, "AttributeDuplicated")
	TitleDuplicated     = flag.InitFlag(&utils.ErrorFlagCounter, "TitleDuplicated")
	IdDuplicated        = flag.InitFlag(&utils.ErrorFlagCounter, "IdDuplicated")
	AttributeForbidden  = flag.InitFlag(&utils.ErrorFlagCounter, "AttributeForbidden")

	LeafElementHasChild = flag.InitFlag(&utils.ErrorFlagCounter, "LeafElementHasChild")
	NotUniqueChild      = flag.InitFlag(&utils.ErrorFlagCounter, "NotUniqueChild")
	IsNotMIME           = flag.InitFlag(&utils.ErrorFlagCounter, "IsNotMIME")
	IriNotValid         = flag.InitFlag(&utils.ErrorFlagCounter, "IriNotValid")
	IriNotAbsolute      = flag.InitFlag(&utils.ErrorFlagCounter, "IriNotAbsolute")
	NotXMLMediaType     = flag.InitFlag(&utils.ErrorFlagCounter, "NotXMLMediaType")

	SourcedContentElementNotEmpty = flag.InitFlag(&utils.ErrorFlagCounter, "SourcedContentElementNotEmpty")
	EntryWithIdAndDateDuplicated  = flag.InitFlag(&utils.ErrorFlagCounter, "EntryWithIdAndDateDuplicated")
	ContentTypeIsNotValid         = flag.InitFlag(&utils.ErrorFlagCounter, "ContentTypeIsNotValid")
	XHTMLElementNotNamespaced     = flag.InitFlag(&utils.ErrorFlagCounter, "XHTMLElementNotNamespaced")
	XHTMLEncodeToStringError      = flag.InitFlag(&utils.ErrorFlagCounter, "XHTMLEncodeToStringError")
	XHTMLRootNodeNotDiv           = flag.InitFlag(&utils.ErrorFlagCounter, "XHTMLRootNodeNotDiv")
	DateFormat                    = flag.InitFlag(&utils.ErrorFlagCounter, "DateFormat")
	CannotFlush                   = flag.InitFlag(&utils.ErrorFlagCounter, "CannotFlush")
	RelNotValid                   = flag.InitFlag(&utils.ErrorFlagCounter, "RelNotValid")
	LinkAlternateDuplicated       = flag.InitFlag(&utils.ErrorFlagCounter, "LinkAlternateDuplicated")
	NoContentOrAlternateLink      = flag.InitFlag(&utils.ErrorFlagCounter, "NoContentOrAlternateLink")
	NotPositiveNumber             = flag.InitFlag(&utils.ErrorFlagCounter, "NotPositiveNumber")
)
