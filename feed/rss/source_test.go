package rss

import (
	"github.com/JLoup/xml/utils"
	"fmt"
	"testing"
)

func NewTestSource(url, content string) *Source {
	s := NewSource()

	s.Url.Value = url
	s.Content.Value = content

	return s
}

type testSource struct {
	XML            string
	ExpectedError  utils.ParserError
	ExpectedSource *Source
}

func testSourceValidator(actual utils.Visitor, expected utils.Visitor) error {
	s1 := actual.(*Source)
	s2 := expected.(*Source)

	if s1.Url.Value != s2.Url.Value {
		return fmt.Errorf("Url is invalid '%s' (expected) vs '%s'", s2.Url.Value, s1.Url.Value)
	}

	if s1.Content.Value != s2.Content.Value {
		return fmt.Errorf("Content is invalid '%s' (expected) vs '%s'", s2.Content.Value, s1.Content.Value)
	}

	return nil
}

func testSourceConstructor() utils.Visitor {
	return NewSource()
}

func _TestSourceToTestVisitor(t testSource) utils.TestVisitor {
	testVisitor := utils.TestVisitor{
		XML:                t.XML,
		ExpectedError:      nil,
		ExpectedVisitor:    t.ExpectedSource,
		VisitorConstructor: testSourceConstructor,
		Validator:          testSourceValidator,
	}

	if t.ExpectedError != nil {
		testVisitor.ExpectedError = t.ExpectedError
	}

	return testVisitor
}
func TestSourceBasic(t *testing.T) {

	var testdata = []testSource{
		{`<source>Tomalak's Realm</source>`,
			utils.NewError(MissingAttribute, ""),
			NewTestSource("", "Tomalak's Realm"),
		},
		{`<source url="http://www.tomaltak.org/links2.xml">Tomalak's Realm</source>`,
			nil,
			NewTestSource("http://www.tomaltak.org/links2.xml", "Tomalak's Realm"),
		},
	}

	nbErrors := 0
	len := len(testdata)
	for _, testsource := range testdata {
		testcase := _TestSourceToTestVisitor(testsource)

		if err := testcase.CheckTestCase(); err != nil {
			t.Errorf("FAIL\n%s\nXML:\n %s\n", err, testcase.XML)
			nbErrors++
		}
	}

	t.Logf("PASS RATIO = %v/%v\n", len-nbErrors, len)
}
