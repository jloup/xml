package thr

import (
	"fmt"
	"testing"

	"github.com/jloup/xml/feed/atom"
	"github.com/jloup/xml/feed/extension"
	xmlutils "github.com/jloup/xml/utils"
)

func NewTestThrEntry(links []*atom.Link, total, href, ref, typ, source string) *atom.Entry {
	e := atom.NewEntry()

	if href != "" || ref != "" || typ != "" || source != "" {
		inreplyto := newInReplyTo()
		inreplyto.Href.Value = href
		inreplyto.Ref.Value = ref
		inreplyto.Type.Value = typ
		inreplyto.Source.Value = source
		e.Extension.Store.Add(_inreplyto, inreplyto)
	}

	if total != "" {
		t := atom.NewBasicElement(nil)
		t.Content.Value = total
		e.Extension.Store.Add(_total, t)
	}

	e.Links = links

	return e
}

func NewTestThrLink(updated, count string) *atom.Link {
	i := atom.NewLink()

	if updated != "" {
		u := newUpdatedAttr()
		u.Set(updated)
		i.Extension.Store.Add(_updated, u)
	}

	if count != "" {
		c := newCountAttr()
		c.Set(count)
		i.Extension.Store.Add(_count, c)
	}

	return i
}

type testThrEntry struct {
	XML              string
	ExpectedError    xmlutils.ParserError
	ExpectedThrEntry *atom.Entry
}

func testInReplyToValidator(actual *InReplyTo, expected *InReplyTo) error {
	if actual.Href.Value != expected.Href.Value {
		return fmt.Errorf("InReplyTo href is invalid '%s' (expected) vs '%s'", expected.Href.Value, actual.Href.Value)
	}

	if actual.Ref.Value != expected.Ref.Value {
		return fmt.Errorf("InReplyTo ref is invalid '%s' (expected) vs '%s'", expected.Ref.Value, actual.Ref.Value)
	}

	if actual.Type.Value != expected.Type.Value {
		return fmt.Errorf("InReplyTo type is invalid '%s' (expected) vs '%s'", expected.Type.Value, actual.Type.Value)
	}

	if actual.Source.Value != expected.Source.Value {
		return fmt.Errorf("InReplyTo source is invalid '%s' (expected) vs '%s'", expected.Source.Value, actual.Source.Value)
	}

	return nil
}

func testThrEntryValidator(actual xmlutils.Visitor, expected xmlutils.Visitor) error {
	e1 := actual.(*atom.Entry)
	e2 := expected.(*atom.Entry)

	irt1, ok1 := GetInReplyTo(e1)
	irt2, ok2 := GetInReplyTo(e2)

	if ok1 != ok2 {
		return fmt.Errorf("THR in-reply-to presence does not fit %v (expected) vs %v", ok2, ok1)
	}

	if ok1 {
		if err := testInReplyToValidator(irt1, irt2); err != nil {
			return err
		}
	}

	total1, ok1 := GetTotal(e1)
	total2, ok2 := GetTotal(e2)

	if ok1 != ok2 {
		return fmt.Errorf("THR total presence does not fit %v (expected) vs %v", ok2, ok1)
	}

	if ok1 {
		if total1.String() != total2.String() {
			return fmt.Errorf("THR total do not match '%s' (expected) vs '%s'", total2.String(), total1.String())
		}
	}

	if len(e1.Links) != len(e2.Links) {
		return fmt.Errorf("Entry does not contain the right count of Links %v (expected) vs %v", len(e2.Links), len(e1.Links))
	}
	for i, _ := range e1.Links {
		e1Count, ok1 := GetCount(e1.Links[i])
		e2Count, ok2 := GetCount(e2.Links[i])

		if ok1 != ok2 {
			return fmt.Errorf("THR count presence does not fit %v (expected) vs %v", ok2, ok1)
		}

		if ok1 {
			if e1Count.String() != e2Count.String() {
				return fmt.Errorf("THR Count do not match '%s' (expected) vs '%s'", e2Count.String(), e1Count.String())
			}
		}

		e1Updated, ok1 := GetUpdated(e1.Links[i])
		e2Updated, ok2 := GetUpdated(e1.Links[i])

		if ok1 != ok2 {
			return fmt.Errorf("THR updated presence does not fit %v (expected) vs %v", ok2, ok1)
		}

		if ok1 {
			if e1Updated.String() != e2Updated.String() {
				return fmt.Errorf("THR Updated do not match '%s' (expected) vs '%s'", e2Updated.String(), e1Updated.String())
			}
		}

	}

	return nil
}

func testThrEntryConstructor() xmlutils.Visitor {
	manager := extension.Manager{}
	AddToManager(&manager)

	return atom.NewEntryExt(manager)
}

