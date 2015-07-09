package thr

import (
	"encoding/xml"

	"github.com/jloup/errors"
	"github.com/jloup/xml/feed/atom"
	"github.com/jloup/xml/feed/extension"
	"github.com/jloup/xml/utils"
)

var _count = xml.Name{Space: NS, Local: "count"}
var _updated = xml.Name{Space: NS, Local: "updated"}

type Count struct {
	extension.BasicAttr
}

func (c *Count) Validate() utils.ParserError {
	errAgg := errors.NewErrorAggregator()
	link, ok := c.Parent.(*atom.Link)
	if !ok {
		errAgg.NewError(utils.NewError(NotInLinkElement, "count attr should be placed in link element"))

	} else {
		if link.Rel.String() != "replies" {
			errAgg.NewError(utils.NewError(LinkNotReplies, "link element should avec a 'replies' rel for this extension"))
		}
	}
	if err := c.Validator(c.Name().Local, c.Content); err != nil {
		errAgg.NewError(err)
	}

	return errAgg.ErrorObject()
}

func newCountAttr() extension.Attr {
	c := Count{extension.NewBasicAttr(_count, utils.IsValidNumber(atom.NotPositiveNumber))}

	return &c
}

type Updated struct {
	extension.BasicAttr
}

func (u *Updated) Validate() utils.ParserError {
	errAgg := errors.NewErrorAggregator()
	link, ok := u.Parent.(*atom.Link)
	if !ok {
		errAgg.NewError(utils.NewError(NotInLinkElement, "updated attr should be placed in link element"))

	} else {
		if link.Rel.String() != "replies" {
			errAgg.NewError(utils.NewError(LinkNotReplies, "link element should avec a 'replies' rel for this extension"))
		}
	}

	return errAgg.ErrorObject()
}

func newUpdatedAttr() extension.Attr {
	u := Updated{extension.NewBasicAttr(_updated, utils.Nop)}
	return &u
}
