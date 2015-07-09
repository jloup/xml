package thr

import (
	"encoding/xml"

	"github.com/jloup/xml/feed/atom"
	"github.com/jloup/xml/feed/extension"
	"github.com/jloup/xml/utils"
)

var _total = xml.Name{Space: NS, Local: "total"}

func newTotalElement() extension.Element {
	t := atom.NewBasicElement(nil)

	t.Content = utils.NewElement("total", "", atom.IsValidLength)

	return t
}
