package thr

import (
	"github.com/JLoup/xml/utils"
	"github.com/JLoup/flag"
)

var (
	LinkNotReplies   = flag.InitFlag(&utils.ErrorFlagCounter, "LinkNotReplies")
	NotInLinkElement = flag.InitFlag(&utils.ErrorFlagCounter, "NotInLinkElement")
)
