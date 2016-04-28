package youtube

import (
	"encoding/xml"

	"github.com/jloup/xml/feed/atom"
	"github.com/jloup/xml/feed/extension"
	xmlutils "github.com/jloup/xml/utils"
)

var _playlistId = xml.Name{Space: NS, Local: "playlistid"}

func newPlaylistIdElement() extension.Element {
	playlistId := atom.NewBasicElement(nil)

	playlistId.Content = xmlutils.NewElement("playlistId", "", xmlutils.Nop)

	return playlistId
}
