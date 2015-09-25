package rss

import (
	"encoding/xml"

	"github.com/jloup/utils"
	"github.com/jloup/xml/feed/extension"
	xmlutils "github.com/jloup/xml/utils"
)

type Image struct {
	Url         *BasicElement
	Title       *BasicElement
	Link        *BasicElement
	Width       *BasicElement
	Height      *BasicElement
	Description *BasicElement

	Extension extension.VisitorExtension
	Parent    xmlutils.Visitor
	depth     xmlutils.DepthWatcher
}

func NewImage() *Image {
	i := Image{depth: xmlutils.NewDepthWatcher()}

	i.Url = NewBasicElement()
	i.Title = NewBasicElement()
	i.Link = NewBasicElement()
	i.Width = NewBasicElement()
	i.Height = NewBasicElement()
	i.Description = NewBasicElement()

	i.init()

	return &i
}

func NewImageExt(manager extension.Manager) *Image {
	i := Image{depth: xmlutils.NewDepthWatcher()}

	i.Url = NewBasicElementExt(manager)
	i.Title = NewBasicElementExt(manager)
	i.Link = NewBasicElementExt(manager)
	i.Width = NewBasicElementExt(manager)
	i.Height = NewBasicElementExt(manager)
	i.Description = NewBasicElementExt(manager)

	i.init()
	i.Extension = extension.InitExtension("image", manager)

	return &i
}

func (i *Image) init() {
	i.Url.Content = xmlutils.NewElement("url", "", xmlutils.Nop)
	i.Url.Content.SetOccurence(xmlutils.NewOccurence("url", xmlutils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	i.Title.Content = xmlutils.NewElement("title", "", xmlutils.Nop)
	i.Title.Content.SetOccurence(xmlutils.NewOccurence("title", xmlutils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	i.Link.Content = xmlutils.NewElement("link", "", xmlutils.Nop)
	i.Link.Content.SetOccurence(xmlutils.NewOccurence("link", xmlutils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	i.Width.Content = xmlutils.NewElement("width", "", xmlutils.Nop)
	i.Width.Content.SetOccurence(xmlutils.NewOccurence("width", xmlutils.UniqueValidator(AttributeDuplicated)))

	i.Height.Content = xmlutils.NewElement("height", "", xmlutils.Nop)
	i.Height.Content.SetOccurence(xmlutils.NewOccurence("height", xmlutils.UniqueValidator(AttributeDuplicated)))

	i.Description.Content = xmlutils.NewElement("description", "", xmlutils.Nop)
	i.Description.Content.SetOccurence(xmlutils.NewOccurence("description", xmlutils.UniqueValidator(AttributeDuplicated)))

	i.Url.Parent = i
	i.Title.Parent = i
	i.Link.Parent = i
	i.Width.Parent = i
	i.Height.Parent = i
	i.Description.Parent = i
}

func (i *Image) reset() {
	i.Url.Content.Reset()
	i.Title.Content.Reset()
	i.Link.Content.Reset()
	i.Width.Content.Reset()
	i.Height.Content.Reset()
	i.Description.Content.Reset()
}

func (i *Image) ProcessStartElement(el xmlutils.StartElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if i.depth.IsRoot() {
		i.reset()
		for _, attr := range el.Attr {
			i.Extension.ProcessAttr(attr, i)
		}

	}

	switch el.Name.Space {
	case "":
		switch el.Name.Local {
		case "url":
			i.Url.Content.IncOccurence()
			i.Url.Reset()
			return i.Url.ProcessStartElement(el)
		case "title":
			i.Title.Content.IncOccurence()
			i.Title.Reset()
			return i.Title.ProcessStartElement(el)
		case "link":
			i.Link.Content.IncOccurence()
			i.Link.Reset()
			return i.Link.ProcessStartElement(el)
		case "width":
			i.Width.Content.IncOccurence()
			i.Width.Reset()
			return i.Width.ProcessStartElement(el)
		case "height":
			i.Height.Content.IncOccurence()
			i.Height.Reset()
			return i.Height.ProcessStartElement(el)
		case "description":
			i.Description.Content.IncOccurence()
			i.Description.Reset()
			return i.Description.ProcessStartElement(el)
		}
	default:
		return i.Extension.ProcessElement(el, i)
	}
	i.depth.Down()

	return i, nil
}

func (i *Image) ProcessEndElement(el xml.EndElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if i.depth.Up() == xmlutils.RootLevel {

		return i.Parent, i.validate()
	}

	return i, nil
}

func (i *Image) ProcessCharData(el xml.CharData) (xmlutils.Visitor, xmlutils.ParserError) {
	return i, nil
}

func (i *Image) validate() xmlutils.ParserError {
	error := utils.NewErrorAggregator()

	xmlutils.ValidateElements("image", &error, i.Url.Content, i.Title.Content, i.Link.Content, i.Width.Content, i.Height.Content, i.Description.Content)
	i.Extension.Validate(&error)

	return error.ErrorObject()
}
