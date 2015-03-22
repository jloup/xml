package atom

import (
	"fmt"
	"testing"

	"github.com/JLoup/xml/helper"
)

func NewTestLink(href, rel, typ, hrelang, title, length string) *Link {
	l := NewLink()
	l.Href.Value = href
	l.Rel.Value = rel
	l.Type.Value = typ
	l.HrefLang.Value = hrelang
	l.Title.Value = title
	l.Length.Value = length

	return l
}

func LinkWithBaseLang(l *Link, lang, base string) *Link {
	l.Lang.Value = lang
	l.Base.Value = base

	return l
}

type testLink struct {
	XML           string
	ExpectedError helper.ParserError
	ExpectedLink  *Link
}

func testLinkValidator(actual helper.Visitor, expected helper.Visitor) error {
	l1 := actual.(*Link)
	l2 := expected.(*Link)

	if l1.Href.Value != l2.Href.Value {
		return fmt.Errorf("Href is invalid '%s' (expected) vs '%s'", l2.Href.Value, l1.Href.Value)
	}

	if l1.Rel.Value != l2.Rel.Value {
		return fmt.Errorf("Rel is invalid '%s' (expected) vs '%s'", l2.Rel.Value, l1.Rel.Value)
	}

	if l1.Type.Value != l2.Type.Value {
		return fmt.Errorf("Type is invalid '%s' (expected) vs '%s'", l2.Type.Value, l1.Type.Value)
	}

	if l1.HrefLang.Value != l2.HrefLang.Value {
		return fmt.Errorf("HrefLang is invalid '%s' (expected) vs '%s'", l2.HrefLang.Value, l1.HrefLang.Value)
	}

	if l1.Title.Value != l2.Title.Value {
		return fmt.Errorf("Title is invalid '%s' (expected) vs '%s'", l2.Title.Value, l1.Title.Value)
	}

	if l1.Length.Value != l2.Length.Value {
		return fmt.Errorf("Length is invalid '%s' (expected) vs '%s'", l2.Length.Value, l1.Length.Value)
	}

	if err := ValidateBaseLang("Link", l1.Base.Value, l1.Lang.Value, l2.Base.Value, l2.Lang.Value); err != nil {
		return err
	}

	return nil
}

func testLinkConstructor() helper.Visitor {
	return NewLink()
}

func testLinkToTestVisitor(t testLink) helper.TestVisitor {
	testVisitor := helper.TestVisitor{
		XML:                t.XML,
		ExpectedError:      nil,
		ExpectedVisitor:    t.ExpectedLink,
		VisitorConstructor: testLinkConstructor,
		Validator:          testLinkValidator,
	}

	if t.ExpectedError != nil {
		testVisitor.ExpectedError = t.ExpectedError
	}

	return testVisitor
}

func TestLinkBasic(t *testing.T) {

	var testdata = []testLink{
		{`<link xml:lang="en-us" xml:base="http://yo.com" href="http://www.go.com" rel="alternate" type="text/html"/>`,
			nil,
			LinkWithBaseLang(NewTestLink("http://www.go.com", "alternate", "text/html", "", "", ""), "en-us", "http://yo.com"),
		},
		{`<link rel="alternate" type="text/html"/>`,
			helper.NewError(MissingAttribute, ""),
			NewTestLink("", "alternate", "text/html", "", "", ""),
		},
		{`<link rel="alternate" type="text/html" href="http://%%%www.go.com"/>`,
			helper.NewError(IriNotValid, ""),
			NewTestLink("http://%%%www.go.com", "alternate", "text/html", "", "", ""),
		},
		{`<link rel="alternate" type="text/html" href="http://www.go.com" length="-56"/>`,
			helper.NewError(NotPositiveNumber, ""),
			NewTestLink("http://www.go.com", "alternate", "text/html", "", "", "-56"),
		},
		{`<link rel="alternate" type="text/html" href="http://www.go.com" length="ab"/>`,
			helper.NewError(NotPositiveNumber, ""),
			NewTestLink("http://www.go.com", "alternate", "text/html", "", "", "ab"),
		},
	}
	nbErrors := 0
	len := len(testdata)
	for _, testlink := range testdata {
		testcase := testLinkToTestVisitor(testlink)

		if err := testcase.CheckTestCase(); err != nil {
			t.Errorf("FAIL\n%s\nXML:\n %s\n", err, testcase.XML)
			nbErrors++
		}
	}

	t.Logf("PASS RATIO = %v/%v\n", len-nbErrors, len)
}
