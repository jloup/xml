package rss

import (
	"github.com/JLoup/xml/helper"
	"github.com/JLoup/flag"
)

var (
	LeafElementHasChild      = flag.InitFlag(&helper.ErrorFlagCounter, "LeafElementHasChild")
	MissingAttribute         = flag.InitFlag(&helper.ErrorFlagCounter, "MissingAttribute")
	AttributeDuplicated      = flag.InitFlag(&helper.ErrorFlagCounter, "AttributeDuplicated")
	XHTMLEncodeToStringError = flag.InitFlag(&helper.ErrorFlagCounter, "XHTMLEncodeToStringError")
	CannotFlush              = flag.InitFlag(&helper.ErrorFlagCounter, "CannotFlush")
	DateFormat               = flag.InitFlag(&helper.ErrorFlagCounter, "DateFormat")
	IriNotValid              = flag.InitFlag(&helper.ErrorFlagCounter, "IriNotValid")
)
