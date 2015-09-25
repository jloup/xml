package rss

import (
	xmlutils "github.com/jloup/xml/utils"
	"fmt"
	"testing"
)

func NewTestItem(
	title *BasicElement,
	link *BasicElement,
	description *UnescapedContent,
	author *BasicElement,
	cat []*Category,
	comments *BasicElement,
	enclosure *Enclosure,
	guid *Guid,
	pubdate *Date,
	source *Source,
) *Item {

	i := NewItem()

	i.Title = title
	i.Link = link
	i.Description = description
	i.Author = author
	i.Categories = cat
	i.Comments = comments
	i.Enclosure = enclosure
	i.Guid = guid
	i.PubDate = pubdate
	i.Source = source

	return i
}

func NewTestBasicElement(content string) *BasicElement {
	b := NewBasicElement()
	b.Content.Value = content

	return b
}

func NewTestUnescapedContent(content string) *UnescapedContent {
	u := NewUnescapedContent()

	u.Content.WriteString(content)
	return u
}

type testItem struct {
	XML           string
	ExpectedError xmlutils.ParserError
	ExpectedItem  *Item
}

func testItemValidator(actual xmlutils.Visitor, expected xmlutils.Visitor) error {
	i1 := actual.(*Item)
	i2 := expected.(*Item)

	if i1.Title.Content.Value != i2.Title.Content.Value {
		return fmt.Errorf("Title is invalid '%s' (expected) vs '%s'", i2.Title.Content.Value, i1.Title.Content.Value)
	}

	if i1.Link.Content.Value != i2.Link.Content.Value {
		return fmt.Errorf("Link is invalid '%s' (expected) vs '%s'", i2.Link.Content.Value, i1.Link.Content.Value)
	}

	if i1.Description.String() != i2.Description.String() {
		return fmt.Errorf("Description is invalid '%s' (expected) vs '%s'", i2.Description.String(), i1.Description.String())
	}

	if i1.Author.Content.Value != i2.Author.Content.Value {
		return fmt.Errorf("Author is invalid '%s' (expected) vs '%s'", i2.Author.Content.Value, i1.Author.Content.Value)
	}

	if len(i1.Categories) != len(i2.Categories) {
		return fmt.Errorf("Item does not contain the right count of Categories %v (expected) vs %v", len(i2.Categories), len(i1.Categories))
	}

	for i, _ := range i1.Categories {
		if err := testCategoryValidator(i1.Categories[i], i2.Categories[i]); err != nil {
			return err
		}

	}

	if i1.Comments.Content.Value != i2.Comments.Content.Value {
		return fmt.Errorf("Comments is invalid '%s' (expected) vs '%s'", i2.Comments.Content.Value, i1.Comments.Content.Value)
	}

	if err := testEnclosureValidator(i1.Enclosure, i2.Enclosure); err != nil {
		return err
	}

	if err := testGuidValidator(i1.Guid, i2.Guid); err != nil {
		return err
	}

	if i1.PubDate.Time.String() != i2.PubDate.Time.String() {
		return fmt.Errorf("PubDate is invalid '%s' (expected) vs '%s'", i2.PubDate.Time.String(), i1.PubDate.Time.String())
	}

	if err := testSourceValidator(i1.Source, i2.Source); err != nil {
		return err
	}

	return nil
}

func testItemConstructor() xmlutils.Visitor {
	return NewItem()
}

func _TestItemToTestVisitor(t testItem) xmlutils.TestVisitor {
	testVisitor := xmlutils.TestVisitor{
		XML:                t.XML,
		ExpectedError:      nil,
		ExpectedVisitor:    t.ExpectedItem,
		VisitorConstructor: testItemConstructor,
		Validator:          testItemValidator,
	}

	if t.ExpectedError != nil {
		testVisitor.ExpectedError = t.ExpectedError
	}

	return testVisitor
}

