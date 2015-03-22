package thr

import (
	"github.com/JLoup/xml/helper"
	"github.com/JLoup/flag"
)

var (
	LinkNotReplies   = flag.InitFlag(&helper.ErrorFlagCounter, "LinkNotReplies")
	NotInLinkElement = flag.InitFlag(&helper.ErrorFlagCounter, "NotInLinkElement")
)
