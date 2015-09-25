package atom

import (
	"fmt"
	"testing"

	xmlutils "github.com/jloup/xml/utils"
)

func NewTestId(iri string) *Id {
	i := NewId()
	i.Content.Value = iri

	return i
}

func IdWithBaseLang(i *Id, lang, base string) *Id {
	i.Lang.Value = lang
	i.Base.Value = base

	return i
}

type testId struct {
	XML           string
	ExpectedError xmlutils.ParserError
	ExpectedId    *Id
}

func testIdValidator(actual xmlutils.Visitor, expected xmlutils.Visitor) error {
	i1 := actual.(*Id)
	i2 := expected.(*Id)

	if i1.Content.Value != i2.Content.Value {
		return fmt.Errorf("Content is invalid '%s' (expected) vs '%s'", i2.Content.Value, i1.Content.Value)
	}

	if err := ValidateBaseLang("Id", i1.Base.Value, i1.Lang.Value, i2.Base.Value, i2.Lang.Value); err != nil {
		return err
	}

	return nil
}

func testIdConstructor() xmlutils.Visitor {
	return NewId()
}

func testIdToTestVisitor(t testId) xmlutils.TestVisitor {
	testVisitor := xmlutils.TestVisitor{
		XML:                t.XML,
		ExpectedError:      nil,
		ExpectedVisitor:    t.ExpectedId,
		VisitorConstructor: testIdConstructor,
		Validator:          testIdValidator,
	}

	if t.ExpectedError != nil {
		testVisitor.ExpectedError = t.ExpectedError
	}

	return testVisitor
}

func TestIdBasic(t *testing.T) {

	var testdata = []testId{
		{`<id xml:lang="en-us" xml:base="http://yo.com">https://www.yo.com</id>`,
			nil,
			IdWithBaseLang(NewTestId("https://www.yo.com"), "en-us", "http://yo.com"),
		},
		{`<id>https://www.yo.com</id>`,
			nil,
			NewTestId("https://www.yo.com"),
		},
		{`<id>https://www.%yo.com</id>`,
			xmlutils.NewError(IriNotAbsolute, ""),
			NewTestId("https://www.%yo.com"),
		},
		{`<id>
      https://www.%yo.com
     </id>`,
			xmlutils.NewError(IriNotAbsolute, ""),
			NewTestId("https://www.%yo.com"),
		},
	}
	nbErrors := 0
	len := len(testdata)
	for _, testid := range testdata {
		testcase := testIdToTestVisitor(testid)

		if err := testcase.CheckTestCase(); err != nil {
			t.Errorf("FAIL\n%s\nXML:\n %s\n", err, testcase.XML)
			nbErrors++
		}
	}

	t.Logf("PASS RATIO = %v/%v\n", len-nbErrors, len)
}
