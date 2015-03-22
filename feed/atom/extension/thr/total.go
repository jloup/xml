package thr

import (
	"encoding/xml"

	"github.com/JLoup/xml/feed/atom"
	"github.com/JLoup/xml/feed/extension"
	"github.com/JLoup/xml/helper"
)

var _total = xml.Name{Space: NS, Local: "total"}

func newTotalElement() extension.Element {
	t := atom.NewBasicElement(nil)

	t.Content = helper.NewElement("total", "", atom.IsValidLength)

	return t
}
