package atom

import (
	"encoding/xml"

	"github.com/jloup/utils"
	"github.com/jloup/xml/feed/extension"
	xmlutils "github.com/jloup/xml/utils"
)

type Person struct {
	CommonAttributes
	Name  *BasicElement
	Uri   *BasicElement
	Email *BasicElement

	name      string
	Extension extension.VisitorExtension
	Parent    xmlutils.Visitor
	depth     xmlutils.DepthWatcher
}

func NewPerson() *Person {
	p := Person{depth: xmlutils.NewDepthWatcher()}

	p.Name = NewBasicElement(&p)
	p.Uri = NewBasicElement(&p)
	p.Email = NewBasicElement(&p)

	p.init()

	return &p
}

func NewPersonExt(manager extension.Manager) *Person {
	p := Person{depth: xmlutils.NewDepthWatcher()}

	p.Name = NewBasicElementExt(&p, manager)
	p.Uri = NewBasicElementExt(&p, manager)
	p.Email = NewBasicElementExt(&p, manager)

	p.Extension = extension.InitExtension("person", manager)

	p.init()

	return &p
}

func (p *Person) init() {

	p.Name.Content = xmlutils.NewElement("name", "", xmlutils.Nop)
	p.Name.Content.SetOccurence(xmlutils.NewOccurence("name", xmlutils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	p.Uri.Content = xmlutils.NewElement("uri", "", IsValidIRI)
	p.Uri.Content.SetOccurence(xmlutils.NewOccurence("uri", xmlutils.UniqueValidator(AttributeDuplicated)))

	p.Email.Content = xmlutils.NewElement("email", "", xmlutils.Nop)
	p.Email.Content.SetOccurence(xmlutils.NewOccurence("email", xmlutils.UniqueValidator(AttributeDuplicated)))

	p.InitCommonAttributes()

}

func (p *Person) reset() {
	p.ResetAttr()
	p.Name.Content.Reset()
	p.Uri.Content.Reset()
	p.Email.Content.Reset()
}

func (p *Person) ProcessStartElement(el xmlutils.StartElement) (xmlutils.Visitor, xmlutils.ParserError) {
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

func (p *Person) ProcessEndElement(el xml.EndElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if p.depth.Up() == xmlutils.RootLevel {

		return p.Parent, p.validate()
	}

	return p, nil
}

func (p *Person) ProcessCharData(el xml.CharData) (xmlutils.Visitor, xmlutils.ParserError) {
	return p, nil
}

func (p *Person) validate() xmlutils.ParserError {

	error := utils.NewErrorAggregator()

	xmlutils.ValidateOccurences(p.name, &error, p.Name.Content.Occurence, p.Uri.Content.Occurence, p.Email.Content.Occurence)
	p.ValidateCommonAttributes(p.name, &error)
	p.Extension.Validate(&error)

	return error.ErrorObject()

}
