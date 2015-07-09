package thr

import (
	"encoding/xml"

	"github.com/jloup/errors"
	"github.com/jloup/xml/feed/atom"
	"github.com/jloup/xml/feed/extension"
	"github.com/jloup/xml/utils"
)

var _inreplyto = xml.Name{Space: NS, Local: "in-reply-to"}

type InReplyTo struct {
	Ref    utils.Element
	Href   utils.Element
	Type   utils.Element
	Source utils.Element

	Parent utils.Visitor
	depth  utils.DepthWatcher
}

func newInReplyTo() *InReplyTo {
	i := InReplyTo{depth: utils.NewDepthWatcher()}

	i.Ref = utils.NewElement("ref", "", atom.IsAbsoluteIRI)
	i.Ref.SetOccurence(utils.NewOccurence("ref", utils.ExistsAndUniqueValidator(atom.MissingAttribute, atom.AttributeDuplicated)))

	i.Href = utils.NewElement("href", "", atom.IsValidIRI)
	i.Href.SetOccurence(utils.NewOccurence("href", utils.UniqueValidator(atom.AttributeDuplicated)))

	i.Type = utils.NewElement("type", "", utils.Nop)
	i.Type.SetOccurence(utils.NewOccurence("type", utils.UniqueValidator(atom.AttributeDuplicated)))

	i.Source = utils.NewElement("source", "", atom.IsValidIRI)
	i.Source.SetOccurence(utils.NewOccurence("source", utils.UniqueValidator(atom.AttributeDuplicated)))

	return &i
}

func (i *InReplyTo) ProcessStartElement(el utils.StartElement) (utils.Visitor, utils.ParserError) {
	if i.depth.Level == 0 {
		for _, attr := range el.Attr {
			switch attr.Name.Space {
			case "":
				switch attr.Name.Local {
				case "ref":
					i.Ref.Value = attr.Value
					i.Ref.IncOccurence()

				case "href":
					i.Href.Value = attr.Value
					i.Href.IncOccurence()

				case "source":
					i.Source.Value = attr.Value
					i.Source.IncOccurence()

				case "type":
					i.Type.Value = attr.Value
					i.Type.IncOccurence()

				}
			}
		}
	}

	i.depth.Down()

	return i, nil
}

func (i *InReplyTo) ProcessEndElement(el xml.EndElement) (utils.Visitor, utils.ParserError) {
	if i.depth.Up() == utils.RootLevel {
		return i.Parent, i.Validate()
	}

	return i, nil
}

func (i *InReplyTo) ProcessCharData(el xml.CharData) (utils.Visitor, utils.ParserError) {
	return i, nil
}

func (i *InReplyTo) Validate() utils.ParserError {
	error := errors.NewErrorAggregator()

	utils.ValidateElements("in-reply-to", &error, i.Ref, i.Href, i.Source, i.Type)

	return error.ErrorObject()
}

func (i *InReplyTo) Name() xml.Name {
	return _inreplyto
}

func (i *InReplyTo) String() string {
	return i.Ref.Value
}

func newInReplyToElement() extension.Element {
	return newInReplyTo()
}

func (i *InReplyTo) SetParent(p utils.Visitor) {
	i.Parent = p
}
