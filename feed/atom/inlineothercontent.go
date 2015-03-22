package atom

import (
	"bytes"
	"encoding/xml"
	"strings"

	"github.com/JLoup/errors"
	"github.com/JLoup/xml/helper"
)

type InlineOtherContent struct {
	Type    helper.Element
	Content *bytes.Buffer

	Encoder  *xml.Encoder
	hasChild bool

	Parent helper.Visitor
	depth  helper.DepthWatcher
}

func NewInlineOtherContent() *InlineOtherContent {
	i := InlineOtherContent{hasChild: false}

	i.Type = helper.NewElement("type", "", IsValidMIME)
	i.Type.SetOccurence(helper.NewOccurence("type", helper.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))
	i.Content = &bytes.Buffer{}

	i.Encoder = xml.NewEncoder(i.Content)
	return &i
}

func (i *InlineOtherContent) ProcessStartElement(el helper.StartElement) (helper.Visitor, helper.ParserError) {
	err := errors.NewErrorAggregator()

	if i.depth.IsRoot() {
		for _, attr := range el.Attr {
			switch attr.Name.Local {
			case "type":
				i.Type.Value = attr.Value
				i.Type.IncOccurence()

			}
		}

	} else {
		if error := i.Encoder.EncodeToken(el); error != nil {
			err.NewError(helper.NewError(XHTMLEncodeToStringError, "cannot encode XHTML"))
		}
		if i.depth.Level > 0 {
			i.hasChild = true
		}
	}

	i.depth.Down()

	return i, nil
}

func (i *InlineOtherContent) ProcessEndElement(el xml.EndElement) (helper.Visitor, helper.ParserError) {
	level := i.depth.Up()

	if level == helper.RootLevel {
		i.Encoder.Flush()
		return i.Parent, i.validate()
	}

	if err := i.Encoder.EncodeToken(el); err != nil {
		return i, helper.NewError(XHTMLEncodeToStringError, "cannot encode XHTML")
	}

	return i, nil
}

func (i *InlineOtherContent) ProcessCharData(el xml.CharData) (helper.Visitor, helper.ParserError) {
	if len(strings.Fields(string(el))) > 0 {
		i.Encoder.EncodeToken(el)
	}
	return i, nil
}

func (i *InlineOtherContent) validate() helper.ParserError {
	error := errors.NewErrorAggregator()

	helper.ValidateElements("inlineothercontent", &error, i.Type)

	if strings.HasPrefix(i.Type.Value, "text/") && i.hasChild {
		error.NewError(helper.NewError(LeafElementHasChild, "text content should not have child"))
	}

	return error.ErrorObject()
}

func (i *InlineOtherContent) HasReadableContent() bool {
	return strings.HasPrefix(i.Type.Value, "text/") || (IsXMLMediaType(i.Type.Name, i.Type.Value) == nil)
}

func (i *InlineOtherContent) String() string {
	return string(i.Content.String())
}
