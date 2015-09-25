package extension

import (
	"encoding/xml"

	xmlutils "github.com/jloup/xml/utils"
)

type BasicAttr struct {
	Content   string
	name      xml.Name
	Validator xmlutils.ElementValidator

	Parent xmlutils.Visitor
}

func NewBasicAttr(name xml.Name, validator xmlutils.ElementValidator) BasicAttr {
	return BasicAttr{name: name, Validator: validator}
}

func (b *BasicAttr) Name() xml.Name {
	return b.name
}

func (b *BasicAttr) String() string {
	return b.Content
}

func (b *BasicAttr) Set(s string) {
	b.Content = s
}

func (b *BasicAttr) SetParent(parent xmlutils.Visitor) {
	b.Parent = parent
}

func (b *BasicAttr) Validate() xmlutils.ParserError {

	if b.Validator == nil {
		return nil
	}
	return b.Validator(b.name.Local, b.Content)
}
