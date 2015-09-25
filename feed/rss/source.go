package rss

import (
	"encoding/xml"
	"strings"

	"github.com/jloup/utils"
	"github.com/jloup/xml/feed/extension"
	xmlutils "github.com/jloup/xml/utils"
)

type Source struct {
	Url     xmlutils.Element
	Content xmlutils.Element

	Extension extension.VisitorExtension
	Parent    xmlutils.Visitor
	depth     xmlutils.DepthWatcher
}

func NewSource() *Source {
	s := Source{depth: xmlutils.NewDepthWatcher()}

	s.Url = xmlutils.NewElement("url", "", xmlutils.Nop)
	s.Url.SetOccurence(xmlutils.NewOccurence("url", xmlutils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	return &s
}

func NewSourceExt(manager extension.Manager) *Source {
	s := NewSource()
	s.Extension = extension.InitExtension("source", manager)

	return s
}

func (s *Source) reset() {
	s.Url.Reset()
}

func (s *Source) ProcessStartElement(el xmlutils.StartElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if s.depth.IsRoot() {
		s.reset()
		for _, attr := range el.Attr {
			switch attr.Name.Space {
			case "":
				switch attr.Name.Local {
				case "url":
					s.Url.Value = attr.Value
					s.Url.IncOccurence()
				}
			default:
				s.Extension.ProcessAttr(attr, s)
			}
		}
	}

	s.depth.Down()

	return s, nil
}

func (s *Source) ProcessEndElement(el xml.EndElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if s.depth.Up() == xmlutils.RootLevel {

		return s.Parent, s.validate()
	}

	return s, nil
}

func (s *Source) ProcessCharData(el xml.CharData) (xmlutils.Visitor, xmlutils.ParserError) {
	s.Content.Value = strings.TrimSpace(string(el))
	return s, nil
}

func (s *Source) validate() xmlutils.ParserError {
	error := utils.NewErrorAggregator()

	xmlutils.ValidateElements("source", &error, s.Url)
	s.Extension.Validate(&error)

	return error.ErrorObject()
}
