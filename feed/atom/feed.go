package atom

import (
	"encoding/xml"
	"fmt"

	"github.com/jloup/utils"
	"github.com/jloup/xml/feed/extension"
	xmlutils "github.com/jloup/xml/utils"
)

type Feed struct {
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
	Entries      []*Entry

	Extension  extension.VisitorExtension
	Occurences xmlutils.OccurenceCollection
	depth      xmlutils.DepthWatcher
	Parent     xmlutils.Visitor
}

func NewFeed() *Feed {
	f := Feed{
		Generator: NewGenerator(),
		Icon:      NewIcon(),
		Id:        NewId(),
		Logo:      NewLogo(),
		Rights:    NewTextConstruct(),
		Subtitle:  NewTextConstruct(),
		Title:     NewTextConstruct(),
		Updated:   NewDate(),

		depth: xmlutils.NewDepthWatcher(),
	}

	f.init()

	return &f
}

func NewFeedExt(manager extension.Manager) *Feed {
	f := Feed{
		Generator: NewGeneratorExt(manager),
		Icon:      NewIconExt(manager),
		Id:        NewIdExt(manager),
		Logo:      NewLogoExt(manager),
		Rights:    NewTextConstructExt(manager),
		Subtitle:  NewTextConstructExt(manager),
		Title:     NewTextConstructExt(manager),
		Updated:   NewDateExt(manager),

		depth: xmlutils.NewDepthWatcher(),
	}

	f.init()
	f.Extension = extension.InitExtension("feed", manager)

	return &f
}
func (f *Feed) init() {

	f.Generator.Parent = f
	f.Icon.Parent = f
	f.Id.Parent = f
	f.Logo.Parent = f
	f.Rights.Parent = f
	f.Subtitle.Parent = f
	f.Title.Parent = f
	f.Updated.Parent = f

	f.InitCommonAttributes()

	f.Occurences = xmlutils.NewOccurenceCollection(
		xmlutils.NewOccurence("generator", xmlutils.UniqueValidator(AttributeDuplicated)),
		xmlutils.NewOccurence("icon", xmlutils.UniqueValidator(AttributeDuplicated)),
		xmlutils.NewOccurence("logo", xmlutils.UniqueValidator(AttributeDuplicated)),
		xmlutils.NewOccurence("id", xmlutils.ExistsAndUniqueValidator(MissingId, IdDuplicated)),
		xmlutils.NewOccurence("rights", xmlutils.UniqueValidator(AttributeDuplicated)),
		xmlutils.NewOccurence("subtitle", xmlutils.UniqueValidator(AttributeDuplicated)),
		xmlutils.NewOccurence("title", xmlutils.ExistsAndUniqueValidator(MissingTitle, TitleDuplicated)),
		xmlutils.NewOccurence("updated", xmlutils.ExistsAndUniqueValidator(MissingDate, AttributeDuplicated)),
	)

}

func (f *Feed) reset() {
	f.ResetAttr()
	f.Occurences.Reset()
}

