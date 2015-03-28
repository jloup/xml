package extension

import (
	"encoding/xml"

	"github.com/JLoup/xml/utils"
)

type BasicAttr struct {
	Content   string
	name      xml.Name
	Validator utils.ElementValidator

	Parent utils.Visitor
}

func NewBasicAttr(name xml.Name, validator utils.ElementValidator) BasicAttr {
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

func (b *BasicAttr) SetParent(parent utils.Visitor) {
	b.Parent = parent
}

func (b *BasicAttr) Validate() utils.ParserError {

	if b.Validator == nil {
		return nil
	}
	return b.Validator(b.name.Local, b.Content)
}
