package rss

import (
	"encoding/xml"

	"github.com/jloup/errors"
	"github.com/jloup/xml/feed/extension"
	"github.com/jloup/xml/utils"
)

type Enclosure struct {
	Url    utils.Element
	Length utils.Element
	Type   utils.Element

	Extension extension.VisitorExtension
	Parent    utils.Visitor
	depth     utils.DepthWatcher
}

func NewEnclosure() *Enclosure {
	e := Enclosure{depth: utils.NewDepthWatcher()}

	e.Url = utils.NewElement("url", "", utils.Nop)
	e.Url.SetOccurence(utils.NewOccurence("url", utils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	e.Length = utils.NewElement("length", "", utils.Nop)
	e.Length.SetOccurence(utils.NewOccurence("length", utils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	e.Type = utils.NewElement("type", "", utils.Nop)
	e.Type.SetOccurence(utils.NewOccurence("type", utils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

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

func (e *Enclosure) ProcessStartElement(el utils.StartElement) (utils.Visitor, utils.ParserError) {
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

func (e *Enclosure) ProcessEndElement(el xml.EndElement) (utils.Visitor, utils.ParserError) {
	if e.depth.Up() == utils.RootLevel {

		return e.Parent, e.validate()
	}

	return e, nil
}

func (e *Enclosure) ProcessCharData(el xml.CharData) (utils.Visitor, utils.ParserError) {
	return e, nil
}

func (e *Enclosure) validate() utils.ParserError {
	error := errors.NewErrorAggregator()

	utils.ValidateElements("enclosure", &error, e.Url, e.Length, e.Type)
	e.Extension.Validate(&error)

	return error.ErrorObject()
}
