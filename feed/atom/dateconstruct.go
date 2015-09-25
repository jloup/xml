package atom

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/jloup/utils"
	"github.com/jloup/xml/feed/extension"
	xmlutils "github.com/jloup/xml/utils"
)

type Date struct {
	CommonAttributes
	Time time.Time

	Extension extension.VisitorExtension
	Parent    xmlutils.Visitor
	depth     xmlutils.DepthWatcher
}

func NewDate() *Date {
	d := Date{depth: xmlutils.NewDepthWatcher()}
	d.depth.SetMaxDepth(1)
	d.InitCommonAttributes()

	return &d
}

func NewDateExt(manager extension.Manager) *Date {
	d := NewDate()

	d.Extension = extension.InitExtension("date", manager)

	return d
}

func (d *Date) ProcessStartElement(el xmlutils.StartElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if d.depth.IsRoot() {
		d.ResetAttr()
		for _, attr := range el.Attr {
			if !d.ProcessAttr(attr) {
				d.Extension.ProcessAttr(attr, d)
			}
		}
	}

	if d.depth.Down() == xmlutils.RootLevel {
		return d.Parent, xmlutils.NewError(LeafElementHasChild, "date construct shoud not have childs")
	}

	return d, nil
}

func (d *Date) ProcessEndElement(el xml.EndElement) (xmlutils.Visitor, xmlutils.ParserError) {
	return d.Parent, d.validate()
}

func (d *Date) ProcessCharData(el xml.CharData) (xmlutils.Visitor, xmlutils.ParserError) {
	var err error
	d.Time, err = time.Parse(time.RFC3339, string(el))

	if err != nil || d.Time.IsZero() {
		return d, xmlutils.NewError(DateFormat, fmt.Sprintf("date not well formatted '%v'", string(el)))
	}
	return d, nil
}

func (d *Date) validate() xmlutils.ParserError {
	error := utils.NewErrorAggregator()

	d.Extension.Validate(&error)
	d.ValidateCommonAttributes("date", &error)

	return error.ErrorObject()
}
