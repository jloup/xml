package thr

import (
	"encoding/xml"

	"github.com/JLoup/xml/feed/atom"
	"github.com/JLoup/xml/feed/extension"
	"github.com/JLoup/xml/utils"
)

var _total = xml.Name{Space: NS, Local: "total"}

func newTotalElement() extension.Element {
	t := atom.NewBasicElement(nil)

	t.Content = utils.NewElement("total", "", atom.IsValidLength)

	return t
}
