package atom

import (
	"encoding/xml"

	"github.com/JLoup/errors"
	"github.com/JLoup/xml/utils"
)

type OutOfLineContent struct {
	Type utils.Element
	Src  utils.Element

	Parent utils.Visitor
	depth  utils.DepthWatcher
}

func NewOutOfLineContent() *OutOfLineContent {
	o := OutOfLineContent{depth: utils.NewDepthWatcher()}

	o.Type = utils.NewElement("type", "", outOfLineTypeIsValid)
	o.Type.SetOccurence(utils.NewOccurence("type", utils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	o.Src = utils.NewElement("src", "", IsValidIRI)
	o.Src.SetOccurence(utils.NewOccurence("src", utils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	o.depth.SetMaxDepth(1)
	return &o
}

func outOfLineTypeIsValid(name, s string) utils.ParserError {
	if err := IsValidMIME(name, s); err != nil {
		return err
	}

	if err := contentTypeIsValid(name, s); err != nil {
		return err
	}

	return nil

}

func (o *OutOfLineContent) ProcessStartElement(el utils.StartElement) (utils.Visitor, utils.ParserError) {
	if o.depth.Down() == utils.MaxDepthReached {
		return o, utils.NewError(SourcedContentElementNotEmpty, "content element should be empty")
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

func (o *OutOfLineContent) ProcessEndElement(el xml.EndElement) (utils.Visitor, utils.ParserError) {

	if o.depth.Up() == utils.RootLevel {
		return o.Parent, o.validate()
	}
	return o, nil
}

func (o *OutOfLineContent) ProcessCharData(el xml.CharData) (utils.Visitor, utils.ParserError) {
	return o, utils.NewError(SourcedContentElementNotEmpty, "content element should be empty")
}

func (o *OutOfLineContent) validate() utils.ParserError {
	error := errors.NewErrorAggregator()

	utils.ValidateElements("out of line", &error, o.Type, o.Src)

	return error.ErrorObject()
}