func _TestThrEntryToTestVisitor(t testThrEntry) xmlutils.TestVisitor {
	customError := xmlutils.NewErrorChecker(xmlutils.DisableAllError)

	customError.EnableErrorChecking("in-reply-to", atom.MissingAttribute)
	customError.EnableErrorChecking("in-reply-to", atom.AttributeDuplicated)
	customError.EnableErrorChecking("entry", atom.AttributeDuplicated)
	customError.EnableErrorChecking("link", LinkNotReplies)
	customError.EnableErrorChecking("link", NotInLinkElement)
	customError.EnableErrorChecking(xmlutils.AllError, atom.NotPositiveNumber)

	testVisitor := xmlutils.TestVisitor{
		XML:                t.XML,
		ExpectedError:      nil,
		ExpectedVisitor:    t.ExpectedThrEntry,
		VisitorConstructor: testThrEntryConstructor,
		Validator:          testThrEntryValidator,
		CustomError:        &customError,
	}

	if t.ExpectedError != nil {
		testVisitor.ExpectedError = t.ExpectedError
	}

	return testVisitor
}

func TestThrEntryBasic(t *testing.T) {

	var testdata = []testThrEntry{
		{`
     <entry xmlns:thr="http://purl.org/syndication/thread/1.0">
       <link rel="replies" thr:count="10" thr:updated="2003-12-13T18:30:02Z" href="http://example.org/2003/12/13/atom03"/>
       <thr:in-reply-to
         type="application/xhtml+xml"
         href="http://www.example.org/entries/1"
         ref="http://www.example.org/entries/1" />
       <thr:total>44</thr:total>
     </entry>`,
			nil,
			NewTestThrEntry(
				[]*atom.Link{
					NewTestThrLink("2003-12-13T18:30:02Z", "10"),
				},
				"44",
				"http://www.example.org/entries/1",
				"http://www.example.org/entries/1",
				"application/xhtml+xml",
				"",
			),
		},
		{`
     <entry>
       <link rel="replies" thr:count="10" thr:updated="2003-12-13T18:30:02Z" href="http://example.org/2003/12/13/atom03"/>
       <thr:in-reply-to
         type="application/xhtml+xml"
         href="http://www.example.org/entries/1"
         ref="http://www.example.org/entries/1" />
       <thr:total>44</thr:total>
     </entry>`,
			nil,
			NewTestThrEntry(
				[]*atom.Link{
					NewTestThrLink("", ""),
				},
				"",
				"",
				"",
				"",
				"",
			),
		},
		{`
     <entry xmlns:thr="http://purl.org/syndication/thread/1.0">
       <link thr:count="10" thr:updated="2003-12-13T18:30:02Z" href="http://example.org/2003/12/13/atom03"/>
     </entry>`,
			xmlutils.NewError(LinkNotReplies, ""),
			NewTestThrEntry(
				[]*atom.Link{
					NewTestThrLink("2003-12-13T18:30:02Z", "10"),
				},
				"",
				"",
				"",
				"",
				"",
			),
		},
		{`
     <entry xmlns:thr="http://purl.org/syndication/thread/1.0">
       <rights thr:count="10" thr:updated="2003-12-13T18:30:02Z" href="http://example.org/2003/12/13/atom03"/>
       <thr:in-reply-to
         type="application/xhtml+xml"
         href="http://www.example.org/entries/1"
         ref="http://www.example.org/entries/1" />
       <thr:total>44</thr:total>
     </entry>`,
			nil,
			NewTestThrEntry(
				nil,
				"44",
				"http://www.example.org/entries/1",
				"http://www.example.org/entries/1",
				"application/xhtml+xml",
				"",
			),
		},
		{`
     <entry xmlns:thr="http://purl.org/syndication/thread/1.0">
       <link rel="replies" thr:count="1ss0" thr:updated="2003-12-13T18:30:02Z" href="http://example.org/2003/12/13/atom03"/>
     </entry>`,
			xmlutils.NewError(atom.NotPositiveNumber, ""),
			NewTestThrEntry(
				[]*atom.Link{
					NewTestThrLink("2003-12-13T18:30:02Z", "1ss0"),
				},
				"",
				"",
				"",
				"",
				"",
			),
		},
		{`
     <entry xmlns:thr="http://purl.org/syndication/thread/1.0">
       <thr:total>4s4</thr:total>
     </entry>`,
			xmlutils.NewError(atom.NotPositiveNumber, ""),
			NewTestThrEntry(
				nil,
				"4s4",
				"",
				"",
				"",
				"",
			),
		},
		{`
     <entry xmlns:thr="http://purl.org/syndication/thread/1.0">
       <thr:in-reply-to
         type="application/xhtml+xml"
         href="http://www.example.org/entries/1"
         ref="http://www.example.org/entries/1" />
       <thr:in-reply-to
         type="application/xhtml+xml"
         href="http://www.example.org/entries/2"
         ref="http://www.example.org/entries/2" />
     </entry>`,
			xmlutils.NewError(atom.AttributeDuplicated, ""),
			NewTestThrEntry(
				nil,
				"",
				"http://www.example.org/entries/1",
				"http://www.example.org/entries/1",
				"application/xhtml+xml",
				"",
			),
		},
	}

	nbErrors := 0
	len := len(testdata)
	for _, testentry := range testdata {
		testcase := _TestThrEntryToTestVisitor(testentry)

		if err := testcase.CheckTestCase(); err != nil {
			t.Errorf("FAIL\n%s\nXML:\n %s\n", err, testcase.XML)
			nbErrors++
		}
	}
	t.Logf("PASS RATIO = %v/%v\n", len-nbErrors, len)
}
