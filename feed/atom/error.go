package atom

import (
	xmlutils "github.com/jloup/xml/utils"
	"github.com/jloup/utils"
)

var (
	MissingAttribute = utils.InitFlag(&xmlutils.ErrorFlagCounter, "MissingAttribute")
	MissingDate      = utils.InitFlag(&xmlutils.ErrorFlagCounter, "MissingDate")
	MissingAuthor    = utils.InitFlag(&xmlutils.ErrorFlagCounter, "MissingAuthor")
	MissingId        = utils.InitFlag(&xmlutils.ErrorFlagCounter, "MissingId")
	MissingTitle     = utils.InitFlag(&xmlutils.ErrorFlagCounter, "MissingTitle")
	MissingSummary   = utils.InitFlag(&xmlutils.ErrorFlagCounter, "MissingSummary")
	MissingSelfLink  = utils.InitFlag(&xmlutils.ErrorFlagCounter, "MissingSelfLink")

	AttributeDuplicated = utils.InitFlag(&xmlutils.ErrorFlagCounter, "AttributeDuplicated")
	TitleDuplicated     = utils.InitFlag(&xmlutils.ErrorFlagCounter, "TitleDuplicated")
	IdDuplicated        = utils.InitFlag(&xmlutils.ErrorFlagCounter, "IdDuplicated")
	AttributeForbidden  = utils.InitFlag(&xmlutils.ErrorFlagCounter, "AttributeForbidden")

	LeafElementHasChild = utils.InitFlag(&xmlutils.ErrorFlagCounter, "LeafElementHasChild")
	NotUniqueChild      = utils.InitFlag(&xmlutils.ErrorFlagCounter, "NotUniqueChild")
	IsNotMIME           = utils.InitFlag(&xmlutils.ErrorFlagCounter, "IsNotMIME")
	IriNotValid         = utils.InitFlag(&xmlutils.ErrorFlagCounter, "IriNotValid")
	IriNotAbsolute      = utils.InitFlag(&xmlutils.ErrorFlagCounter, "IriNotAbsolute")
	NotXMLMediaType     = utils.InitFlag(&xmlutils.ErrorFlagCounter, "NotXMLMediaType")

	SourcedContentElementNotEmpty = utils.InitFlag(&xmlutils.ErrorFlagCounter, "SourcedContentElementNotEmpty")
	EntryWithIdAndDateDuplicated  = utils.InitFlag(&xmlutils.ErrorFlagCounter, "EntryWithIdAndDateDuplicated")
	ContentTypeIsNotValid         = utils.InitFlag(&xmlutils.ErrorFlagCounter, "ContentTypeIsNotValid")
	XHTMLElementNotNamespaced     = utils.InitFlag(&xmlutils.ErrorFlagCounter, "XHTMLElementNotNamespaced")
	XHTMLEncodeToStringError      = utils.InitFlag(&xmlutils.ErrorFlagCounter, "XHTMLEncodeToStringError")
	XHTMLRootNodeNotDiv           = utils.InitFlag(&xmlutils.ErrorFlagCounter, "XHTMLRootNodeNotDiv")
	DateFormat                    = utils.InitFlag(&xmlutils.ErrorFlagCounter, "DateFormat")
	CannotFlush                   = utils.InitFlag(&xmlutils.ErrorFlagCounter, "CannotFlush")
	RelNotValid                   = utils.InitFlag(&xmlutils.ErrorFlagCounter, "RelNotValid")
	LinkAlternateDuplicated       = utils.InitFlag(&xmlutils.ErrorFlagCounter, "LinkAlternateDuplicated")
	NoContentOrAlternateLink      = utils.InitFlag(&xmlutils.ErrorFlagCounter, "NoContentOrAlternateLink")
	NotPositiveNumber             = utils.InitFlag(&xmlutils.ErrorFlagCounter, "NotPositiveNumber")
)
