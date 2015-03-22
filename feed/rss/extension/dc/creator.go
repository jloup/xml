package dc

import (
	"encoding/xml"

	"github.com/JLoup/xml/feed/extension"
	"github.com/JLoup/xml/feed/rss"
	"github.com/JLoup/xml/helper"
)

var CREATOR = xml.Name{Space: NS, Local: "creator"}

func NewCreatorElement() extension.Element {
	c := rss.NewBasicElement()

	c.Content = helper.NewElement("creator", "", helper.Nop)

	return c
}
