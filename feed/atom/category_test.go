package atom

import (
	"fmt"
	"testing"

	"github.com/JLoup/xml/helper"
)

func ValidateBaseLang(name, base, lang, eBase, eLang string) error {
	if lang != eLang {
		return fmt.Errorf("%s lang is invalid '%s' (expected) vs '%s'", name, eLang, lang)
	}

	if base != eBase {
		return fmt.Errorf("%s base is invalid '%s' (expected) vs '%s'", name, eBase, base)
	}

	return nil
}

func NewTestCategory(scheme, term, label string) *Category {
	c := NewCategory()
	c.Scheme.Value = scheme
	c.Label.Value = label
	c.Term.Value = term

	return c
}

func CategoryWithBaseLang(c *Category, lang, base string) *Category {
	c.Lang.Value = lang
	c.Base.Value = base

	return c
}

type testCategory struct {
	XML              string
	ExpectedError    helper.ParserError
	ExpectedCategory *Category
}

func testCategoryValidator(actual helper.Visitor, expected helper.Visitor) error {
	c1 := actual.(*Category)
	c2 := expected.(*Category)

	if c1.Term.Value != c2.Term.Value {
		return fmt.Errorf("Category term is invalid '%s' (expected) vs '%s'", c2.Term.Value, c1.Term.Value)
	}

	if c1.Label.Value != c2.Label.Value {
		return fmt.Errorf("Category label is invalid '%s' (expected) vs '%s'", c2.Label.Value, c1.Label.Value)
	}

	if c1.Scheme.Value != c2.Scheme.Value {
		return fmt.Errorf("Category scheme is invalid '%s' (expected) vs '%s'", c2.Scheme.Value, c1.Scheme.Value)
	}

	if err := ValidateBaseLang("Category", c1.Base.Value, c1.Lang.Value, c2.Base.Value, c2.Lang.Value); err != nil {
		return err
	}

	return nil
}

func testCategoryConstructor() helper.Visitor {
	return NewCategory()
}

func testCategoryToTestVisitor(t testCategory) helper.TestVisitor {
	testVisitor := helper.TestVisitor{
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
		{`<category xml:lang="en-us" xml:base="http://yo.com" scheme='https://www.yo.com' term='music' />`,
			nil,
			CategoryWithBaseLang(NewTestCategory("https://www.yo.com", "music", ""), "en-us", "http://yo.com"),
		},

		{`<category scheme='https://www.yo.com' term='music' />`,
			nil,
			NewTestCategory("https://www.yo.com", "music", ""),
		},
		{`<category scheme='https://www.tbray.org/ongoing/What/' />`,
			helper.NewError(MissingAttribute, ""),
			NewTestCategory("https://www.tbray.org/ongoing/What/", "", ""),
		},
		{`<category term="music" label='house' />`,
			nil,
			NewTestCategory("", "music", "house"),
		},
		{`<category term="music" label='house' scheme='http://www.%scheme.com'/>`,
			helper.NewError(IriNotValid, ""),
			NewTestCategory("http://www.%scheme.com", "music", "house"),
		},
	}
	nbErrors := 0
	len := len(testdata)
	for _, testcategory := range testdata {
		testcase := testCategoryToTestVisitor(testcategory)

		if err := testcase.CheckTestCase(); err != nil {
			t.Errorf("FAIL\n%s\nXML:\n %s\n", err, testcase.XML)
			nbErrors++
		}
	}

	t.Logf("PASS RATIO = %v/%v\n", len-nbErrors, len)
}
