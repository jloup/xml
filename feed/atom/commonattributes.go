package atom

import (
	"encoding/xml"

	"github.com/jloup/errors"
	"github.com/jloup/xml/utils"
)

type CommonAttributes struct {
	Base utils.Element
	Lang utils.Element
}

func (c *CommonAttributes) InitCommonAttributes() {
	c.Base = utils.NewElement("base", "", IsValidIRI)
	c.Base.SetOccurence(utils.NewOccurence("base", utils.UniqueValidator(AttributeDuplicated)))

	c.Lang = utils.NewElement("lang", "", utils.Nop)
	c.Lang.SetOccurence(utils.NewOccurence("lang", utils.UniqueValidator(AttributeDuplicated)))

}

func (c *CommonAttributes) ProcessAttr(attr xml.Attr) bool {
	switch attr.Name.Space {
	case utils.XML_NS:
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
	utils.ValidateElements(parentName, errorAgg, c.Base, c.Lang)
}
