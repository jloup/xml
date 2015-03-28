package extension

import (
	"encoding/xml"

	"github.com/JLoup/xml/utils"
)

// attribute extension interface - e.g. thr:count attribute
type Attr interface {
	Name() xml.Name
	String() string
	Set(s string)
	SetParent(parent utils.Visitor)
	Validate() utils.ParserError
}

type AttrConstructor func() Attr

// standalone element extension interface - e.g. thr:in-reply-to element
type Element interface {
	utils.Visitor
	Name() xml.Name
	String() string
	SetParent(p utils.Visitor)
	Validate() utils.ParserError
}

type ElementConstructor func() Element
