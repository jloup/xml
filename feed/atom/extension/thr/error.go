package thr

import (
	"github.com/jloup/xml/utils"
	"github.com/jloup/flag"
)

var (
	LinkNotReplies   = flag.InitFlag(&utils.ErrorFlagCounter, "LinkNotReplies")
	NotInLinkElement = flag.InitFlag(&utils.ErrorFlagCounter, "NotInLinkElement")
)
