package atom

import (
	"encoding/xml"

	"github.com/jloup/errors"
	"github.com/jloup/xml/feed/extension"
	"github.com/jloup/xml/utils"
)

type Person struct {
	CommonAttributes
	Name  *BasicElement
	Uri   *BasicElement
	Email *BasicElement

	name      string
	Extension extension.VisitorExtension
	Parent    utils.Visitor
	depth     utils.DepthWatcher
}

func NewPerson() *Person {
	p := Person{depth: utils.NewDepthWatcher()}

	p.Name = NewBasicElement(&p)
	p.Uri = NewBasicElement(&p)
	p.Email = NewBasicElement(&p)

	p.init()

	return &p
}

func NewPersonExt(manager extension.Manager) *Person {
	p := Person{depth: utils.NewDepthWatcher()}

	p.Name = NewBasicElementExt(&p, manager)
	p.Uri = NewBasicElementExt(&p, manager)
	p.Email = NewBasicElementExt(&p, manager)

	p.Extension = extension.InitExtension("person", manager)

	p.init()

	return &p
}

func (p *Person) init() {

	p.Name.Content = utils.NewElement("name", "", utils.Nop)
	p.Name.Content.SetOccurence(utils.NewOccurence("name", utils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	p.Uri.Content = utils.NewElement("uri", "", IsValidIRI)
	p.Uri.Content.SetOccurence(utils.NewOccurence("uri", utils.UniqueValidator(AttributeDuplicated)))

	p.Email.Content = utils.NewElement("email", "", utils.Nop)
	p.Email.Content.SetOccurence(utils.NewOccurence("email", utils.UniqueValidator(AttributeDuplicated)))

	p.InitCommonAttributes()

}

func (p *Person) reset() {
	p.ResetAttr()
	p.Name.Content.Reset()
	p.Uri.Content.Reset()
	p.Email.Content.Reset()
}

func (p *Person) ProcessStartElement(el utils.StartElement) (utils.Visitor, utils.ParserError) {
	if p.depth.IsRoot() {
		p.reset()
		p.name = el.Name.Local
		p.Extension = extension.InitExtension(p.name, p.Extension.Manager)
		for _, attr := range el.Attr {
			if !p.ProcessAttr(attr) {
				p.Extension.ProcessAttr(attr, p)
			}
		}

	}

	switch el.Name.Space {
	case "", "http://www.w3.org/2005/atom":
		switch el.Name.Local {
		case "name":
			p.Name.Content.IncOccurence()
			p.Name.Reset()
			return p.Name.ProcessStartElement(el)
		case "uri":
			p.Uri.Content.IncOccurence()
			p.Uri.Reset()
			return p.Uri.ProcessStartElement(el)
		case "email":
			p.Email.Content.IncOccurence()
			p.Email.Reset()
			return p.Email.ProcessStartElement(el)
		}
	default:
		return p.Extension.ProcessElement(el, p)
	}

	p.depth.Down()

	return p, nil
}

func (p *Person) ProcessEndElement(el xml.EndElement) (utils.Visitor, utils.ParserError) {
	if p.depth.Up() == utils.RootLevel {

		return p.Parent, p.validate()
	}

	return p, nil
}

func (p *Person) ProcessCharData(el xml.CharData) (utils.Visitor, utils.ParserError) {
	return p, nil
}

func (p *Person) validate() utils.ParserError {

	error := errors.NewErrorAggregator()

	utils.ValidateOccurences(p.name, &error, p.Name.Content.Occurence, p.Uri.Content.Occurence, p.Email.Content.Occurence)
	p.ValidateCommonAttributes(p.name, &error)
	p.Extension.Validate(&error)

	return error.ErrorObject()

}
