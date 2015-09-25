package rss

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/jloup/utils"
	"github.com/jloup/xml/feed/extension"
	xmlutils "github.com/jloup/xml/utils"
)

type Date struct {
	Time       time.Time
	RawContent string

	Extension extension.VisitorExtension
	Parent    xmlutils.Visitor
	depth     xmlutils.DepthWatcher
}

func NewDate() *Date {
	d := Date{depth: xmlutils.NewDepthWatcher()}

	d.depth.SetMaxDepth(1)
	return &d
}

func NewDateExt(manager extension.Manager) *Date {
	d := NewDate()
	d.depth.SetMaxDepth(1)
	d.Extension = extension.InitExtension("date", manager)

	return d
}

func (d *Date) ProcessStartElement(el xmlutils.StartElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if d.depth.IsRoot() {
		for _, attr := range el.Attr {
			d.Extension.ProcessAttr(attr, d)
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

var rssDateFormat = []string{"Mon, 02 Jan 2006 15:04:05 MST", "Mon, _2 Jan 2006 15:04:05 -0700"}

func (d *Date) ProcessCharData(el xml.CharData) (xmlutils.Visitor, xmlutils.ParserError) {
	var err error

	d.RawContent = string(el)
	for _, dateFormat := range rssDateFormat {
		d.Time, err = time.Parse(dateFormat, string(el))
		if err == nil && !d.Time.IsZero() {
			return d, nil
		}
	}

	return d, xmlutils.NewError(DateFormat, fmt.Sprintf("date not well formatted '%v'", string(el)))
}

func (d *Date) validate() xmlutils.ParserError {
	error := utils.NewErrorAggregator()

	d.Extension.Validate(&error)

	return error.ErrorObject()
}
