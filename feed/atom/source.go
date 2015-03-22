package atom

import (
	"encoding/xml"
	"fmt"

	"github.com/JLoup/errors"
	"github.com/JLoup/xml/feed/extension"
	"github.com/JLoup/xml/helper"
)

type Source struct {
	CommonAttributes
	Authors      []*Person
	Categories   []*Category
	Contributors []*Person
	Generator    *Generator
	Icon         *Icon
	Id           *Id
	Links        []*Link
	Logo         *Logo
	Rights       *TextConstruct
	Subtitle     *TextConstruct
	Title        *TextConstruct
	Updated      *Date

	Extension  extension.VisitorExtension
	Occurences helper.OccurenceCollection
	depth      helper.DepthWatcher
	Parent     helper.Visitor
}

func NewSource() *Source {
	s := Source{
		Generator: NewGenerator(),
		Icon:      NewIcon(),
		Id:        NewId(),
		Logo:      NewLogo(),
		Rights:    NewTextConstruct(),
		Subtitle:  NewTextConstruct(),
		Title:     NewTextConstruct(),
		Updated:   NewDate(),

		depth: helper.NewDepthWatcher(),
	}

	s.init()

	return &s
}

func NewSourceExt(manager extension.Manager) *Source {
	s := Source{
		Generator: NewGeneratorExt(manager),
		Icon:      NewIconExt(manager),
		Id:        NewIdExt(manager),
		Logo:      NewLogoExt(manager),
		Rights:    NewTextConstructExt(manager),
		Subtitle:  NewTextConstructExt(manager),
		Title:     NewTextConstructExt(manager),
		Updated:   NewDateExt(manager),

		depth: helper.NewDepthWatcher(),
	}

	s.init()
	s.Extension = extension.InitExtension("source", manager)

	return &s
}
func (s *Source) init() {

	s.Generator.Parent = s
	s.Icon.Parent = s
	s.Id.Parent = s
	s.Logo.Parent = s
	s.Rights.Parent = s
	s.Subtitle.Parent = s
	s.Title.Parent = s
	s.Updated.Parent = s

	s.Occurences = helper.NewOccurenceCollection(
		helper.NewOccurence("generator", helper.UniqueValidator(AttributeDuplicated)),
		helper.NewOccurence("icon", helper.UniqueValidator(AttributeDuplicated)),
		helper.NewOccurence("logo", helper.UniqueValidator(AttributeDuplicated)),
		helper.NewOccurence("id", helper.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)),
		helper.NewOccurence("rights", helper.UniqueValidator(AttributeDuplicated)),
		helper.NewOccurence("subtitle", helper.UniqueValidator(AttributeDuplicated)),
		helper.NewOccurence("title", helper.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)),
		helper.NewOccurence("updated", helper.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)),
	)

	s.InitCommonAttributes()
}

func (s *Source) reset() {
	s.ResetAttr()
	s.Occurences.Reset()
}

func (s *Source) ProcessStartElement(el helper.StartElement) (helper.Visitor, helper.ParserError) {
	if s.depth.IsRoot() {
		s.reset()
		for _, attr := range el.Attr {
			if !s.ProcessAttr(attr) {
				s.Extension.ProcessAttr(attr, s)
			}
		}

	}

	switch el.Name.Space {
	case "", "http://www.w3.org/2005/atom":
		switch el.Name.Local {
		case "author":
			author := NewPersonExt(s.Extension.Manager)
			author.Parent = s
			s.Authors = append(s.Authors, author)
			return author.ProcessStartElement(el)

		case "category":
			category := NewCategoryExt(s.Extension.Manager)
			category.Parent = s
			s.Categories = append(s.Categories, category)
			return category.ProcessStartElement(el)

		case "contributor":
			contributor := NewPersonExt(s.Extension.Manager)
			contributor.Parent = s
			s.Contributors = append(s.Contributors, contributor)
			return contributor.ProcessStartElement(el)

		case "generator":
			s.Occurences.Inc("generator")
			return s.Generator.ProcessStartElement(el)

		case "icon":
			s.Occurences.Inc("icon")
			return s.Icon.ProcessStartElement(el)

		case "logo":
			s.Occurences.Inc("logo")
			return s.Logo.ProcessStartElement(el)

		case "subtitle":
			s.Occurences.Inc("subtitle")
			return s.Subtitle.ProcessStartElement(el)

		case "id":
			s.Occurences.Inc("id")
			return s.Id.ProcessStartElement(el)

		case "updated":
			s.Occurences.Inc("updated")
			return s.Updated.ProcessStartElement(el)

		case "link":
			link := NewLinkExt(s.Extension.Manager)
			link.Parent = s
			s.Links = append(s.Links, link)

			return link.ProcessStartElement(el)

		case "rights":
			s.Occurences.Inc("rights")
			return s.Rights.ProcessStartElement(el)

		case "title":
			s.Occurences.Inc("title")
			return s.Title.ProcessStartElement(el)

		case "entry":
			s.depth.Down()
			return s, helper.NewError(AttributeForbidden, "source should not contain entry elements")
		}
	default:
		return s.Extension.ProcessElement(el, s)
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
	return s, nil
}

func (s *Source) hasAuthor() bool {
	return len(s.Authors) > 0
}

func (s *Source) validate() helper.ParserError {
	error := errors.NewErrorAggregator()

	helper.ValidateOccurenceCollection("source", &error, s.Occurences)
	s.Extension.Validate(&error)
	s.validateLinks(&error)
	s.ValidateCommonAttributes("source", &error)

	return error.ErrorObject()
}

func (s *Source) validateLinks(err *errors.ErrorAggregator) {
	combinations := make([]string, 0)

	for _, link := range s.Links {
		if link.Rel.Value == "alternate" {
			s := link.Type.Value + link.HrefLang.Value
			unique := true

			for _, comb := range combinations {
				if s == comb {
					err.NewError(helper.NewError(LinkAlternateDuplicated, fmt.Sprintf("Alternate Link duplicated: hreflang '%s' type '%s'", link.HrefLang.Value, link.Type.Value)))
					unique = false
				}
			}

			if unique {
				combinations = append(combinations, s)
			}
		}
	}

}
