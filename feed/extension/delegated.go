package extension

import (
	"encoding/xml"

	"github.com/JLoup/xml/helper"
)

// attribute extension interface - e.g. thr:count attribute
type Attr interface {
	Name() xml.Name
	String() string
	Set(s string)
	SetParent(parent helper.Visitor)
	Validate() helper.ParserError
}

type AttrConstructor func() Attr

// standalone element extension interface - e.g. thr:in-reply-to element
type Element interface {
	helper.Visitor
	Name() xml.Name
	String() string
	SetParent(p helper.Visitor)
	Validate() helper.ParserError
}

type ElementConstructor func() Element
