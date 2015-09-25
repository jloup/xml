package atom

import (
	"fmt"
	"testing"

	xmlutils "github.com/jloup/xml/utils"
)

func NewTestLogo(iri string) *Logo {
	l := NewLogo()
	l.Iri.Value = iri

	return l
}

func LogoWithBaseLang(l *Logo, lang, base string) *Logo {
	l.Lang.Value = lang
	l.Base.Value = base

	return l
}

type testLogo struct {
	XML           string
	ExpectedError xmlutils.ParserError
	ExpectedLogo  *Logo
}

func testLogoValidator(actual xmlutils.Visitor, expected xmlutils.Visitor) error {
	l1 := actual.(*Logo)
	l2 := expected.(*Logo)

	if l1.Iri.Value != l2.Iri.Value {
		return fmt.Errorf("Logo Iri is invalid '%s' (expected) vs '%s'", l2.Iri.Value, l1.Iri.Value)
	}

	if err := ValidateBaseLang("Logo", l1.Base.Value, l1.Lang.Value, l2.Base.Value, l2.Lang.Value); err != nil {
		return err
	}

	return nil
}

func testLogoConstructor() xmlutils.Visitor {
	return NewLogo()
}

func testLogoToTestVisitor(t testLogo) xmlutils.TestVisitor {
	testVisitor := xmlutils.TestVisitor{
		XML:                t.XML,
		ExpectedError:      nil,
		ExpectedVisitor:    t.ExpectedLogo,
		VisitorConstructor: testLogoConstructor,
		Validator:          testLogoValidator,
	}

	if t.ExpectedError != nil {
		testVisitor.ExpectedError = t.ExpectedError
	}

	return testVisitor
}

func TestLogoBasic(t *testing.T) {

	var testdata = []testLogo{
		{`<logo xml:lang="en-us" xml:base="http://yo.com">https://www.yo.com</logo>`,
			nil,
			LogoWithBaseLang(NewTestLogo("https://www.yo.com"), "en-us", "http://yo.com"),
		},
		{`<logo>https://www.yo.com</logo>`,
			nil,
			NewTestLogo("https://www.yo.com"),
		},
		{`<logo>https://www.%yo.com</logo>`,
			xmlutils.NewError(IriNotValid, ""),
			NewTestLogo("https://www.%yo.com"),
		},
	}
	nbErrors := 0
	len := len(testdata)
	for _, testlogo := range testdata {
		testcase := testLogoToTestVisitor(testlogo)

		if err := testcase.CheckTestCase(); err != nil {
			t.Errorf("FAIL\n%s\nXML:\n %s\n", err, testcase.XML)
			nbErrors++
		}
	}

	t.Logf("PASS RATIO = %v/%v\n", len-nbErrors, len)
}
