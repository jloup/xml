//Package atom implements functions to build an object from atom document and check it against specification
package atom

import (
	"encoding/xml"
	"fmt"

	"github.com/JLoup/errors"
	"github.com/JLoup/xml/feed/extension"
	"github.com/JLoup/xml/helper"
)

type Entry struct {
	CommonAttributes
	Authors      []*Person
	Categories   []*Category
	Content      *Content
	Contributors []*Person
	Id           *Id
	Links        []*Link
	Published    *Date
	Rights       *TextConstruct
	Source       *Source
	Summary      *TextConstruct
	Title        *TextConstruct
	Updated      *Date

	Extension  extension.VisitorExtension
	Occurences helper.OccurenceCollection
	depth      helper.DepthWatcher
	Parent     helper.Visitor
}

func NewEntry() *Entry {
	e := Entry{
		Content:   NewContent(),
		Id:        NewId(),
		Published: NewDate(),
		Rights:    NewTextConstruct(),
		Source:    NewSource(),
		Summary:   NewTextConstruct(),
		Title:     NewTextConstruct(),
		Updated:   NewDate(),

		depth: helper.NewDepthWatcher(),
	}

	e.init()

	return &e
}

func NewEntryExt(manager extension.Manager) *Entry {
	e := Entry{
		Content:   NewContentExt(manager),
		Id:        NewIdExt(manager),
		Published: NewDateExt(manager),
		Rights:    NewTextConstructExt(manager),
		Source:    NewSourceExt(manager),
		Summary:   NewTextConstructExt(manager),
		Title:     NewTextConstructExt(manager),
		Updated:   NewDateExt(manager),

		depth: helper.NewDepthWatcher(),
	}

	e.init()
	e.Extension = extension.InitExtension("entry", manager)

	return &e
}

func (e *Entry) init() {

	e.Content.Parent = e
	e.Id.Parent = e
	e.Published.Parent = e
	e.Rights.Parent = e
	e.Source.Parent = e
	e.Summary.Parent = e
	e.Title.Parent = e
	e.Updated.Parent = e

	e.InitCommonAttributes()

	e.Occurences = helper.NewOccurenceCollection(
		helper.NewOccurence("content", helper.UniqueValidator(AttributeDuplicated)),
		helper.NewOccurence("id", helper.ExistsAndUniqueValidator(MissingId, IdDuplicated)),
		helper.NewOccurence("published", helper.UniqueValidator(AttributeDuplicated)),
		helper.NewOccurence("rights", helper.UniqueValidator(AttributeDuplicated)),
		helper.NewOccurence("source", helper.UniqueValidator(AttributeDuplicated)),
		helper.NewOccurence("summary", helper.UniqueValidator(AttributeDuplicated)),
		helper.NewOccurence("title", helper.ExistsAndUniqueValidator(MissingTitle, TitleDuplicated)),
		helper.NewOccurence("updated", helper.ExistsAndUniqueValidator(MissingDate, AttributeDuplicated)),
	)

}

func (e *Entry) reset() {
	e.ResetAttr()
	e.Occurences.Reset()
}

func (e *Entry) ProcessStartElement(el helper.StartElement) (helper.Visitor, helper.ParserError) {
	if e.depth.IsRoot() {
		e.reset()
		for _, attr := range el.Attr {
			if !e.ProcessAttr(attr) {
				e.Extension.ProcessAttr(attr, e)
			}
		}
	}

	switch el.Name.Space {
	case "", "http://www.w3.org/2005/atom":
		switch el.Name.Local {
		case "author":
			author := NewPersonExt(e.Extension.Manager)
			author.Parent = e
			e.Authors = append(e.Authors, author)
			return author.ProcessStartElement(el)

		case "category":
			category := NewCategoryExt(e.Extension.Manager)
			category.Parent = e
			e.Categories = append(e.Categories, category)
			return category.ProcessStartElement(el)

		case "content":
			e.Occurences.Inc("content")
			return e.Content.ProcessStartElement(el)

		case "contributor":
			contributor := NewPersonExt(e.Extension.Manager)
			contributor.Parent = e
			e.Contributors = append(e.Contributors, contributor)
			return contributor.ProcessStartElement(el)

		case "id":
			e.Occurences.Inc("id")
			return e.Id.ProcessStartElement(el)

		case "link":
			link := NewLinkExt(e.Extension.Manager)
			link.Parent = e
			e.Links = append(e.Links, link)
			return link.ProcessStartElement(el)

		case "published":
			e.Occurences.Inc("published")
			return e.Published.ProcessStartElement(el)

		case "rights":
			e.Occurences.Inc("rights")
			return e.Rights.ProcessStartElement(el)

		case "source":
			e.Occurences.Inc("source")
			return e.Source.ProcessStartElement(el)

		case "summary":
			e.Occurences.Inc("summary")
			return e.Summary.ProcessStartElement(el)

		case "title":
			e.Occurences.Inc("title")
			return e.Title.ProcessStartElement(el)

		case "updated":
			e.Occurences.Inc("updated")
			return e.Updated.ProcessStartElement(el)
		}
	default:
		return e.Extension.ProcessElement(el, e)

	}

	e.depth.Down()
	return e, nil
}

func (e *Entry) ProcessEndElement(el xml.EndElement) (helper.Visitor, helper.ParserError) {
	if e.depth.Up() == helper.RootLevel {
		return e.Parent, e.validate()
	}

	return e, nil
}

func (e *Entry) ProcessCharData(el xml.CharData) (helper.Visitor, helper.ParserError) {
	return e, nil
}

func (e *Entry) validate() helper.ParserError {
	error := errors.NewErrorAggregator()

	e.validateLinks(&error)
	e.validateAuthors(&error)
	helper.ValidateOccurenceCollection("entry", &error, e.Occurences)
	e.Extension.Validate(&error)
	e.ValidateCommonAttributes("entry", &error)

	if e.Occurences.Count("summary") == 0 && !e.Content.HasReadableContent() {
		error.NewError(helper.NewError(MissingSummary, "Summary must be provided or Content must contain XML media type, XHTML or text"))
	}

	return error.ErrorObject()
}

func (e *Entry) hasAuthor() bool {
	return len(e.Authors) > 0 || e.Source.hasAuthor()
}

func (e *Entry) validateAuthors(err *errors.ErrorAggregator) {
	if e.Parent == nil && !e.hasAuthor() {
		err.NewError(helper.NewError(MissingAuthor, "entry should contain at least one author"))
	}

}

func (e *Entry) validateLinks(err *errors.ErrorAggregator) {
	combinations := make([]string, 0)
	hasAlternateRel := false

	for _, link := range e.Links {
		if link.Rel.Value == "alternate" {
			hasAlternateRel = true
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

	if e.Occurences.Count("content") == 0 && !hasAlternateRel {
		err.NewError(helper.NewError(NoContentOrAlternateLink, "Entry should have either a Content element or a Link with alternate type"))
	}
}
