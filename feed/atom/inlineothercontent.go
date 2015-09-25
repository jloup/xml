package atom

import (
	"bytes"
	"encoding/xml"
	"strings"

	"github.com/jloup/utils"
	xmlutils "github.com/jloup/xml/utils"
)

type InlineOtherContent struct {
	Type    xmlutils.Element
	Content *bytes.Buffer

	Encoder  *xml.Encoder
	hasChild bool

	Parent xmlutils.Visitor
	depth  xmlutils.DepthWatcher
}

func NewInlineOtherContent() *InlineOtherContent {
	i := InlineOtherContent{hasChild: false}

	i.Type = xmlutils.NewElement("type", "", IsValidMIME)
	i.Type.SetOccurence(xmlutils.NewOccurence("type", xmlutils.ExistsAndUniqueValidator(MissingAttribute, AttributeDuplicated)))
	i.Content = &bytes.Buffer{}

	i.Encoder = xml.NewEncoder(i.Content)
	return &i
}

func (i *InlineOtherContent) ProcessStartElement(el xmlutils.StartElement) (xmlutils.Visitor, xmlutils.ParserError) {
	err := utils.NewErrorAggregator()

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
			err.NewError(xmlutils.NewError(XHTMLEncodeToStringError, "cannot encode XHTML"))
		}
		if i.depth.Level > 0 {
			i.hasChild = true
		}
	}

	i.depth.Down()

	return i, nil
}

func (i *InlineOtherContent) ProcessEndElement(el xml.EndElement) (xmlutils.Visitor, xmlutils.ParserError) {
	level := i.depth.Up()

	if level == xmlutils.RootLevel {
		i.Encoder.Flush()
		return i.Parent, i.validate()
	}

	if err := i.Encoder.EncodeToken(el); err != nil {
		return i, xmlutils.NewError(XHTMLEncodeToStringError, "cannot encode XHTML")
	}

	return i, nil
}

func (i *InlineOtherContent) ProcessCharData(el xml.CharData) (xmlutils.Visitor, xmlutils.ParserError) {
	if len(strings.Fields(string(el))) > 0 {
		i.Encoder.EncodeToken(el)
	}
	return i, nil
}

func (i *InlineOtherContent) validate() xmlutils.ParserError {
	error := utils.NewErrorAggregator()

	xmlutils.ValidateElements("inlineothercontent", &error, i.Type)

	if strings.HasPrefix(i.Type.Value, "text/") && i.hasChild {
		error.NewError(xmlutils.NewError(LeafElementHasChild, "text content should not have child"))
	}

	return error.ErrorObject()
}

func (i *InlineOtherContent) HasReadableContent() bool {
	return strings.HasPrefix(i.Type.Value, "text/") || (IsXMLMediaType(i.Type.Name, i.Type.Value) == nil)
}

func (i *InlineOtherContent) String() string {
	return string(i.Content.String())
}
