package rss

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/JLoup/errors"
	"github.com/JLoup/xml/feed/extension"
	"github.com/JLoup/xml/helper"
)

type Date struct {
	Time       time.Time
	RawContent string

	Extension extension.VisitorExtension
	Parent    helper.Visitor
	depth     helper.DepthWatcher
}

func NewDate() *Date {
	d := Date{depth: helper.NewDepthWatcher()}

	d.depth.SetMaxDepth(1)
	return &d
}

func NewDateExt(manager extension.Manager) *Date {
	d := NewDate()
	d.depth.SetMaxDepth(1)
	d.Extension = extension.InitExtension("date", manager)

	return d
}

func (d *Date) ProcessStartElement(el helper.StartElement) (helper.Visitor, helper.ParserError) {
	if d.depth.IsRoot() {
		for _, attr := range el.Attr {
			d.Extension.ProcessAttr(attr, d)
		}
	}

	if d.depth.Down() == helper.RootLevel {
		return d.Parent, helper.NewError(LeafElementHasChild, "date construct shoud not have childs")
	}

	return d, nil
}

func (d *Date) ProcessEndElement(el xml.EndElement) (helper.Visitor, helper.ParserError) {
	return d.Parent, d.validate()
}

var rssDateFormat = []string{"Mon, 02 Jan 2006 15:04:05 MST", "Mon, _2 Jan 2006 15:04:05 -0700"}

func (d *Date) ProcessCharData(el xml.CharData) (helper.Visitor, helper.ParserError) {
	var err error

	d.RawContent = string(el)
	for _, dateFormat := range rssDateFormat {
		d.Time, err = time.Parse(dateFormat, string(el))
		if err == nil && !d.Time.IsZero() {
			return d, nil
		}
	}

	return d, helper.NewError(DateFormat, fmt.Sprintf("date not well formatted '%v'", string(el)))
}

func (d *Date) validate() helper.ParserError {
	error := errors.NewErrorAggregator()

	d.Extension.Validate(&error)

	return error.ErrorObject()
}
