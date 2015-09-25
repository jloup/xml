package atom

import (
	"encoding/xml"

	"github.com/jloup/utils"
	xmlutils "github.com/jloup/xml/utils"
)

type OutOfLineContent struct {
	Type xmlutils.Element
	Src  xmlutils.Element

	Parent xmlutils.Visitor
	depth  xmlutils.DepthWatcher
}

func NewOutOfLineContent() *OutOfLineContent {
	o := OutOfLineContent{depth: xmlutils.NewDepthWatcher()}

	o.Type = xmlutils.NewElement("type", "", outOfLineTypeIsValid)
	o.Type.SetOccurence(xmlutils.NewOccurence("type", xmlutils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	o.Src = xmlutils.NewElement("src", "", IsValidIRI)
	o.Src.SetOccurence(xmlutils.NewOccurence("src", xmlutils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))

	o.depth.SetMaxDepth(1)
	return &o
}

func outOfLineTypeIsValid(name, s string) xmlutils.ParserError {
	if err := IsValidMIME(name, s); err != nil {
		return err
	}

	if err := contentTypeIsValid(name, s); err != nil {
		return err
	}

	return nil

}

func (o *OutOfLineContent) ProcessStartElement(el xmlutils.StartElement) (xmlutils.Visitor, xmlutils.ParserError) {
	if o.depth.Down() == xmlutils.MaxDepthReached {
		return o, xmlutils.NewError(SourcedContentElementNotEmpty, "content element should be empty")
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

func (o *OutOfLineContent) ProcessEndElement(el xml.EndElement) (xmlutils.Visitor, xmlutils.ParserError) {

	if o.depth.Up() == xmlutils.RootLevel {
		return o.Parent, o.validate()
	}
	return o, nil
}

func (o *OutOfLineContent) ProcessCharData(el xml.CharData) (xmlutils.Visitor, xmlutils.ParserError) {
	return o, xmlutils.NewError(SourcedContentElementNotEmpty, "content element should be empty")
}

func (o *OutOfLineContent) validate() xmlutils.ParserError {
	error := utils.NewErrorAggregator()

	xmlutils.ValidateElements("out of line", &error, o.Type, o.Src)

	return error.ErrorObject()
}
