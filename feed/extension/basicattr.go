package extension

import (
	"encoding/xml"

	"github.com/JLoup/xml/helper"
)

type BasicAttr struct {
	Content   string
	name      xml.Name
	Validator helper.ElementValidator

	Parent helper.Visitor
}

func NewBasicAttr(name xml.Name, validator helper.ElementValidator) BasicAttr {
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

func (b *BasicAttr) SetParent(parent helper.Visitor) {
	b.Parent = parent
}

func (b *BasicAttr) Validate() helper.ParserError {

	if b.Validator == nil {
		return nil
	}
	return b.Validator(b.name.Local, b.Content)
}
