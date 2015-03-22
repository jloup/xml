package rss

import (
	"fmt"
	"testing"
	"time"

	"github.com/JLoup/xml/helper"
)

func NewTestDate(date string) *Date {
	d := NewDate()
	d.Time, _ = time.Parse(rssDateFormat[0], date)

	return d
}

type testDate struct {
	XML           string
	ExpectedError helper.ParserError
	ExpectedDate  *Date
}

func testDateValidator(actual helper.Visitor, expected helper.Visitor) error {
	d1 := actual.(*Date)
	d2 := expected.(*Date)

	if d1.Time.Equal(d2.Time) == false {
		return fmt.Errorf("Time is invalid '%s' (expected) vs '%s'", d2.Time.String(), d1.Time.String())
	}

	return nil
}

func testDateConstructor() helper.Visitor {
	return NewDate()
}

func testDateToTestVisitor(t testDate) helper.TestVisitor {
	testVisitor := helper.TestVisitor{
		XML:                t.XML,
		ExpectedError:      nil,
		ExpectedVisitor:    t.ExpectedDate,
		VisitorConstructor: testDateConstructor,
		Validator:          testDateValidator,
	}

	if t.ExpectedError != nil {
		testVisitor.ExpectedError = t.ExpectedError
	}

	return testVisitor
}

func TestDateBasic(t *testing.T) {

	var testdata = []testDate{
		{`<updated>Tue, 20 Sep 2010 16:02:50 GMT</updated>`,
			nil,
			NewTestDate("Tue, 20 Sep 2010 16:02:50 GMT"),
		},
		{`<updated>Tue, 20 Sep 2010 16:02:50</updated>`,
			helper.NewError(DateFormat, ""),
			NewTestDate("0"),
		},
		{`<updated>2003-12-13T18:30:02.25</updated>`,
			helper.NewError(DateFormat, ""),
			NewTestDate("0"),
		},
	}
	nbErrors := 0
	len := len(testdata)
	for _, testdate := range testdata {
		testcase := testDateToTestVisitor(testdate)

		if err := testcase.CheckTestCase(); err != nil {
			t.Errorf("FAIL\n%s\nXML:\n %s\n", err, testcase.XML)
			nbErrors++
		}
	}

	t.Logf("PASS RATIO = %v/%v\n", len-nbErrors, len)
}
