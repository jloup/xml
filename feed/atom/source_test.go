package atom

import (
	"fmt"
	"testing"

	"github.com/JLoup/xml/utils"
)

func NewTestSource(
	authors []*Person,
	cat []*Category,
	contributors []*Person,
	generator *Generator,
	icon *Icon,
	id *Id,
	links []*Link,
	logo *Logo,
	rights, sub, title *TextConstruct,
	updated *Date,
) *Source {

	s := NewSource()

	s.Authors = authors
	s.Categories = cat
	s.Contributors = contributors
	s.Generator = generator
	s.Icon = icon
	s.Id = id
	s.Links = links
	s.Logo = logo
	s.Rights = rights
	s.Subtitle = sub
	s.Title = title
	s.Updated = updated

	return s
}

func SourceWithBaseLang(s *Source, lang, base string) *Source {
	s.Lang.Value = lang
	s.Base.Value = base

	return s
}

type testSource struct {
	XML            string
	ExpectedError  utils.ParserError
	ExpectedSource *Source
}

func testSourceValidator(actual utils.Visitor, expected utils.Visitor) error {
	s1 := actual.(*Source)
	s2 := expected.(*Source)

	if len(s1.Authors) != len(s2.Authors) {
		return fmt.Errorf("Source does not contain the right count of Authors %v (expected) vs %v", len(s2.Authors), len(s1.Authors))
	}

	for i, _ := range s1.Authors {
		if err := testPersonValidator(s1.Authors[i], s2.Authors[i]); err != nil {
			return err
		}

	}

	if len(s1.Categories) != len(s2.Categories) {
		return fmt.Errorf("Source does not contain the right count of Categories %v (expected) vs %v", len(s2.Categories), len(s1.Categories))
	}

	for i, _ := range s1.Categories {
		if err := testCategoryValidator(s1.Categories[i], s2.Categories[i]); err != nil {
			return err
		}

	}

	if len(s1.Contributors) != len(s2.Contributors) {
		return fmt.Errorf("Source does not contain the right count of Contributors %v (expected) vs %v", len(s2.Contributors), len(s1.Contributors))
	}

	for i, _ := range s1.Contributors {
		if err := testPersonValidator(s1.Contributors[i], s2.Contributors[i]); err != nil {
			return err
		}

	}

	if err := testGeneratorValidator(s1.Generator, s2.Generator); err != nil {
		return err
	}

	if err := testIconValidator(s1.Icon, s2.Icon); err != nil {
		return err
	}

	if err := testIdValidator(s1.Id, s2.Id); err != nil {
		return err
	}

	if len(s1.Links) != len(s2.Links) {
		return fmt.Errorf("Source does not contain the right count of Links %v (expected) vs %v", len(s2.Links), len(s1.Links))
	}

	for i, _ := range s1.Links {
		if err := testLinkValidator(s1.Links[i], s2.Links[i]); err != nil {
			return err
		}

	}

	if err := testLogoValidator(s1.Logo, s2.Logo); err != nil {
		return err
	}

	if err := testTextConstructValidator(s1.Rights, s2.Rights); err != nil {
		return err
	}

	if err := testTextConstructValidator(s1.Subtitle, s2.Subtitle); err != nil {
		return err
	}

	if err := testTextConstructValidator(s1.Title, s2.Title); err != nil {
		return err
	}

	if err := testDateValidator(s1.Updated, s2.Updated); err != nil {
		return err
	}

	if err := ValidateBaseLang("Source", s1.Base.Value, s1.Lang.Value, s2.Base.Value, s2.Lang.Value); err != nil {
		return err
	}

	return nil
}

func testSourceConstructor() utils.Visitor {
	return NewSource()
}

func _TestSourceToTestVisitor(t testSource) utils.TestVisitor {
	testVisitor := utils.TestVisitor{
		XML:                t.XML,
		ExpectedError:      nil,
		ExpectedVisitor:    t.ExpectedSource,
		VisitorConstructor: testSourceConstructor,
		Validator:          testSourceValidator,
	}

	if t.ExpectedError != nil {
		testVisitor.ExpectedError = t.ExpectedError
	}

	return testVisitor
}

