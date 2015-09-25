package rss

import (
	xmlutils "github.com/jloup/xml/utils"
	"fmt"
	"testing"
)

func NewTestEnclosure(url, length, typ string) *Enclosure {
	e := NewEnclosure()

	e.Url.Value = url
	e.Length.Value = length
	e.Type.Value = typ

	return e
}

type testEnclosure struct {
	XML               string
	ExpectedError     xmlutils.ParserError
	ExpectedEnclosure *Enclosure
}

func testEnclosureValidator(actual xmlutils.Visitor, expected xmlutils.Visitor) error {
	e1 := actual.(*Enclosure)
	e2 := expected.(*Enclosure)

	if e1.Url.Value != e2.Url.Value {
		return fmt.Errorf("Url is invalid '%s' (expected) vs '%s'", e2.Url.Value, e1.Url.Value)
	}

	if e1.Length.Value != e2.Length.Value {
		return fmt.Errorf("Length is invalid '%s' (expected) vs '%s'", e2.Length.Value, e1.Length.Value)
	}

	if e1.Type.Value != e2.Type.Value {
		return fmt.Errorf("Type is invalid '%s' (expected) vs '%s'", e2.Type.Value, e1.Type.Value)
	}
	return nil
}

func testEnclosureConstructor() xmlutils.Visitor {
	return NewEnclosure()
}

func _TestEnclosureToTestVisitor(t testEnclosure) xmlutils.TestVisitor {
	testVisitor := xmlutils.TestVisitor{
		XML:                t.XML,
		ExpectedError:      nil,
		ExpectedVisitor:    t.ExpectedEnclosure,
		VisitorConstructor: testEnclosureConstructor,
		Validator:          testEnclosureValidator,
	}

	if t.ExpectedError != nil {
		testVisitor.ExpectedError = t.ExpectedError
	}

	return testVisitor
}
func TestEnclosureBasic(t *testing.T) {

	var testdata = []testEnclosure{
		{`<enclosure url="http://www.scripting.com/mp3s/weatherReportSuite.mp3" length="12216320" type="audio/mpeg" />`,
			nil,
			NewTestEnclosure("http://www.scripting.com/mp3s/weatherReportSuite.mp3", "12216320", "audio/mpeg"),
		},
		{`<enclosure length="12216320" type="audio/mpeg" />`,
			xmlutils.NewError(MissingAttribute, ""),
			NewTestEnclosure("", "12216320", "audio/mpeg"),
		},
		{`<enclosure url="http://www.scripting.com/mp3s/weatherReportSuite.mp3" type="audio/mpeg" />`,
			xmlutils.NewError(MissingAttribute, ""),
			NewTestEnclosure("http://www.scripting.com/mp3s/weatherReportSuite.mp3", "", "audio/mpeg"),
		},
		{`<enclosure url="http://www.scripting.com/mp3s/weatherReportSuite.mp3" length="12216320" />`,
			xmlutils.NewError(MissingAttribute, ""),
			NewTestEnclosure("http://www.scripting.com/mp3s/weatherReportSuite.mp3", "12216320", ""),
		},
	}

	nbErrors := 0
	len := len(testdata)
	for _, testenclosure := range testdata {
		testcase := _TestEnclosureToTestVisitor(testenclosure)

		if err := testcase.CheckTestCase(); err != nil {
			t.Errorf("FAIL\n%s\nXML:\n %s\n", err, testcase.XML)
			nbErrors++
		}
	}

	t.Logf("PASS RATIO = %v/%v\n", len-nbErrors, len)
}
