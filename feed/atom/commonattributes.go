package atom

import (
	"encoding/xml"

	"github.com/JLoup/errors"
	"github.com/JLoup/xml/helper"
)

type CommonAttributes struct {
	Base helper.Element
	Lang helper.Element
}

func (c *CommonAttributes) InitCommonAttributes() {
	c.Base = helper.NewElement("base", "", IsValidIRI)
	c.Base.SetOccurence(helper.NewOccurence("base", helper.UniqueValidator(AttributeDuplicated)))

	c.Lang = helper.NewElement("lang", "", helper.Nop)
	c.Lang.SetOccurence(helper.NewOccurence("lang", helper.UniqueValidator(AttributeDuplicated)))

}

func (c *CommonAttributes) ProcessAttr(attr xml.Attr) bool {
	switch attr.Name.Space {
	case helper.XML_NS:
		switch attr.Name.Local {

		case "lang":
			c.Lang.Value = attr.Value
			c.Lang.IncOccurence()
			return true

		case "base":
			c.Base.Value = attr.Value
			c.Base.IncOccurence()
			return true
		}
	}

	return false
}

func (c *CommonAttributes) ResetAttr() {
	c.Base.Reset()
	c.Lang.Reset()
}

func (c *CommonAttributes) ValidateCommonAttributes(parentName string, errorAgg *errors.ErrorAggregator) {
	helper.ValidateElements(parentName, errorAgg, c.Base, c.Lang)
}
