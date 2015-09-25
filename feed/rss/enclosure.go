package rss

import (
	"encoding/xml"

	"github.com/jloup/utils"
	"github.com/jloup/xml/feed/extension"
	xmlutils "github.com/jloup/xml/utils"
)

type Enclosure struct {
	Url    xmlutils.Element
	Length xmlutils.Element
	Type   xmlutils.Element

	Extension extension.VisitorExtension
	Parent    xmlutils.Visitor
	depth     xmlutils.DepthWatcher
}

func NewEnclosure() *Enclosure {
	e := Enclosure{depth: xmlutils.NewDepthWatcher()}

	e.Url = xmlutils.NewElement("url", "", xmlutils.Nop)
	e.Url.SetOccurence(xmlutils.NewOccurence("url", xmlutils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	e.Length = xmlutils.NewElement("length", "", xmlutils.Nop)
	e.Length.SetOccurence(xmlutils.NewOccurence("length", xmlutils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	e.Type = xmlutils.NewElement("type", "", xmlutils.Nop)
	e.Type.SetOccurence(xmlutils.NewOccurence("type", xmlutils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	return &e
}

func NewEnclosureExt(manager extension.Manager) *Enclosure {
	e := NewEnclosure()
	e.Extension = extension.InitExtension("enclosure", manager)

	return e
}

func (e *Enclosure) reset() {
	e.Url.Reset()
	e.Length.Reset()
	e.Type.Reset()
}

func (e *Enclosure) ProcessStartElement(el xmlutils.StartElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if e.depth.IsRoot() {
		e.reset()
		for _, attr := range el.Attr {
			switch attr.Name.Space {
			case "":
				switch attr.Name.Local {
				case "url":
					e.Url.Value = attr.Value
					e.Url.IncOccurence()

				case "length":
					e.Length.Value = attr.Value
					e.Length.IncOccurence()

				case "type":
					e.Type.Value = attr.Value
					e.Type.IncOccurence()
				}
			default:
				e.Extension.ProcessAttr(attr, e)
			}
		}
	}

	e.depth.Down()

	return e, nil
}

func (e *Enclosure) ProcessEndElement(el xml.EndElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if e.depth.Up() == xmlutils.RootLevel {

		return e.Parent, e.validate()
	}

	return e, nil
}

func (e *Enclosure) ProcessCharData(el xml.CharData) (xmlutils.Visitor, xmlutils.ParserError) {
	return e, nil
}

func (e *Enclosure) validate() xmlutils.ParserError {
	error := utils.NewErrorAggregator()

	xmlutils.ValidateElements("enclosure", &error, e.Url, e.Length, e.Type)
	e.Extension.Validate(&error)

	return error.ErrorObject()
}
