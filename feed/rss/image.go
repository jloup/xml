package rss

import (
	"encoding/xml"

	"github.com/JLoup/errors"
	"github.com/JLoup/xml/feed/extension"
	"github.com/JLoup/xml/utils"
)

type Image struct {
	Url         *BasicElement
	Title       *BasicElement
	Link        *BasicElement
	Width       *BasicElement
	Height      *BasicElement
	Description *BasicElement

	Extension extension.VisitorExtension
	Parent    utils.Visitor
	depth     utils.DepthWatcher
}

func NewImage() *Image {
	i := Image{depth: utils.NewDepthWatcher()}

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
	i := Image{depth: utils.NewDepthWatcher()}

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
	i.Url.Content = utils.NewElement("url", "", utils.Nop)
	i.Url.Content.SetOccurence(utils.NewOccurence("url", utils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	i.Title.Content = utils.NewElement("title", "", utils.Nop)
	i.Title.Content.SetOccurence(utils.NewOccurence("title", utils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	i.Link.Content = utils.NewElement("link", "", utils.Nop)
	i.Link.Content.SetOccurence(utils.NewOccurence("link", utils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	i.Width.Content = utils.NewElement("width", "", utils.Nop)
	i.Width.Content.SetOccurence(utils.NewOccurence("width", utils.UniqueValidator(AttributeDuplicated)))

	i.Height.Content = utils.NewElement("height", "", utils.Nop)
	i.Height.Content.SetOccurence(utils.NewOccurence("height", utils.UniqueValidator(AttributeDuplicated)))

	i.Description.Content = utils.NewElement("description", "", utils.Nop)
	i.Description.Content.SetOccurence(utils.NewOccurence("description", utils.UniqueValidator(AttributeDuplicated)))

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

func (i *Image) ProcessStartElement(el utils.StartElement) (utils.Visitor, utils.ParserError) {
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

func (i *Image) ProcessEndElement(el xml.EndElement) (utils.Visitor, utils.ParserError) {
	if i.depth.Up() == utils.RootLevel {

		return i.Parent, i.validate()
	}

	return i, nil
}

func (i *Image) ProcessCharData(el xml.CharData) (utils.Visitor, utils.ParserError) {
	return i, nil
}

func (i *Image) validate() utils.ParserError {
	error := errors.NewErrorAggregator()

	utils.ValidateElements("image", &error, i.Url.Content, i.Title.Content, i.Link.Content, i.Width.Content, i.Height.Content, i.Description.Content)
	i.Extension.Validate(&error)

	return error.ErrorObject()
}
