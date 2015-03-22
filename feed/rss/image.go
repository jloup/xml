package rss

import (
	"encoding/xml"

	"github.com/JLoup/errors"
	"github.com/JLoup/xml/feed/extension"
	"github.com/JLoup/xml/helper"
)

type Image struct {
	Url         *BasicElement
	Title       *BasicElement
	Link        *BasicElement
	Width       *BasicElement
	Height      *BasicElement
	Description *BasicElement

	Extension extension.VisitorExtension
	Parent    helper.Visitor
	depth     helper.DepthWatcher
}

func NewImage() *Image {
	i := Image{depth: helper.NewDepthWatcher()}

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
	i := Image{depth: helper.NewDepthWatcher()}

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
	i.Url.Content = helper.NewElement("url", "", helper.Nop)
	i.Url.Content.SetOccurence(helper.NewOccurence("url", helper.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	i.Title.Content = helper.NewElement("title", "", helper.Nop)
	i.Title.Content.SetOccurence(helper.NewOccurence("title", helper.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	i.Link.Content = helper.NewElement("link", "", helper.Nop)
	i.Link.Content.SetOccurence(helper.NewOccurence("link", helper.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	i.Width.Content = helper.NewElement("width", "", helper.Nop)
	i.Width.Content.SetOccurence(helper.NewOccurence("width", helper.UniqueValidator(AttributeDuplicated)))

	i.Height.Content = helper.NewElement("height", "", helper.Nop)
	i.Height.Content.SetOccurence(helper.NewOccurence("height", helper.UniqueValidator(AttributeDuplicated)))

	i.Description.Content = helper.NewElement("description", "", helper.Nop)
	i.Description.Content.SetOccurence(helper.NewOccurence("description", helper.UniqueValidator(AttributeDuplicated)))

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

func (i *Image) ProcessStartElement(el helper.StartElement) (helper.Visitor, helper.ParserError) {
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

func (i *Image) ProcessEndElement(el xml.EndElement) (helper.Visitor, helper.ParserError) {
	if i.depth.Up() == helper.RootLevel {

		return i.Parent, i.validate()
	}

	return i, nil
}

func (i *Image) ProcessCharData(el xml.CharData) (helper.Visitor, helper.ParserError) {
	return i, nil
}

func (i *Image) validate() helper.ParserError {
	error := errors.NewErrorAggregator()

	helper.ValidateElements("image", &error, i.Url.Content, i.Title.Content, i.Link.Content, i.Width.Content, i.Height.Content, i.Description.Content)
	i.Extension.Validate(&error)

	return error.ErrorObject()
}
