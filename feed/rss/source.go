package rss

import (
	"encoding/xml"
	"strings"

	"github.com/jloup/errors"
	"github.com/jloup/xml/feed/extension"
	"github.com/jloup/xml/utils"
)

type Source struct {
	Url     utils.Element
	Content utils.Element

	Extension extension.VisitorExtension
	Parent    utils.Visitor
	depth     utils.DepthWatcher
}

func NewSource() *Source {
	s := Source{depth: utils.NewDepthWatcher()}

	s.Url = utils.NewElement("url", "", utils.Nop)
	s.Url.SetOccurence(utils.NewOccurence("url", utils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

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

func (s *Source) ProcessStartElement(el utils.StartElement) (utils.Visitor, utils.ParserError) {
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

func (s *Source) ProcessEndElement(el xml.EndElement) (utils.Visitor, utils.ParserError) {
	if s.depth.Up() == utils.RootLevel {

		return s.Parent, s.validate()
	}

	return s, nil
}

func (s *Source) ProcessCharData(el xml.CharData) (utils.Visitor, utils.ParserError) {
	s.Content.Value = strings.TrimSpace(string(el))
	return s, nil
}

func (s *Source) validate() utils.ParserError {
	error := errors.NewErrorAggregator()

	utils.ValidateElements("source", &error, s.Url)
	s.Extension.Validate(&error)

	return error.ErrorObject()
}
