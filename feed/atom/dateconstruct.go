package atom

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/JLoup/errors"
	"github.com/JLoup/xml/feed/extension"
	"github.com/JLoup/xml/utils"
)

type Date struct {
	CommonAttributes
	Time time.Time

	Extension extension.VisitorExtension
	Parent    utils.Visitor
	depth     utils.DepthWatcher
}

func NewDate() *Date {
	d := Date{depth: utils.NewDepthWatcher()}
	d.depth.SetMaxDepth(1)
	d.InitCommonAttributes()

	return &d
}

func NewDateExt(manager extension.Manager) *Date {
	d := NewDate()

	d.Extension = extension.InitExtension("date", manager)

	return d
}

func (d *Date) ProcessStartElement(el utils.StartElement) (utils.Visitor, utils.ParserError) {
	if d.depth.IsRoot() {
		d.ResetAttr()
		for _, attr := range el.Attr {
			if !d.ProcessAttr(attr) {
				d.Extension.ProcessAttr(attr, d)
			}
		}
	}

	if d.depth.Down() == utils.RootLevel {
		return d.Parent, utils.NewError(LeafElementHasChild, "date construct shoud not have childs")
	}

	return d, nil
}

func (d *Date) ProcessEndElement(el xml.EndElement) (utils.Visitor, utils.ParserError) {
	return d.Parent, d.validate()
}

func (d *Date) ProcessCharData(el xml.CharData) (utils.Visitor, utils.ParserError) {
	var err error
	d.Time, err = time.Parse(time.RFC3339, string(el))

	if err != nil || d.Time.IsZero() {
		return d, utils.NewError(DateFormat, fmt.Sprintf("date not well formatted '%v'", string(el)))
	}
	return d, nil
}

func (d *Date) validate() utils.ParserError {
	error := errors.NewErrorAggregator()

	d.Extension.Validate(&error)
	d.ValidateCommonAttributes("date", &error)

	return error.ErrorObject()
}