func TestItemBasic(t *testing.T) {

	var testdata = []testItem{
		{`
                 <item>
                   <title>Star City</title>
                   <link>http://liftoff.msfc.nasa.gov/news/2003/news-starcity.asp</link>
                   <description>How do Americans get ready to work with Russians aboard the International Space Station? They take a crash course in culture, "yo" language and protocol at Russia's <a href="http://howe.iki.rssi.ru/GCTC/gctc_e.htm">Star City</a>.</description>
                   <pubDate>Tue, 03 Jun 2003 09:39:21 GMT</pubDate>
                   <guid>http://liftoff.msfc.nasa.gov/2003/06/03.html#item573</guid>
		   <category domain="http://www.fool.com/cusips">MSFT</category>
		   <category>MUSIC</category>
                 </item>`,
			nil,
			NewTestItem(
				NewTestBasicElement("Star City"),
				NewTestBasicElement("http://liftoff.msfc.nasa.gov/news/2003/news-starcity.asp"),
				NewTestUnescapedContent("How do Americans get ready to work with Russians aboard the International Space Station? They take a crash course in culture, \"yo\" language and protocol at Russia's <a href=\"http://howe.iki.rssi.ru/GCTC/gctc_e.htm\">Star City</a>."),
				NewBasicElement(),
				[]*Category{
					NewTestCategory("http://www.fool.com/cusips", "MSFT"),
					NewTestCategory("", "MUSIC"),
				},
				NewBasicElement(),
				NewEnclosure(),
				NewTestGuid("true", "http://liftoff.msfc.nasa.gov/2003/06/03.html#item573"),
				NewTestDate("Tue, 03 Jun 2003 09:39:21 GMT"),
				NewSource(),
			),
		},
		{`
                 <item>
		  <title>Star City</title>
                 </item>`,
			nil,
			NewTestItem(
				NewTestBasicElement("Star City"),
				NewBasicElement(),
				NewUnescapedContent(),
				NewBasicElement(),
				nil,
				NewBasicElement(),
				NewEnclosure(),
				NewGuid(),
				NewDate(),
				NewSource(),
			),
		},
		{`
                 <item>
		   <description>How do Americans get ready to work with Russians aboard the International Space Station? They take a crash course in culture, "yo" language and protocol at Russia's <a href="http://howe.iki.rssi.ru/GCTC/gctc_e.htm">Star City</a>.</description>
                 </item>`,
			nil,
			NewTestItem(
				NewBasicElement(),
				NewBasicElement(),
				NewTestUnescapedContent("How do Americans get ready to work with Russians aboard the International Space Station? They take a crash course in culture, \"yo\" language and protocol at Russia's <a href=\"http://howe.iki.rssi.ru/GCTC/gctc_e.htm\">Star City</a>."),
				NewBasicElement(),
				nil,
				NewBasicElement(),
				NewEnclosure(),
				NewGuid(),
				NewDate(),
				NewSource(),
			),
		},
		{`
                 <item>
                 </item>`,
			xmlutils.NewError(MissingAttribute, ""),
			NewTestItem(
				NewBasicElement(),
				NewBasicElement(),
				NewUnescapedContent(),
				NewBasicElement(),
				nil,
				NewBasicElement(),
				NewEnclosure(),
				NewGuid(),
				NewDate(),
				NewSource(),
			),
		},
		{`
                 <item>
		  <title>Star City</title>
                  <title>Star City 2</title>
		  </item>`,
			xmlutils.NewError(AttributeDuplicated, ""),
			NewTestItem(
				NewTestBasicElement("Star City 2"),
				NewBasicElement(),
				NewUnescapedContent(),
				NewBasicElement(),
				nil,
				NewBasicElement(),
				NewEnclosure(),
				NewGuid(),
				NewDate(),
				NewSource(),
			),
		},
	}

	nbErrors := 0
	len := len(testdata)
	for _, testitem := range testdata {
		testcase := _TestItemToTestVisitor(testitem)

		if err := testcase.CheckTestCase(); err != nil {
			t.Errorf("FAIL\n%s\nXML:\n %s\n", err, testcase.XML)
			nbErrors++
		}
	}

	t.Logf("PASS RATIO = %v/%v\n", len-nbErrors, len)
}
