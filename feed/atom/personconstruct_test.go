package atom

import (
	"fmt"
	"testing"

	"github.com/JLoup/xml/utils"
)

func NewTestPerson(name, uri, email string) *Person {
	p := NewPerson()

	p.Name.Content.Value = name
	p.Uri.Content.Value = uri
	p.Email.Content.Value = email

	return p
}

func PersonWithBaseLang(p *Person, lang, base string) *Person {
	p.Lang.Value = lang
	p.Base.Value = base

	return p
}

type testPerson struct {
	XML            string
	ExpectedError  utils.ParserError
	ExpectedPerson *Person
}

func testPersonValidator(actual utils.Visitor, expected utils.Visitor) error {
	p1 := actual.(*Person)
	p2 := expected.(*Person)

	if p1.Name.Content.Value != p2.Name.Content.Value {
		return fmt.Errorf("%s name is invalid '%s' (expected) vs '%s'", p1.name, p2.Name.Content.Value, p1.Name.Content.Value)
	}

	if p1.Email.Content.Value != p2.Email.Content.Value {
		return fmt.Errorf("%s email is invalid '%s' (expected) vs '%s'", p1.name, p2.Email.Content.Value, p1.Email.Content.Value)
	}

	if p1.Uri.Content.Value != p2.Uri.Content.Value {
		return fmt.Errorf("%s uri is invalid '%s' (expected) vs '%s'", p1.name, p2.Uri.Content.Value, p1.Uri.Content.Value)
	}

	if err := ValidateBaseLang(p1.name, p1.Base.Value, p1.Lang.Value, p2.Base.Value, p2.Lang.Value); err != nil {
		return err
	}

	return nil
}

func testPersonConstructor() utils.Visitor {
	return NewPerson()
}

func _TestPersonToTestVisitor(t testPerson) utils.TestVisitor {
	testVisitor := utils.TestVisitor{
		XML:                t.XML,
		ExpectedError:      nil,
		ExpectedVisitor:    t.ExpectedPerson,
		VisitorConstructor: testPersonConstructor,
		Validator:          testPersonValidator,
	}

	if t.ExpectedError != nil {
		testVisitor.ExpectedError = t.ExpectedError
	}

	return testVisitor
}

var testdata = []testPerson{
	{`
	<author xmlns='http://www.w3.org/2005/Atom' xml:lang="en-us" xml:base="http://yo.com">
     <name>Tim Bray</name>
     <email>tim.bray@machin.com</email>
    </author>`,
		nil,
		PersonWithBaseLang(NewTestPerson("Tim Bray", "", "tim.bray@machin.com"), "en-us", "http://yo.com"),
	},
	{`
    <author xmlns='http://www.w3.org/2005/Atom'>
     <name>Tim Bray</name>
     <email>tim.bray@machin.com</email>
    </author>`,
		nil,
		NewTestPerson("Tim Bray", "", "tim.bray@machin.com"),
	},
	{`<author xmlns='http://www.w3.org/2005/Atom'>
     <uri>http://yo.com</uri>
     <email>tim.bray@machin.com</email>
    </author>`,
		utils.NewError(MissingAttribute, ""),
		NewTestPerson("", "http://yo.com", "tim.bray@machin.com"),
	},
	{`
    <author xmlns='http://www.w3.org/2005/Atom'>
     <name>Tim Bray</name>
     <name>Tim Bray</name>
     <email>tim.bray@machin.com</email>
    </author>`,
		utils.NewError(AttributeDuplicated, ""),
		NewTestPerson("Tim Bray", "", "tim.bray@machin.com"),
	},
	{`
    <author xmlns='http://www.w3.org/2005/Atom'>
     <name>Tim Bray</name>
     <email><a>tim.bray@machin.com</a></email>
    </author>`,
		utils.NewError(LeafElementHasChild, ""),
		NewTestPerson("Tim Bray", "", "tim.bray@machin.com"),
	},
	{`
    <author xmlns='http://www.w3.org/2005/Atom'>
     <name>Tim Bray</name>
     <email>tim.bray@machin.com</email>
     <uri>http://www.yo.com</uri>
    </author>`,
		nil,
		NewTestPerson("Tim Bray", "http://www.yo.com", "tim.bray@machin.com"),
	},
	{`
    <author xmlns='http://www.w3.org/2005/Atom'>
     <name>Tim Bray</name>
     <email>tim.bray@machin.com</email>
     <uri>http://www%%.yo.com</uri>
    </author>`,
		utils.NewError(IriNotValid, ""),
		NewTestPerson("Tim Bray", "http://www%%.yo.com", "tim.bray@machin.com"),
	},
	{`
    <author xmlns='http://www.w3.org/2005/Atom'>
    </author>`,
		utils.NewError(MissingAttribute, ""),
		NewTestPerson("", "", ""),
	},
}

func TestPersonBasic(t *testing.T) {
	nbErrors := 0
	len := len(testdata)
	for _, testperson := range testdata {
		testcase := _TestPersonToTestVisitor(testperson)

		if err := testcase.CheckTestCase(); err != nil {
			t.Errorf("FAIL\n%s\nXML:\n %s\n", err, testcase.XML)
			nbErrors++
		}
	}

	t.Logf("PASS RATIO = %v/%v\n", len-nbErrors, len)
}