func TestSourceBasic(t *testing.T) {

	var testdata = []testSource{
		{`
  <source xml:lang="en-us" xml:base="http://yo.com" xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <id>tag:example.org,2003:3</id>  
  </source>`,
			nil,
			SourceWithBaseLang(NewTestSource(
				nil,
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				nil,
				NewLogo(),
				NewTextConstruct(),
				NewTextConstruct(),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
			), "en-us", "http://yo.com"),
		},
		{`
  <source xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <id>tag:example.org,2003:3</id>  
    <author>
      <name>jean</name>
    </author>
    <author>
      <name>pierre</name>
    </author>
  </source>`,
			nil,
			NewTestSource(
				[]*Person{
					NewTestPerson("jean", "", ""),
					NewTestPerson("pierre", "", ""),
				},
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				nil,
				NewLogo(),
				NewTextConstruct(),
				NewTextConstruct(),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
			),
		},
		{`
  <source xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <id>tag:example.org,2003:3</id>  
    <category term='music' />
    <category term='book' />
  </source>`,
			nil,
			NewTestSource(
				nil,
				[]*Category{
					NewTestCategory("", "music", ""),
					NewTestCategory("", "book", ""),
				},
				nil,
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				nil,
				NewLogo(),
				NewTextConstruct(),
				NewTextConstruct(),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
			),
		},
		{`
  <source xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <id>tag:example.org,2003:3</id>  
    <contributor>
      <name>jean</name>
    </contributor>
    <contributor>
      <name>pierre</name>
    </contributor>
  </source>`,
			nil,
			NewTestSource(
				nil,
				nil,
				[]*Person{
					NewTestPerson("jean", "", ""),
					NewTestPerson("pierre", "", ""),
				},
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				nil,
				NewLogo(),
				NewTextConstruct(),
				NewTextConstruct(),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
			),
		},
		{`
  <source xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <id>tag:example.org,2003:3</id>  
    <generator uri="http://there.com" version="4.0">my generator</generator>
    <generator uri="http://there2.com" version="5.0">my generator bis</generator>
  </source>`,
			utils.NewError(AttributeDuplicated, ""),
			NewTestSource(
				nil,
				nil,
				nil,
				NewTestGenerator("5.0", "http://there2.com", "my generator bis"),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				nil,
				NewLogo(),
				NewTextConstruct(),
				NewTextConstruct(),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
			),
		},
		{`
  <source xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <id>tag:example.org,2003:3</id>  
    <icon>https://www.yo.com</icon>
    <icon>https://www.yo2.com</icon>
  </source>`,
			utils.NewError(AttributeDuplicated, ""),
			NewTestSource(
				nil,
				nil,
				nil,
				NewGenerator(),
				NewTestIcon("https://www.yo2.com"),
				NewTestId("tag:example.org,2003:3"),
				nil,
				NewLogo(),
				NewTextConstruct(),
				NewTextConstruct(),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
			),
		},
		{`
  <source xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <id>tag:example.org,2003:3</id>  
    <id>tag:example.org,2010:3</id>
  </source>`,
			utils.NewError(AttributeDuplicated, ""),
			NewTestSource(
				nil,
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2010:3"),
				nil,
				NewLogo(),
				NewTextConstruct(),
				NewTextConstruct(),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
			),
		},
		{`
  <source xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
  </source>`,
			utils.NewError(MissingAttribute, ""),
			NewTestSource(
				nil,
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewId(),
				nil,
				NewLogo(),
				NewTextConstruct(),
				NewTextConstruct(),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
			),
		},
		{`
  <source xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <id>tag:example.org,2003:3</id>  
    <logo>https://www.yo.com</logo>
    <logo>https://www.yo2.com</logo>
  </source>`,
			utils.NewError(AttributeDuplicated, ""),
			NewTestSource(
				nil,
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				nil,
				NewTestLogo("https://www.yo2.com"),
				NewTextConstruct(),
				NewTextConstruct(),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
			),
		},
		{`
  <source xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <id>tag:example.org,2003:3</id> 
    <rights>Copyright (c) 2003, Mark Pilgrim</rights>
    <rights>Copyright (c) 2006, Mark Pilgrim</rights>
  </source>`,
			utils.NewError(AttributeDuplicated, ""),
			NewTestSource(
				nil,
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				nil,
				NewLogo(),
				NewTestTextConstruct("", "Copyright (c) 2006, Mark Pilgrim"),
				NewTextConstruct(),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
			),
		},
		{`
  <source xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <id>tag:example.org,2003:3</id> 
    <subtitle type="html">First</subtitle>
    <subtitle type="html">Second</subtitle>
  </source>`,
			utils.NewError(AttributeDuplicated, ""),
			NewTestSource(
				nil,
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				nil,
				NewLogo(),
				NewTextConstruct(),
				NewTestTextConstruct("", "Second"),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
			),
		},
		{`
  <source xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <id>tag:example.org,2003:3</id> 
    <title type="html">Second title</title>
  </source>`,
			utils.NewError(AttributeDuplicated, ""),
			NewTestSource(
				nil,
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				nil,
				NewLogo(),
				NewTextConstruct(),
				NewTextConstruct(),
				NewTestTextConstruct("", "Second title"),
				NewTestDate("2005-07-31T12:29:29Z"),
			),
		},
		{`
  <source xmlns="http://www.w3.org/2005/Atom">
    <updated>2005-07-31T12:29:29Z</updated>
    <id>tag:example.org,2003:3</id> 
  </source>`,
			utils.NewError(MissingAttribute, ""),
			NewTestSource(
				nil,
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				nil,
				NewLogo(),
				NewTextConstruct(),
				NewTextConstruct(),
				NewTextConstruct(),
				NewTestDate("2005-07-31T12:29:29Z"),
			),
		},
		{`
  <source xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <id>tag:example.org,2003:3</id> 
  </source>`,
			utils.NewError(MissingAttribute, ""),
			NewTestSource(
				nil,
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				nil,
				NewLogo(),
				NewTextConstruct(),
				NewTextConstruct(),
				NewTestTextConstruct("", "dive into mark"),
				NewDate(),
			),
		},
		{`
  <source xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <updated>2010-07-31T12:29:29Z</updated>
    <id>tag:example.org,2003:3</id> 
  </source>`,
			utils.NewError(AttributeDuplicated, ""),
			NewTestSource(
				nil,
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				nil,
				NewLogo(),
				NewTextConstruct(),
				NewTextConstruct(),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2010-07-31T12:29:29Z"),
			),
		},
		{`
  <source xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <id>tag:example.org,2003:3</id> 
    <entry></entry>
  </source>`,
			utils.NewError(AttributeForbidden, ""),
			NewTestSource(
				nil,
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				nil,
				NewLogo(),
				NewTextConstruct(),
				NewTextConstruct(),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
			),
		},
		{`
  <source xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <id>tag:example.org,2003:3</id> 
    <link href="http://example.com" hreflang="en" type="application/xhtml+xml"/>
    <link href="http://example.com" hreflang="en" type="application/xhtml+xml"/>
  </source>`,
			utils.NewError(LinkAlternateDuplicated, ""),
			NewTestSource(
				nil,
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				[]*Link{
					NewTestLink("http://example.com", "alternate", "application/xhtml+xml", "en", "", ""),
					NewTestLink("http://example.com", "alternate", "application/xhtml+xml", "en", "", ""),
				},
				NewLogo(),
				NewTextConstruct(),
				NewTextConstruct(),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
			),
		},
		{`
  <source xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <id>tag:example.org,2003:3</id> 
    <link href="http://example.com" type="application/xhtml+xml"/>
    <link href="http://example.com" type="application/xhtml+xml"/>
  </source>`,
			utils.NewError(LinkAlternateDuplicated, ""),
			NewTestSource(
				nil,
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				[]*Link{
					NewTestLink("http://example.com", "alternate", "application/xhtml+xml", "", "", ""),
					NewTestLink("http://example.com", "alternate", "application/xhtml+xml", "", "", ""),
				},
				NewLogo(),
				NewTextConstruct(),
				NewTextConstruct(),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
			),
		},
		{`
  <source xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <id>tag:example.org,2003:3</id> 
    <link href="http://example.com" hreflang="en" type="application/xhtml+xml"/>
    <link href="http://example.com" hreflang="fr" type="application/xhtml+xml"/>
  </source>`,
			nil,
			NewTestSource(
				nil,
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				[]*Link{
					NewTestLink("http://example.com", "alternate", "application/xhtml+xml", "en", "", ""),
					NewTestLink("http://example.com", "alternate", "application/xhtml+xml", "fr", "", ""),
				},
				NewLogo(),
				NewTextConstruct(),
				NewTextConstruct(),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
			),
		},
	}

	nbErrors := 0
	len := len(testdata)
	for _, testsource := range testdata {
		testcase := _TestSourceToTestVisitor(testsource)

		if err := testcase.CheckTestCase(); err != nil {
			t.Errorf("FAIL\n%s\nXML:\n %s\n", err, testcase.XML)
			nbErrors++
		}
	}

	t.Logf("PASS RATIO = %v/%v\n", len-nbErrors, len)
}
