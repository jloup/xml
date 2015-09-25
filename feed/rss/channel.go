//Package rss implements functions to build an object from RSS document and check it against specification
package rss

import (
	"encoding/xml"

	"github.com/jloup/utils"
	"github.com/jloup/xml/feed/extension"
	xmlutils "github.com/jloup/xml/utils"
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
	Parent     xmlutils.Visitor
	Extension  extension.VisitorExtension
	depth      xmlutils.DepthWatcher
	Occurences xmlutils.OccurenceCollection
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

		depth: xmlutils.NewDepthWatcher(),
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

		depth: xmlutils.NewDepthWatcher(),
	}

	c.init()
	c.Extension = extension.InitExtension("channel", manager)

	return &c
}

func (c *Channel) init() {

	c.Title.Content = xmlutils.NewElement("title", "", xmlutils.Nop)
	c.Link.Content = xmlutils.NewElement("link", "", IsValidIRI)
	c.Language.Content = xmlutils.NewElement("language", "", xmlutils.Nop)
	c.Copyright.Content = xmlutils.NewElement("copyright", "", xmlutils.Nop)
	c.ManagingEditor.Content = xmlutils.NewElement("managingeditor", "", xmlutils.Nop)
	c.Webmaster.Content = xmlutils.NewElement("webmaster", "", xmlutils.Nop)
	c.Generator.Content = xmlutils.NewElement("generator", "", xmlutils.Nop)
	c.Docs.Content = xmlutils.NewElement("docs", "", xmlutils.Nop)
	c.Ttl.Content = xmlutils.NewElement("ttl", "", xmlutils.Nop)
	c.Rating.Content = xmlutils.NewElement("rating", "", xmlutils.Nop)
	c.SkipHours.Content = xmlutils.NewElement("skiphours", "", xmlutils.Nop)
	c.SkipDays.Content = xmlutils.NewElement("skipdays", "", xmlutils.Nop)

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

	c.Occurences = xmlutils.NewOccurenceCollection(
		xmlutils.NewOccurence("title", xmlutils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)),
		xmlutils.NewOccurence("link", xmlutils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)),
		xmlutils.NewOccurence("description", xmlutils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)),
		xmlutils.NewOccurence("language", xmlutils.UniqueValidator(AttributeDuplicated)),
		xmlutils.NewOccurence("copyright", xmlutils.UniqueValidator(AttributeDuplicated)),
		xmlutils.NewOccurence("managingeditor", xmlutils.UniqueValidator(AttributeDuplicated)),
		xmlutils.NewOccurence("webmaster", xmlutils.UniqueValidator(AttributeDuplicated)),
		xmlutils.NewOccurence("pubdate", xmlutils.UniqueValidator(AttributeDuplicated)),
		xmlutils.NewOccurence("lastbuilddate", xmlutils.UniqueValidator(AttributeDuplicated)),
		xmlutils.NewOccurence("generator", xmlutils.UniqueValidator(AttributeDuplicated)),
		xmlutils.NewOccurence("docs", xmlutils.UniqueValidator(AttributeDuplicated)),
		xmlutils.NewOccurence("cloud", xmlutils.UniqueValidator(AttributeDuplicated)),
		xmlutils.NewOccurence("ttl", xmlutils.UniqueValidator(AttributeDuplicated)),
		xmlutils.NewOccurence("image", xmlutils.UniqueValidator(AttributeDuplicated)),
		xmlutils.NewOccurence("rating", xmlutils.UniqueValidator(AttributeDuplicated)),
		xmlutils.NewOccurence("skiphours", xmlutils.UniqueValidator(AttributeDuplicated)),
		xmlutils.NewOccurence("skipdays", xmlutils.UniqueValidator(AttributeDuplicated)),
	)

}

func (c *Channel) reset() {
	c.Occurences.Reset()
}

func (c *Channel) ProcessStartElement(el xmlutils.StartElement) (xmlutils.Visitor, xmlutils.ParserError) {
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

func (c *Channel) ProcessEndElement(el xml.EndElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if c.depth.Up() == xmlutils.RootLevel {
		return c.Parent, c.validate()
	}

	return c, nil
}

func (c *Channel) ProcessCharData(el xml.CharData) (xmlutils.Visitor, xmlutils.ParserError) {
	return c, nil
}

func (c *Channel) validate() xmlutils.ParserError {
	err := utils.NewErrorAggregator()

	xmlutils.ValidateOccurenceCollection("channel", &err, c.Occurences)
	c.Extension.Validate(&err)

	return err.ErrorObject()
}
