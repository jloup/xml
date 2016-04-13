package youtube

import (
	"encoding/xml"

	"github.com/jloup/xml/feed/atom"
	"github.com/jloup/xml/feed/extension"
	xmlutils "github.com/jloup/xml/utils"
)

var _channelId = xml.Name{Space: NS, Local: "channelid"}

func newChannelIdElement() extension.Element {
	channelId := atom.NewBasicElement(nil)

	channelId.Content = xmlutils.NewElement("channelId", "", xmlutils.Nop)

	return channelId
}
