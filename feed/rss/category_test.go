package rss

import (
	"github.com/JLoup/xml/utils"
	"fmt"
	"testing"
)

func NewTestCategory(domain, content string) *Category {
	c := NewCategory()

	c.Domain.Value = domain
	c.Content.Value = content

	return c
}

type testCategory struct {
	XML              string
	ExpectedError    utils.ParserError
	ExpectedCategory *Category
}

func testCategoryValidator(actual utils.Visitor, expected utils.Visitor) error {
	c1 := actual.(*Category)
	c2 := expected.(*Category)

	if c1.Domain.Value != c2.Domain.Value {
		return fmt.Errorf("Domain is invalid '%s' (expected) vs '%s'", c2.Domain.Value, c1.Domain.Value)
	}

	if c1.Content.Value != c2.Content.Value {
		return fmt.Errorf("Content is invalid '%s' (expected) vs '%s'", c2.Content.Value, c1.Content.Value)
	}

	return nil
}

func testCategoryConstructor() utils.Visitor {
	return NewCategory()
}

func _TestCategoryToTestVisitor(t testCategory) utils.TestVisitor {
	testVisitor := utils.TestVisitor{
		XML:                t.XML,
		ExpectedError:      nil,
		ExpectedVisitor:    t.ExpectedCategory,
		VisitorConstructor: testCategoryConstructor,
		Validator:          testCategoryValidator,
	}

	if t.ExpectedError != nil {
		testVisitor.ExpectedError = t.ExpectedError
	}

	return testVisitor
}
func TestCategoryBasic(t *testing.T) {

	var testdata = []testCategory{
		{`<category domain="http://www.fool.com/cusips">MSFT</category>`,
			nil,
			NewTestCategory("http://www.fool.com/cusips", "MSFT"),
		},
		{`<category>MSFT</category>`,
			nil,
			NewTestCategory("", "MSFT"),
		},
	}

	nbErrors := 0
	len := len(testdata)
	for _, testcategory := range testdata {
		testcase := _TestCategoryToTestVisitor(testcategory)

		if err := testcase.CheckTestCase(); err != nil {
			t.Errorf("FAIL\n%s\nXML:\n %s\n", err, testcase.XML)
			nbErrors++
		}
	}

	t.Logf("PASS RATIO = %v/%v\n", len-nbErrors, len)
}
