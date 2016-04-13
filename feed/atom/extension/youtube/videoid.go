package youtube

import (
	"encoding/xml"

	"github.com/jloup/xml/feed/atom"
	"github.com/jloup/xml/feed/extension"
	xmlutils "github.com/jloup/xml/utils"
)

var _videoId = xml.Name{Space: NS, Local: "videoid"}

func newVideoIdElement() extension.Element {
	videoId := atom.NewBasicElement(nil)

	videoId.Content = xmlutils.NewElement("videoId", "", xmlutils.Nop)

	return videoId
}
