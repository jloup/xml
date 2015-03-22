//Package rss implements functions to build an object from RSS document and check it against specification
package rss

import (
	"encoding/xml"

	"github.com/JLoup/errors"
	"github.com/JLoup/xml/feed/extension"
	"github.com/JLoup/xml/helper"
)

type Channel struct {
	Title          *BasicElement
	Link           *BasicElement
	Description    *UnescapedContent
	Language       *BasicElement
	Copyright      *BasicElement
	ManagingEditor *BasicElement
	Webmaster      *BasicElement
	PubDate        *Date
	LastBuildDate  *Date
	Categories     []*Category
	Generator      *BasicElement
	Docs           *BasicElement
	Cloud          *Cloud
	Ttl            *BasicElement
	Image          *Image
	Rating         *BasicElement
	SkipHours      *BasicElement
	SkipDays       *BasicElement

	Items      []*Item
	Parent     helper.Visitor
	Extension  extension.VisitorExtension
	depth      helper.DepthWatcher
	Occurences helper.OccurenceCollection
}

func NewChannel() *Channel {
	c := Channel{
		Title:          NewBasicElement(),
		Link:           NewBasicElement(),
		Description:    NewUnescapedContent(),
		Language:       NewBasicElement(),
		Copyright:      NewBasicElement(),
		ManagingEditor: NewBasicElement(),
		Webmaster:      NewBasicElement(),
		PubDate:        NewDate(),
		LastBuildDate:  NewDate(),
		Generator:      NewBasicElement(),
		Docs:           NewBasicElement(),
		Cloud:          NewCloud(),
		Ttl:            NewBasicElement(),
		Image:          NewImage(),
		Rating:         NewBasicElement(),
		SkipHours:      NewBasicElement(),
		SkipDays:       NewBasicElement(),

		depth: helper.NewDepthWatcher(),
	}

	c.init()

	return &c
}

func NewChannelExt(manager extension.Manager) *Channel {
	c := Channel{
		Title:          NewBasicElementExt(manager),
		Link:           NewBasicElementExt(manager),
		Description:    NewUnescapedContentExt(manager),
		Language:       NewBasicElementExt(manager),
		Copyright:      NewBasicElementExt(manager),
		ManagingEditor: NewBasicElementExt(manager),
		Webmaster:      NewBasicElementExt(manager),
		PubDate:        NewDateExt(manager),
		LastBuildDate:  NewDateExt(manager),
		Generator:      NewBasicElementExt(manager),
		Docs:           NewBasicElementExt(manager),
		Cloud:          NewCloudExt(manager),
		Ttl:            NewBasicElementExt(manager),
		Image:          NewImageExt(manager),
		Rating:         NewBasicElementExt(manager),
		SkipHours:      NewBasicElementExt(manager),
		SkipDays:       NewBasicElementExt(manager),

		depth: helper.NewDepthWatcher(),
	}

	c.init()
	c.Extension = extension.InitExtension("channel", manager)

	return &c
}

func (c *Channel) init() {

	c.Title.Content = helper.NewElement("title", "", helper.Nop)
	c.Link.Content = helper.NewElement("link", "", IsValidIRI)
	c.Language.Content = helper.NewElement("language", "", helper.Nop)
	c.Copyright.Content = helper.NewElement("copyright", "", helper.Nop)
	c.ManagingEditor.Content = helper.NewElement("managingeditor", "", helper.Nop)
	c.Webmaster.Content = helper.NewElement("webmaster", "", helper.Nop)
	c.Generator.Content = helper.NewElement("generator", "", helper.Nop)
	c.Docs.Content = helper.NewElement("docs", "", helper.Nop)
	c.Ttl.Content = helper.NewElement("ttl", "", helper.Nop)
	c.Rating.Content = helper.NewElement("rating", "", helper.Nop)
	c.SkipHours.Content = helper.NewElement("skiphours", "", helper.Nop)
	c.SkipDays.Content = helper.NewElement("skipdays", "", helper.Nop)

	c.Title.Parent = c
	c.Link.Parent = c
	c.Description.Parent = c
	c.Language.Parent = c
	c.Copyright.Parent = c
	c.ManagingEditor.Parent = c
	c.Webmaster.Parent = c
	c.PubDate.Parent = c
	c.LastBuildDate.Parent = c
	c.Generator.Parent = c
	c.Docs.Parent = c
	c.Cloud.Parent = c
	c.Ttl.Parent = c
	c.Image.Parent = c
	c.Rating.Parent = c
	c.SkipHours.Parent = c
	c.SkipDays.Parent = c

	c.Occurences = helper.NewOccurenceCollection(
		helper.NewOccurence("title", helper.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)),
		helper.NewOccurence("link", helper.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)),
		helper.NewOccurence("description", helper.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)),
		helper.NewOccurence("language", helper.UniqueValidator(AttributeDuplicated)),
		helper.NewOccurence("copyright", helper.UniqueValidator(AttributeDuplicated)),
		helper.NewOccurence("managingeditor", helper.UniqueValidator(AttributeDuplicated)),
		helper.NewOccurence("webmaster", helper.UniqueValidator(AttributeDuplicated)),
		helper.NewOccurence("pubdate", helper.UniqueValidator(AttributeDuplicated)),
		helper.NewOccurence("lastbuilddate", helper.UniqueValidator(AttributeDuplicated)),
		helper.NewOccurence("generator", helper.UniqueValidator(AttributeDuplicated)),
		helper.NewOccurence("docs", helper.UniqueValidator(AttributeDuplicated)),
		helper.NewOccurence("cloud", helper.UniqueValidator(AttributeDuplicated)),
		helper.NewOccurence("ttl", helper.UniqueValidator(AttributeDuplicated)),
		helper.NewOccurence("image", helper.UniqueValidator(AttributeDuplicated)),
		helper.NewOccurence("rating", helper.UniqueValidator(AttributeDuplicated)),
		helper.NewOccurence("skiphours", helper.UniqueValidator(AttributeDuplicated)),
		helper.NewOccurence("skipdays", helper.UniqueValidator(AttributeDuplicated)),
	)

}

