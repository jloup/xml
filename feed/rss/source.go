package rss

import (
	"encoding/xml"
	"strings"

	"github.com/JLoup/errors"
	"github.com/JLoup/xml/feed/extension"
	"github.com/JLoup/xml/helper"
)

type Source struct {
	Url     helper.Element
	Content helper.Element

	Extension extension.VisitorExtension
	Parent    helper.Visitor
	depth     helper.DepthWatcher
}

func NewSource() *Source {
	s := Source{depth: helper.NewDepthWatcher()}

	s.Url = helper.NewElement("url", "", helper.Nop)
	s.Url.SetOccurence(helper.NewOccurence("url", helper.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

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

func (s *Source) ProcessStartElement(el helper.StartElement) (helper.Visitor, helper.ParserError) {
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

func (s *Source) ProcessEndElement(el xml.EndElement) (helper.Visitor, helper.ParserError) {
	if s.depth.Up() == helper.RootLevel {

		return s.Parent, s.validate()
	}

	return s, nil
}

func (s *Source) ProcessCharData(el xml.CharData) (helper.Visitor, helper.ParserError) {
	s.Content.Value = strings.TrimSpace(string(el))
	return s, nil
}

func (s *Source) validate() helper.ParserError {
	error := errors.NewErrorAggregator()

	helper.ValidateElements("source", &error, s.Url)
	s.Extension.Validate(&error)

	return error.ErrorObject()
}
