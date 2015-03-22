package atom

import (
	"encoding/xml"

	"github.com/JLoup/errors"
	"github.com/JLoup/xml/helper"
)

type OutOfLineContent struct {
	Type helper.Element
	Src  helper.Element

	Parent helper.Visitor
	depth  helper.DepthWatcher
}

func NewOutOfLineContent() *OutOfLineContent {
	o := OutOfLineContent{depth: helper.NewDepthWatcher()}

	o.Type = helper.NewElement("type", "", outOfLineTypeIsValid)
	o.Type.SetOccurence(helper.NewOccurence("type", helper.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	o.Src = helper.NewElement("src", "", IsValidIRI)
	o.Src.SetOccurence(helper.NewOccurence("src", helper.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	o.depth.SetMaxDepth(1)
	return &o
}

func outOfLineTypeIsValid(name, s string) helper.ParserError {
	if err := IsValidMIME(name, s); err != nil {
		return err
	}

	if err := contentTypeIsValid(name, s); err != nil {
		return err
	}

	return nil

}

func (o *OutOfLineContent) ProcessStartElement(el helper.StartElement) (helper.Visitor, helper.ParserError) {
	if o.depth.Down() == helper.MaxDepthReached {
		return o, helper.NewError(SourcedContentElementNotEmpty, "content element should be empty")
	}

	for _, attr := range el.Attr {
		switch attr.Name.Space {
		case "":
			switch attr.Name.Local {
			case "type":
				o.Type.Value = attr.Value
				o.Type.IncOccurence()

			case "src":
				o.Src.Value = attr.Value
				o.Src.IncOccurence()
			}
		}
	}

	return o, nil
}

func (o *OutOfLineContent) ProcessEndElement(el xml.EndElement) (helper.Visitor, helper.ParserError) {

	if o.depth.Up() == helper.RootLevel {
		return o.Parent, o.validate()
	}
	return o, nil
}

func (o *OutOfLineContent) ProcessCharData(el xml.CharData) (helper.Visitor, helper.ParserError) {
	return o, helper.NewError(SourcedContentElementNotEmpty, "content element should be empty")
}

func (o *OutOfLineContent) validate() helper.ParserError {
	error := errors.NewErrorAggregator()

	helper.ValidateElements("out of line", &error, o.Type, o.Src)

	return error.ErrorObject()
}
