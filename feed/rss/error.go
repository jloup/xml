package rss

import (
	xmlutils "github.com/jloup/xml/utils"
	"github.com/jloup/utils"
)

var (
	LeafElementHasChild      = utils.InitFlag(&xmlutils.ErrorFlagCounter, "LeafElementHasChild")
	MissingAttribute         = utils.InitFlag(&xmlutils.ErrorFlagCounter, "MissingAttribute")
	AttributeDuplicated      = utils.InitFlag(&xmlutils.ErrorFlagCounter, "AttributeDuplicated")
	XHTMLEncodeToStringError = utils.InitFlag(&xmlutils.ErrorFlagCounter, "XHTMLEncodeToStringError")
	CannotFlush              = utils.InitFlag(&xmlutils.ErrorFlagCounter, "CannotFlush")
	DateFormat               = utils.InitFlag(&xmlutils.ErrorFlagCounter, "DateFormat")
	IriNotValid              = utils.InitFlag(&xmlutils.ErrorFlagCounter, "IriNotValid")
)
