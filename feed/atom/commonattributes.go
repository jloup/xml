package atom

import (
	"encoding/xml"

	"github.com/jloup/utils"
	xmlutils "github.com/jloup/xml/utils"
)

type CommonAttributes struct {
	Base xmlutils.Element
	Lang xmlutils.Element
}

func (c *CommonAttributes) InitCommonAttributes() {
	c.Base = xmlutils.NewElement("base", "", IsValidIRI)
	c.Base.SetOccurence(xmlutils.NewOccurence("base", xmlutils.UniqueValidator(AttributeDuplicated)))

	c.Lang = xmlutils.NewElement("lang", "", xmlutils.Nop)
	c.Lang.SetOccurence(xmlutils.NewOccurence("lang", xmlutils.UniqueValidator(AttributeDuplicated)))

}

func (c *CommonAttributes) ProcessAttr(attr xml.Attr) bool {
	switch attr.Name.Space {
	case xmlutils.XML_NS:
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

func (c *CommonAttributes) ValidateCommonAttributes(parentName string, errorAgg *utils.ErrorAggregator) {
	xmlutils.ValidateElements(parentName, errorAgg, c.Base, c.Lang)
}
