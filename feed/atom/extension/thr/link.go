package thr

import (
	"encoding/xml"

	"github.com/JLoup/errors"
	"github.com/JLoup/xml/feed/atom"
	"github.com/JLoup/xml/feed/extension"
	"github.com/JLoup/xml/helper"
)

var _count = xml.Name{Space: NS, Local: "count"}
var _updated = xml.Name{Space: NS, Local: "updated"}

type Count struct {
	extension.BasicAttr
}

func (c *Count) Validate() helper.ParserError {
	errAgg := errors.NewErrorAggregator()
	link, ok := c.Parent.(*atom.Link)
	if !ok {
		errAgg.NewError(helper.NewError(NotInLinkElement, "count attr should be placed in link element"))

	} else {
		if link.Rel.String() != "replies" {
			errAgg.NewError(helper.NewError(LinkNotReplies, "link element should avec a 'replies' rel for this extension"))
		}
	}
	if err := c.Validator(c.Name().Local, c.Content); err != nil {
		errAgg.NewError(err)
	}

	return errAgg.ErrorObject()
}

func newCountAttr() extension.Attr {
	c := Count{extension.NewBasicAttr(_count, helper.IsValidNumber(atom.NotPositiveNumber))}

	return &c
}

type Updated struct {
	extension.BasicAttr
}

func (u *Updated) Validate() helper.ParserError {
	errAgg := errors.NewErrorAggregator()
	link, ok := u.Parent.(*atom.Link)
	if !ok {
		errAgg.NewError(helper.NewError(NotInLinkElement, "updated attr should be placed in link element"))

	} else {
		if link.Rel.String() != "replies" {
			errAgg.NewError(helper.NewError(LinkNotReplies, "link element should avec a 'replies' rel for this extension"))
		}
	}

	return errAgg.ErrorObject()
}

func newUpdatedAttr() extension.Attr {
	u := Updated{extension.NewBasicAttr(_updated, helper.Nop)}
	return &u
}
