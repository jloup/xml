package atom

import (
	"fmt"
	"testing"

	"github.com/JLoup/xml/utils"
)

func NewTestTextConstruct(xhtml, plaintext string) *TextConstruct {
	t := NewTextConstruct()

	t.PlainText.Content = plaintext

	t.XHTML.Content.WriteString(xhtml)

	return t
}

func TextConstructWithBaseLang(t *TextConstruct, lang, base string) *TextConstruct {
	t.Lang.Value = lang
	t.Base.Value = base

	return t
}

type testTextConstruct struct {
	XML                   string
	ExpectedError         utils.ParserError
	ExpectedTextConstruct *TextConstruct
}

func testTextConstructValidator(actual utils.Visitor, expected utils.Visitor) error {
	t1 := actual.(*TextConstruct)
	t2 := expected.(*TextConstruct)

	if t1.PlainText.String() != t2.PlainText.String() {
		return fmt.Errorf("%s Text content is invalid '%s' (expected) vs '%s'", t1.name, t2.PlainText.String(), t1.PlainText.String())
	}

	if t1.XHTML.String() != t2.XHTML.String() {
		return fmt.Errorf("%s XHTML content is invalid '%s' (expected) vs '%s'", t1.name, t2.XHTML.String(), t1.XHTML.String())
	}

	if err := ValidateBaseLang(t1.name, t1.Base.Value, t1.Lang.Value, t2.Base.Value, t2.Lang.Value); err != nil {
		return err
	}

	return nil
}

func testTextConstructConstructor() utils.Visitor {
	return NewTextConstruct()
}

func _TestTextConstructToTestVisitor(t testTextConstruct) utils.TestVisitor {
	testVisitor := utils.TestVisitor{
		XML:                t.XML,
		ExpectedError:      nil,
		ExpectedVisitor:    t.ExpectedTextConstruct,
		VisitorConstructor: testTextConstructConstructor,
		Validator:          testTextConstructValidator,
	}

	if t.ExpectedError != nil {
		testVisitor.ExpectedError = t.ExpectedError
	}

	return testVisitor
}

func TestTextConstructBasic(t *testing.T) {

	var testdata = []testTextConstruct{
		{`
    <summary type="xhtml" xml:lang="en-us" xml:base="http://yo.com"><div xmlns="http://www.w3.org/1999/xhtml">This ' is <b>XHTML</b> content</div></summary>`,
			nil,
			TextConstructWithBaseLang(NewTestTextConstruct("<div>This ' is <b>XHTML</b> content</div>", ""), "en-us", "http://yo.com"),
		},
		{`
    <summary type="xhtml"><div xmlns="http://www.w3.org/1999/xhtml">This is is <b>XHTML</b> content</div><div lambda="3" xmlns="http://www.w3.org/1999/xhtml">this is not usual child</div></summary>`,
			utils.NewError(NotUniqueChild, ""),
			NewTestTextConstruct("<div>This is is <b>XHTML</b> content</div><div lambda=\"3\">this is not usual child</div>", ""),
		},
		{`
    <summary type="xhtml"><div xmlns="http://www.w3.org/1999/xhtml">This ' is <b>XHTML</b> content</div></summary>`,
			nil,
			NewTestTextConstruct("<div>This ' is <b>XHTML</b> content</div>", ""),
		},
		{`
    <summary type="xhtml"><div>This is is <b>XHTML</b> content</div></summary>`,
			utils.NewError(XHTMLElementNotNamespaced, ""),
			NewTestTextConstruct("<div>This is is <b>XHTML</b> content</div>", ""),
		},
		{`
    <summary>This is plain text</summary>`,
			nil,
			NewTestTextConstruct("", "This is plain text"),
		},
		{`
    <summary type="text">This is plain text</summary>`,
			nil,
			NewTestTextConstruct("", "This is plain text"),
		},
		{`
    <summary type="text"><div><a>CHILD</a></div><div>CHILD2</div></summary>`,
			utils.NewError(LeafElementHasChild, ""),
			NewTestTextConstruct("", "CHILD2"),
		},
		{`
    <summary type="html">This is plain text</summary>`,
			nil,
			NewTestTextConstruct("", "This is plain text"),
		},
		{`
    <summary type="xhtml" xmlns:h="http://www.w3.org/1999/xhtml"><h:div>This is is <h:b>XHTML</h:b> content</h:div></summary>`,
			nil,
			NewTestTextConstruct("<div>This is is <b>XHTML</b> content</div>", ""),
		},
		{`
    <summary type="xhtml" xmlns:h="http://www.w3.org/1999/xhtml"><h:p>This is is <h:b>XHTML</h:b> content</h:p></summary>`,
			utils.NewError(XHTMLRootNodeNotDiv, ""),
			NewTestTextConstruct("<p>This is is <b>XHTML</b> content</p>", ""),
		},
	}

	nbErrors := 0
	len := len(testdata)
	for _, testtextconstruct := range testdata {
		testcase := _TestTextConstructToTestVisitor(testtextconstruct)

		if err := testcase.CheckTestCase(); err != nil {
			t.Errorf("FAIL\n%s\nXML:\n %s\n", err, testcase.XML)
			nbErrors++
		}
	}

	t.Logf("PASS RATIO = %v/%v\n", len-nbErrors, len)
}
