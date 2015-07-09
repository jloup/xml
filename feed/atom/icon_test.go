package atom

import (
	"fmt"
	"testing"

	"github.com/jloup/xml/utils"
)

func NewTestIcon(iri string) *Icon {
	i := NewIcon()
	i.Iri.Value = iri

	return i
}

func IconWithBaseLang(i *Icon, lang, base string) *Icon {
	i.Lang.Value = lang
	i.Base.Value = base

	return i
}

type testIcon struct {
	XML           string
	ExpectedError utils.ParserError
	ExpectedIcon  *Icon
}

func testIconValidator(actual utils.Visitor, expected utils.Visitor) error {
	i1 := actual.(*Icon)
	i2 := expected.(*Icon)

	if i1.Iri.Value != i2.Iri.Value {
		return fmt.Errorf("Content is invalid '%s' (expected) vs '%s'", i2.Iri, i1.Iri)
	}

	if err := ValidateBaseLang("Icon", i1.Base.Value, i1.Lang.Value, i2.Base.Value, i2.Lang.Value); err != nil {
		return err
	}

	return nil
}
func testIconConstructor() utils.Visitor {
	return NewIcon()
}

func testIconToTestVisitor(t testIcon) utils.TestVisitor {
	testVisitor := utils.TestVisitor{
		XML:                t.XML,
		ExpectedError:      nil,
		ExpectedVisitor:    t.ExpectedIcon,
		VisitorConstructor: testIconConstructor,
		Validator:          testIconValidator,
	}

	if t.ExpectedError != nil {
		testVisitor.ExpectedError = t.ExpectedError
	}

	return testVisitor
}

func TestIconBasic(t *testing.T) {

	var testdata = []testIcon{
		{`<icon xml:lang="en-us" xml:base="http://yo.com">https://www.yo.com</icon>`,
			nil,
			IconWithBaseLang(NewTestIcon("https://www.yo.com"), "en-us", "http://yo.com"),
		},
		{`<icon>https://www.yo.com</icon>`,
			nil,
			NewTestIcon("https://www.yo.com"),
		},
		{`<icon>https://www.%yo.com</icon>`,
			utils.NewError(IriNotValid, ""),
			NewTestIcon("https://www.%yo.com"),
		},
	}
	nbErrors := 0
	len := len(testdata)
	for _, testicon := range testdata {
		testcase := testIconToTestVisitor(testicon)

		if err := testcase.CheckTestCase(); err != nil {
			t.Errorf("FAIL\n%s\nXML:\n %s\n", err, testcase.XML)
			nbErrors++
		}
	}

	t.Logf("PASS RATIO = %v/%v\n", len-nbErrors, len)
}