func (c *Channel) reset() {
	c.Occurences.Reset()
}

func (c *Channel) ProcessStartElement(el helper.StartElement) (helper.Visitor, helper.ParserError) {
	if c.depth.IsRoot() {
		c.reset()
		for _, attr := range el.Attr {
			c.Extension.ProcessAttr(attr, c)
		}
	}

	switch el.Name.Space {
	case "":
		switch el.Name.Local {
		case "title":
			c.Occurences.Inc("title")
			return c.Title.ProcessStartElement(el)

		case "link":
			c.Occurences.Inc("link")
			return c.Link.ProcessStartElement(el)

		case "description":
			c.Occurences.Inc("description")
			return c.Description, nil

		case "language":
			c.Occurences.Inc("language")
			return c.Language.ProcessStartElement(el)

		case "copyright":
			c.Occurences.Inc("copyright")
			return c.Copyright.ProcessStartElement(el)

		case "managingeditor":
			c.Occurences.Inc("managingeditor")
			return c.ManagingEditor.ProcessStartElement(el)

		case "webmaster":
			c.Occurences.Inc("webmaster")
			return c.Webmaster.ProcessStartElement(el)

		case "pubdate":
			c.Occurences.Inc("pubdate")
			return c.PubDate.ProcessStartElement(el)

		case "lastbuilddate":
			c.Occurences.Inc("lastbuilddate")
			return c.LastBuildDate.ProcessStartElement(el)

		case "category":
			category := NewCategoryExt(c.Extension.Manager)
			category.Parent = c
			c.Categories = append(c.Categories, category)
			return category.ProcessStartElement(el)

		case "generator":
			c.Occurences.Inc("generator")
			return c.Generator.ProcessStartElement(el)

		case "docs":
			c.Occurences.Inc("docs")
			return c.Docs.ProcessStartElement(el)

		case "cloud":
			c.Occurences.Inc("cloud")
			return c.Cloud.ProcessStartElement(el)

		case "ttl":
			c.Occurences.Inc("ttl")
			return c.Ttl.ProcessStartElement(el)

		case "image":
			c.Occurences.Inc("image")
			return c.Image.ProcessStartElement(el)

		case "rating":
			c.Occurences.Inc("rating")
			return c.Rating.ProcessStartElement(el)

		case "skiphours":
			c.Occurences.Inc("skiphours")
			return c.SkipHours.ProcessStartElement(el)

		case "skipdays":
			c.Occurences.Inc("skipdays")
			return c.SkipDays.ProcessStartElement(el)
		case "item":
			item := NewItemExt(c.Extension.Manager)
			item.Parent = c
			c.Items = append(c.Items, item)
			return item.ProcessStartElement(el)
		}
	default:
		return c.Extension.ProcessElement(el, c)
	}

	c.depth.Down()

	return c, nil
}

func (c *Channel) ProcessEndElement(el xml.EndElement) (helper.Visitor, helper.ParserError) {
	if c.depth.Up() == helper.RootLevel {
		return c.Parent, c.validate()
	}

	return c, nil
}

func (c *Channel) ProcessCharData(el xml.CharData) (helper.Visitor, helper.ParserError) {
	return c, nil
}

func (c *Channel) validate() helper.ParserError {
	err := errors.NewErrorAggregator()

	helper.ValidateOccurenceCollection("channel", &err, c.Occurences)
	c.Extension.Validate(&err)

	return err.ErrorObject()
}
