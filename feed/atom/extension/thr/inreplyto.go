package thr

import (
	"encoding/xml"

	"github.com/JLoup/errors"
	"github.com/JLoup/xml/feed/atom"
	"github.com/JLoup/xml/feed/extension"
	"github.com/JLoup/xml/helper"
)

var _inreplyto = xml.Name{Space: NS, Local: "in-reply-to"}

type InReplyTo struct {
	Ref    helper.Element
	Href   helper.Element
	Type   helper.Element
	Source helper.Element

	Parent helper.Visitor
	depth  helper.DepthWatcher
}

func newInReplyTo() *InReplyTo {
	i := InReplyTo{depth: helper.NewDepthWatcher()}

	i.Ref = helper.NewElement("ref", "", atom.IsAbsoluteIRI)
	i.Ref.SetOccurence(helper.NewOccurence("ref", helper.ExistsAndUniqueValidator(atom.MissingAttribute, atom.AttributeDuplicated)))

	i.Href = helper.NewElement("href", "", atom.IsValidIRI)
	i.Href.SetOccurence(helper.NewOccurence("href", helper.UniqueValidator(atom.AttributeDuplicated)))

	i.Type = helper.NewElement("type", "", helper.Nop)
	i.Type.SetOccurence(helper.NewOccurence("type", helper.UniqueValidator(atom.AttributeDuplicated)))

	i.Source = helper.NewElement("source", "", atom.IsValidIRI)
	i.Source.SetOccurence(helper.NewOccurence("source", helper.UniqueValidator(atom.AttributeDuplicated)))

	return &i
}

func (i *InReplyTo) ProcessStartElement(el helper.StartElement) (helper.Visitor, helper.ParserError) {
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

func (i *InReplyTo) ProcessEndElement(el xml.EndElement) (helper.Visitor, helper.ParserError) {
	if i.depth.Up() == helper.RootLevel {
		return i.Parent, i.Validate()
	}

	return i, nil
}

func (i *InReplyTo) ProcessCharData(el xml.CharData) (helper.Visitor, helper.ParserError) {
	return i, nil
}

func (i *InReplyTo) Validate() helper.ParserError {
	error := errors.NewErrorAggregator()

	helper.ValidateElements("in-reply-to", &error, i.Ref, i.Href, i.Source, i.Type)

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

func (i *InReplyTo) SetParent(p helper.Visitor) {
	i.Parent = p
}
