package atom

import (
	"fmt"
	"testing"
	"time"

	"github.com/JLoup/xml/utils"
)

func NewTestDate(date string) *Date {
	d := NewDate()
	d.Time, _ = time.Parse(time.RFC3339, date)

	return d
}

func DateWithBaseLang(date *Date, lang, base string) *Date {
	date.Lang.Value = lang
	date.Base.Value = base

	return date
}

type testDate struct {
	XML           string
	ExpectedError utils.ParserError
	ExpectedDate  *Date
}

func testDateValidator(actual utils.Visitor, expected utils.Visitor) error {
	d1 := actual.(*Date)
	d2 := expected.(*Date)

	if d1.Time.Equal(d2.Time) == false {
		return fmt.Errorf("Time is invalid '%s' (expected) vs '%s'", d2.Time.String(), d1.Time.String())
	}

	if err := ValidateBaseLang("Date", d1.Base.Value, d1.Lang.Value, d2.Base.Value, d2.Lang.Value); err != nil {
		return err
	}

	return nil
}

func testDateConstructor() utils.Visitor {
	return NewDate()
}

func testDateToTestVisitor(t testDate) utils.TestVisitor {
	testVisitor := utils.TestVisitor{
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
		{`<updated xml:lang="en-us" xml:base="http://yo.com">2003-12-13T18:30:02Z</updated>`,
			nil,
			DateWithBaseLang(NewTestDate("2003-12-13T18:30:02Z"), "en-us", "http://yo.com"),
		},
		{`<updated>2003-12-13T18:30:02Z</updated>`,
			nil,
			NewTestDate("2003-12-13T18:30:02Z"),
		},
		{`<updated>2003-12-13T18:30:02.25Z</updated>`,
			nil,
			NewTestDate("2003-12-13T18:30:02.25Z"),
		},
		{`<updated>2003-12-13T18:30:02+01:00</updated>`,
			nil,
			NewTestDate("2003-12-13T18:30:02+01:00"),
		},
		{`<updated>2003-12-13T18:30:02.25+01:00</updated>`,
			nil,
			NewTestDate("2003-12-13T18:30:02.25+01:00"),
		},
		{`<updated>2003-12-13t18:30:02z</updated>`,
			utils.NewError(DateFormat, ""),
			NewTestDate("0"),
		},
		{`<updated>2003-12-13T18:30:02.25</updated>`,
			utils.NewError(DateFormat, ""),
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
