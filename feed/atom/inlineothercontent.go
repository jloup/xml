package atom

import (
	"bytes"
	"encoding/xml"
	"strings"

	"github.com/JLoup/errors"
	"github.com/JLoup/xml/utils"
)

type InlineOtherContent struct {
	Type    utils.Element
	Content *bytes.Buffer

	Encoder  *xml.Encoder
	hasChild bool

	Parent utils.Visitor
	depth  utils.DepthWatcher
}

func NewInlineOtherContent() *InlineOtherContent {
	i := InlineOtherContent{hasChild: false}

	i.Type = utils.NewElement("type", "", IsValidMIME)
	i.Type.SetOccurence(utils.NewOccurence("type", utils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))
	i.Content = &bytes.Buffer{}

	i.Encoder = xml.NewEncoder(i.Content)
	return &i
}

func (i *InlineOtherContent) ProcessStartElement(el utils.StartElement) (utils.Visitor, utils.ParserError) {
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
			err.NewError(utils.NewError(XHTMLEncodeToStringError, "cannot encode XHTML"))
		}
		if i.depth.Level > 0 {
			i.hasChild = true
		}
	}

	i.depth.Down()

	return i, nil
}

func (i *InlineOtherContent) ProcessEndElement(el xml.EndElement) (utils.Visitor, utils.ParserError) {
	level := i.depth.Up()

	if level == utils.RootLevel {
		i.Encoder.Flush()
		return i.Parent, i.validate()
	}

	if err := i.Encoder.EncodeToken(el); err != nil {
		return i, utils.NewError(XHTMLEncodeToStringError, "cannot encode XHTML")
	}

	return i, nil
}

func (i *InlineOtherContent) ProcessCharData(el xml.CharData) (utils.Visitor, utils.ParserError) {
	if len(strings.Fields(string(el))) > 0 {
		i.Encoder.EncodeToken(el)
	}
	return i, nil
}

func (i *InlineOtherContent) validate() utils.ParserError {
	error := errors.NewErrorAggregator()

	utils.ValidateElements("inlineothercontent", &error, i.Type)

	if strings.HasPrefix(i.Type.Value, "text/") && i.hasChild {
		error.NewError(utils.NewError(LeafElementHasChild, "text content should not have child"))
	}

	return error.ErrorObject()
}

func (i *InlineOtherContent) HasReadableContent() bool {
	return strings.HasPrefix(i.Type.Value, "text/") || (IsXMLMediaType(i.Type.Name, i.Type.Value) == nil)
}

func (i *InlineOtherContent) String() string {
	return string(i.Content.String())
}
