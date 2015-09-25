package thr

import (
	xmlutils "github.com/jloup/xml/utils"
	"github.com/jloup/utils"
)

var (
	LinkNotReplies   = utils.InitFlag(&xmlutils.ErrorFlagCounter, "LinkNotReplies")
	NotInLinkElement = utils.InitFlag(&xmlutils.ErrorFlagCounter, "NotInLinkElement")
)
