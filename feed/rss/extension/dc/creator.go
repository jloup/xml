package dc

import (
	"encoding/xml"

	"github.com/JLoup/xml/feed/extension"
	"github.com/JLoup/xml/feed/rss"
	"github.com/JLoup/xml/utils"
)

var CREATOR = xml.Name{Space: NS, Local: "creator"}

func NewCreatorElement() extension.Element {
	c := rss.NewBasicElement()

	c.Content = utils.NewElement("creator", "", utils.Nop)

	return c
}
