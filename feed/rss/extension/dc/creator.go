package dc

import (
	"encoding/xml"

	"github.com/jloup/xml/feed/extension"
	"github.com/jloup/xml/feed/rss"
	xmlutils "github.com/jloup/xml/utils"
)

var CREATOR = xml.Name{Space: NS, Local: "creator"}

func NewCreatorElement() extension.Element {
	c := rss.NewBasicElement()

	c.Content = xmlutils.NewElement("creator", "", xmlutils.Nop)

	return c
}
