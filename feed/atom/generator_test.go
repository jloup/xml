package atom

import (
	"fmt"
	"testing"

	xmlutils "github.com/jloup/xml/utils"
)

func NewTestGenerator(version, uri, content string) *Generator {
	g := NewGenerator()
	g.Version.Value = version
	g.Uri.Value = uri
	g.Content = content

	return g
}

func GeneratorWithBaseLang(g *Generator, lang, base string) *Generator {
	g.Lang.Value = lang
	g.Base.Value = base

	return g
}

type testGenerator struct {
	XML               string
	ExpectedError     xmlutils.ParserError
	ExpectedGenerator *Generator
}

func testGeneratorValidator(actual xmlutils.Visitor, expected xmlutils.Visitor) error {
	g1 := actual.(*Generator)
	g2 := expected.(*Generator)

	if g1.Uri.Value != g2.Uri.Value {
		return fmt.Errorf("Uri is invalid '%s' (expected) vs '%s'", g2.Uri.Value, g1.Uri.Value)
	}

	if g1.Version.Value != g2.Version.Value {
		return fmt.Errorf("Version is invalid '%s' (expected) vs '%s'", g2.Version.Value, g1.Version.Value)
	}

	if g1.Content != g2.Content {
		return fmt.Errorf("Content is invalid '%s' (expected) vs '%s'", g2.Content, g1.Content)
	}

	if err := ValidateBaseLang("Generator", g1.Base.Value, g1.Lang.Value, g2.Base.Value, g2.Lang.Value); err != nil {
		return err
	}

	return nil
}

func testGeneratorConstructor() xmlutils.Visitor {
	return NewGenerator()
}

func testGeneratorToTestVisitor(t testGenerator) xmlutils.TestVisitor {
	testVisitor := xmlutils.TestVisitor{
		XML:                t.XML,
		ExpectedError:      nil,
		ExpectedVisitor:    t.ExpectedGenerator,
		VisitorConstructor: testGeneratorConstructor,
		Validator:          testGeneratorValidator,
	}

	if t.ExpectedError != nil {
		testVisitor.ExpectedError = t.ExpectedError
	}

	return testVisitor
}

func TestGeneratorBasic(t *testing.T) {

	var testdata = []testGenerator{
		{`<generator uri="http://there.com" version="4.0" xml:lang="en-us" xml:base="http://yo.com">my generator</generator>`,
			nil,
			GeneratorWithBaseLang(NewTestGenerator("4.0", "http://there.com", "my generator"), "en-us", "http://yo.com"),
		},
		{`<generator uri="http://there.com" version="4.0">my generator</generator>`,
			nil,
			NewTestGenerator("4.0", "http://there.com", "my generator"),
		},
		{`<generator uri="http://there.com" version="4.0"><name>CHILD</name></generator>`,
			xmlutils.NewError(LeafElementHasChild, ""),
			NewTestGenerator("4.0", "http://there.com", "CHILD"),
		},
		{`<generator uri="http://%there.com" version="4.0">my generator</generator>`,
			xmlutils.NewError(IriNotValid, ""),
			NewTestGenerator("4.0", "http://%there.com", "my generator"),
		},
	}
	nbErrors := 0
	len := len(testdata)
	for _, testgenerator := range testdata {
		testcase := testGeneratorToTestVisitor(testgenerator)

		if err := testcase.CheckTestCase(); err != nil {
			t.Errorf("FAIL\n%s\nXML:\n %s\n", err, testcase.XML)
			nbErrors++
		}
	}

	t.Logf("PASS RATIO = %v/%v\n", len-nbErrors, len)
}