func (f *Feed) ProcessStartElement(el xmlutils.StartElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if f.depth.IsRoot() {
		f.reset()
		for _, attr := range el.Attr {
			if !f.ProcessAttr(attr) {
				f.Extension.ProcessAttr(attr, f)
			}
		}
	}

	switch el.Name.Space {
	case "", "http://www.w3.org/2005/atom":
		switch el.Name.Local {
		case "author":
			author := NewPersonExt(f.Extension.Manager)
			author.Parent = f
			f.Authors = append(f.Authors, author)
			return author.ProcessStartElement(el)

		case "category":
			category := NewCategoryExt(f.Extension.Manager)
			category.Parent = f
			f.Categories = append(f.Categories, category)
			return category.ProcessStartElement(el)

		case "contributor":
			contributor := NewPersonExt(f.Extension.Manager)
			contributor.Parent = f
			f.Contributors = append(f.Contributors, contributor)
			return contributor.ProcessStartElement(el)

		case "generator":
			f.Occurences.Inc("generator")
			return f.Generator.ProcessStartElement(el)

		case "icon":
			f.Occurences.Inc("icon")
			return f.Icon.ProcessStartElement(el)

		case "logo":
			f.Occurences.Inc("logo")
			return f.Logo.ProcessStartElement(el)

		case "subtitle":
			f.Occurences.Inc("subtitle")
			return f.Subtitle.ProcessStartElement(el)

		case "id":
			f.Occurences.Inc("id")
			return f.Id.ProcessStartElement(el)

		case "updated":
			f.Occurences.Inc("updated")
			return f.Updated.ProcessStartElement(el)

		case "entry":
			entry := NewEntryExt(f.Extension.Manager)
			entry.Parent = f
			f.Entries = append(f.Entries, entry)
			return entry.ProcessStartElement(el)

		case "link":
			link := NewLinkExt(f.Extension.Manager)
			link.Parent = f
			f.Links = append(f.Links, link)

			return link.ProcessStartElement(el)

		case "rights":
			f.Occurences.Inc("rights")
			return f.Rights.ProcessStartElement(el)

		case "title":
			f.Occurences.Inc("title")
			return f.Title.ProcessStartElement(el)
		}
	default:
		return f.Extension.ProcessElement(el, f)

	}

	f.depth.Down()
	return f, nil
}

func (f *Feed) ProcessEndElement(el xml.EndElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if f.depth.Up() == xmlutils.RootLevel {
		return f.Parent, f.validate()
	}

	return f, nil
}

func (f *Feed) ProcessCharData(el xml.CharData) (xmlutils.Visitor, xmlutils.ParserError) {
	return f, nil
}

func (f *Feed) validate() xmlutils.ParserError {
	error := utils.NewErrorAggregator()

	xmlutils.ValidateOccurenceCollection("feed", &error, f.Occurences)
	f.Extension.Validate(&error)

	f.validateLinks(&error)
	f.validateAuthors(&error)
	f.validateEntries(&error)
	f.ValidateCommonAttributes("feed", &error)

	return error.ErrorObject()
}

func (f *Feed) validateEntries(err *utils.ErrorAggregator) {
	combinations := make([]string, 0)

	for _, entry := range f.Entries {
		s := entry.Id.Content.Value + entry.Updated.Time.String()
		unique := true
		for _, comb := range combinations {

			if s == comb {
				err.NewError(xmlutils.NewError(EntryWithIdAndDateDuplicated, fmt.Sprintf("Entries are duplicated: id '%s' updated '%s'", entry.Id.Content.Value, entry.Updated.Time.String())))
				unique = false
			}
		}

		if unique {
			combinations = append(combinations, s)
		}
	}
}

func (f *Feed) validateAuthors(err *utils.ErrorAggregator) {
	if len(f.Authors) > 0 {
		return
	}

	count := 0
	for _, entry := range f.Entries {

		if !entry.hasAuthor() {
			count += 1
		}
	}

	if count > 0 || len(f.Entries) == 0 {
		err.NewError(xmlutils.NewError(MissingAuthor, fmt.Sprintf("%v entry(ies) are missing author reference", count)))
	}
}

func (f *Feed) validateLinks(err *utils.ErrorAggregator) {
	combinations := make([]string, 0)
	hasSelf := false

	for _, link := range f.Links {
		if link.Rel.Value == "alternate" {
			s := link.Type.Value + link.HrefLang.Value
			unique := true

			for _, comb := range combinations {
				if s == comb {
					err.NewError(xmlutils.NewError(LinkAlternateDuplicated, fmt.Sprintf("Alternate Link duplicated: hreflang '%s' type '%s'", link.HrefLang.Value, link.Type.Value)))
					unique = false
				}
			}

			if unique {
				combinations = append(combinations, s)
			}
		} else if link.Rel.Value == "self" {
			hasSelf = true
		}
	}

	if !hasSelf {
		err.NewError(xmlutils.NewError(MissingSelfLink, "Feed must have a link with rel attribute set to 'self'"))
	}
}
