package thr

import (
	"encoding/xml"

	"github.com/jloup/xml/feed/atom"
	"github.com/jloup/xml/feed/extension"
	xmlutils "github.com/jloup/xml/utils"
)

var _total = xml.Name{Space: NS, Local: "total"}

func newTotalElement() extension.Element {
	t := atom.NewBasicElement(nil)

	t.Content = xmlutils.NewElement("total", "", atom.IsValidLength)

	return t
}
