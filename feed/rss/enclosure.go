package rss

import (
	"encoding/xml"

	"github.com/JLoup/errors"
	"github.com/JLoup/xml/feed/extension"
	"github.com/JLoup/xml/helper"
)

type Enclosure struct {
	Url    helper.Element
	Length helper.Element
	Type   helper.Element

	Extension extension.VisitorExtension
	Parent    helper.Visitor
	depth     helper.DepthWatcher
}

func NewEnclosure() *Enclosure {
	e := Enclosure{depth: helper.NewDepthWatcher()}

	e.Url = helper.NewElement("url", "", helper.Nop)
	e.Url.SetOccurence(helper.NewOccurence("url", helper.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	e.Length = helper.NewElement("length", "", helper.Nop)
	e.Length.SetOccurence(helper.NewOccurence("length", helper.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	e.Type = helper.NewElement("type", "", helper.Nop)
	e.Type.SetOccurence(helper.NewOccurence("type", helper.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

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

func (e *Enclosure) ProcessStartElement(el helper.StartElement) (helper.Visitor, helper.ParserError) {
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

func (e *Enclosure) ProcessEndElement(el xml.EndElement) (helper.Visitor, helper.ParserError) {
	if e.depth.Up() == helper.RootLevel {

		return e.Parent, e.validate()
	}

	return e, nil
}

func (e *Enclosure) ProcessCharData(el xml.CharData) (helper.Visitor, helper.ParserError) {
	return e, nil
}

func (e *Enclosure) validate() helper.ParserError {
	error := errors.NewErrorAggregator()

	helper.ValidateElements("enclosure", &error, e.Url, e.Length, e.Type)
	e.Extension.Validate(&error)

	return error.ErrorObject()
}
