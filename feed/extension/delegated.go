package extension

import (
	"encoding/xml"

	xmlutils "github.com/jloup/xml/utils"
)

// attribute extension interface - e.g. thr:count attribute
type Attr interface {
	Name() xml.Name
	String() string
	Set(s string)
	SetParent(parent xmlutils.Visitor)
	Validate() xmlutils.ParserError
}

type AttrConstructor func() Attr

// standalone element extension interface - e.g. thr:in-reply-to element
type Element interface {
	xmlutils.Visitor
	Name() xml.Name
	String() string
	SetParent(p xmlutils.Visitor)
	Validate() xmlutils.ParserError
}

type ElementConstructor func() Element
