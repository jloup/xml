package rss

import (
	"encoding/xml"

	"github.com/jloup/errors"
	"github.com/jloup/xml/feed/extension"
	"github.com/jloup/xml/utils"
)

type Item struct {
	Title       *BasicElement
	Link        *BasicElement
	Description *UnescapedContent
	Author      *BasicElement
	Categories  []*Category
	Comments    *BasicElement
	Enclosure   *Enclosure
	Guid        *Guid
	PubDate     *Date
	Source      *Source

	Extension  extension.VisitorExtension
	Parent     utils.Visitor
	depth      utils.DepthWatcher
	Occurences utils.OccurenceCollection
}

func NewItem() *Item {
	i := Item{
		Title:       NewBasicElement(),
		Link:        NewBasicElement(),
		Description: NewUnescapedContent(),
		Author:      NewBasicElement(),
		Comments:    NewBasicElement(),
		Enclosure:   NewEnclosure(),
		Guid:        NewGuid(),
		PubDate:     NewDate(),
		Source:      NewSource(),

		depth: utils.NewDepthWatcher(),
	}

	i.init()

	return &i
}

func NewItemExt(manager extension.Manager) *Item {
	i := Item{
		Title:       NewBasicElementExt(manager),
		Link:        NewBasicElementExt(manager),
		Description: NewUnescapedContentExt(manager),
		Author:      NewBasicElementExt(manager),
		Comments:    NewBasicElementExt(manager),
		Enclosure:   NewEnclosureExt(manager),
		Guid:        NewGuidExt(manager),
		PubDate:     NewDateExt(manager),
		Source:      NewSourceExt(manager),

		depth: utils.NewDepthWatcher(),
	}

	i.init()
	i.Extension = extension.InitExtension("item", manager)

	return &i
}

func (i *Item) init() {

	i.Title.Content = utils.NewElement("title", "", utils.Nop)
	i.Link.Content = utils.NewElement("link", "", utils.Nop)
	i.Author.Content = utils.NewElement("author", "", utils.Nop)
	i.Comments.Content = utils.NewElement("comments", "", utils.Nop)

	i.Title.Parent = i
	i.Link.Parent = i
	i.Description.Parent = i
	i.Author.Parent = i
	i.Comments.Parent = i
	i.Enclosure.Parent = i
	i.Guid.Parent = i
	i.PubDate.Parent = i
	i.Source.Parent = i

	i.Occurences = utils.NewOccurenceCollection(
		utils.NewOccurence("title", utils.UniqueValidator(AttributeDuplicated)),
		utils.NewOccurence("link", utils.UniqueValidator(AttributeDuplicated)),
		utils.NewOccurence("description", utils.UniqueValidator(AttributeDuplicated)),
		utils.NewOccurence("author", utils.UniqueValidator(AttributeDuplicated)),
		utils.NewOccurence("comments", utils.UniqueValidator(AttributeDuplicated)),
		utils.NewOccurence("enclosure", utils.UniqueValidator(AttributeDuplicated)),
		utils.NewOccurence("guid", utils.UniqueValidator(AttributeDuplicated)),
		utils.NewOccurence("pubdate", utils.UniqueValidator(AttributeDuplicated)),
		utils.NewOccurence("source", utils.UniqueValidator(AttributeDuplicated)),
	)
}

func (i *Item) reset() {
	i.Occurences.Reset()
}

func (i *Item) ProcessStartElement(el utils.StartElement) (utils.Visitor, utils.ParserError) {
	if i.depth.IsRoot() {
		i.reset()
		for _, attr := range el.Attr {
			i.Extension.ProcessAttr(attr, i)
		}
	}

	switch el.Name.Space {
	case "":
		switch el.Name.Local {
		case "title":
			i.Occurences.Inc("title")
			return i.Title.ProcessStartElement(el)

		case "link":
			i.Occurences.Inc("link")
			return i.Link.ProcessStartElement(el)

		case "description":
			i.Occurences.Inc("description")
			return i.Description, nil

		case "author":
			i.Occurences.Inc("author")
			return i.Title.ProcessStartElement(el)
		case "category":
			category := NewCategoryExt(i.Extension.Manager)
			category.Parent = i
			i.Categories = append(i.Categories, category)
			return category.ProcessStartElement(el)

		case "comments":
			i.Occurences.Inc("comments")
			return i.Comments.ProcessStartElement(el)

		case "enclosure":
			i.Occurences.Inc("enclosure")
			return i.Enclosure.ProcessStartElement(el)

		case "guid":
			i.Occurences.Inc("guid")
			return i.Guid.ProcessStartElement(el)

		case "pubdate":
			i.Occurences.Inc("pubdate")
			return i.PubDate.ProcessStartElement(el)

		case "source":
			i.Occurences.Inc("source")
			return i.Source.ProcessStartElement(el)
		}
	default:
		return i.Extension.ProcessElement(el, i)
	}

	i.depth.Down()

	return i, nil
}

func (i *Item) ProcessEndElement(el xml.EndElement) (utils.Visitor, utils.ParserError) {
	if i.depth.Up() == utils.RootLevel {
		return i.Parent, i.validate()
	}

	return i, nil
}

func (i *Item) ProcessCharData(el xml.CharData) (utils.Visitor, utils.ParserError) {
	return i, nil
}

func (i *Item) validate() utils.ParserError {
	err := errors.NewErrorAggregator()

	utils.ValidateOccurenceCollection("item", &err, i.Occurences)
	i.Extension.Validate(&err)

	if i.Occurences.Count("description") == 0 && i.Occurences.Count("title") == 0 {
		err.NewError(utils.NewError(MissingAttribute, "item should have at least a title or a description"))
	}

	return err.ErrorObject()
}
