package thr

import (
	"encoding/xml"

	"github.com/jloup/utils"
	"github.com/jloup/xml/feed/atom"
	"github.com/jloup/xml/feed/extension"
	xmlutils "github.com/jloup/xml/utils"
)

var _inreplyto = xml.Name{Space: NS, Local: "in-reply-to"}

type InReplyTo struct {
	Ref    xmlutils.Element
	Href   xmlutils.Element
	Type   xmlutils.Element
	Source xmlutils.Element

	Parent xmlutils.Visitor
	depth  xmlutils.DepthWatcher
}

func newInReplyTo() *InReplyTo {
	i := InReplyTo{depth: xmlutils.NewDepthWatcher()}

	i.Ref = xmlutils.NewElement("ref", "", atom.IsAbsoluteIRI)
	i.Ref.SetOccurence(xmlutils.NewOccurence("ref", xmlutils.ExistsAndUniqueValidator(atom.MissingAttribute, atom.AttributeDuplicated)))

	i.Href = xmlutils.NewElement("href", "", atom.IsValidIRI)
	i.Href.SetOccurence(xmlutils.NewOccurence("href", xmlutils.UniqueValidator(atom.AttributeDuplicated)))

	i.Type = xmlutils.NewElement("type", "", xmlutils.Nop)
	i.Type.SetOccurence(xmlutils.NewOccurence("type", xmlutils.UniqueValidator(atom.AttributeDuplicated)))

	i.Source = xmlutils.NewElement("source", "", atom.IsValidIRI)
	i.Source.SetOccurence(xmlutils.NewOccurence("source", xmlutils.UniqueValidator(atom.AttributeDuplicated)))

	return &i
}

func (i *InReplyTo) ProcessStartElement(el xmlutils.StartElement) (xmlutils.Visitor, xmlutils.ParserError) {
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

func (i *InReplyTo) ProcessEndElement(el xml.EndElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if i.depth.Up() == xmlutils.RootLevel {
		return i.Parent, i.Validate()
	}

	return i, nil
}

func (i *InReplyTo) ProcessCharData(el xml.CharData) (xmlutils.Visitor, xmlutils.ParserError) {
	return i, nil
}

func (i *InReplyTo) Validate() xmlutils.ParserError {
	error := utils.NewErrorAggregator()

	xmlutils.ValidateElements("in-reply-to", &error, i.Ref, i.Href, i.Source, i.Type)

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

func (i *InReplyTo) SetParent(p xmlutils.Visitor) {
	i.Parent = p
}
