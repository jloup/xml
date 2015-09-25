package rss

import (
	xmlutils "github.com/jloup/xml/utils"
	"fmt"
	"testing"
)

func NewTestGuid(isPermalink, content string) *Guid {
	g := NewGuid()

	g.IsPermalink.Value = isPermalink
	g.Content.Value = content

	return g
}

type testGuid struct {
	XML           string
	ExpectedError xmlutils.ParserError
	ExpectedGuid  *Guid
}

func testGuidValidator(actual xmlutils.Visitor, expected xmlutils.Visitor) error {
	g1 := actual.(*Guid)
	g2 := expected.(*Guid)

	if g1.IsPermalink.Value != g2.IsPermalink.Value {
		return fmt.Errorf("IsPermalink is invalid '%s' (expected) vs '%s'", g2.IsPermalink.Value, g1.IsPermalink.Value)
	}

	if g1.Content.Value != g2.Content.Value {
		return fmt.Errorf("Content is invalid '%s' (expected) vs '%s'", g2.Content.Value, g1.Content.Value)
	}

	return nil
}

func testGuidConstructor() xmlutils.Visitor {
	return NewGuid()
}

func _TestGuidToTestVisitor(t testGuid) xmlutils.TestVisitor {
	testVisitor := xmlutils.TestVisitor{
		XML:                t.XML,
		ExpectedError:      nil,
		ExpectedVisitor:    t.ExpectedGuid,
		VisitorConstructor: testGuidConstructor,
		Validator:          testGuidValidator,
	}

	if t.ExpectedError != nil {
		testVisitor.ExpectedError = t.ExpectedError
	}

	return testVisitor
}
func TestGuidBasic(t *testing.T) {

	var testdata = []testGuid{
		{`<guid isPermaLink="false">http://inessential.com/2002/09/01.php#a2</guid>`,
			nil,
			NewTestGuid("false", "http://inessential.com/2002/09/01.php#a2"),
		},
		{`<guid>http://inessential.com/2002/09/01.php#a2</guid>`,
			nil,
			NewTestGuid("true", "http://inessential.com/2002/09/01.php#a2"),
		},
	}

	nbErrors := 0
	len := len(testdata)
	for _, testguid := range testdata {
		testcase := _TestGuidToTestVisitor(testguid)

		if err := testcase.CheckTestCase(); err != nil {
			t.Errorf("FAIL\n%s\nXML:\n %s\n", err, testcase.XML)
			nbErrors++
		}
	}

	t.Logf("PASS RATIO = %v/%v\n", len-nbErrors, len)
}
