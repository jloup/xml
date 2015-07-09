package rss

import (
	"fmt"
	"testing"

	"github.com/jloup/xml/utils"
)

func NewTestChannel(
	title *BasicElement,
	link *BasicElement,
	description *UnescapedContent,
	language *BasicElement,
	copyright *BasicElement,
	managingEditor *BasicElement,
	webmaster *BasicElement,
	pubdate *Date,
	lastbuilddate *Date,
	cat []*Category,
	generator *BasicElement,
	docs *BasicElement,
	cloud *Cloud,
	ttl *BasicElement,
	image *Image,
	rating *BasicElement,
	skiphours *BasicElement,
	skipdays *BasicElement,
	items []*Item,
) *Channel {

	c := NewChannel()

	c.Title = title
	c.Link = link
	c.Description = description
	c.Language = language
	c.Copyright = copyright
	c.ManagingEditor = managingEditor
	c.Webmaster = webmaster
	c.PubDate = pubdate
	c.LastBuildDate = lastbuilddate
	c.Categories = cat
	c.Generator = generator
	c.Docs = docs
	c.Cloud = cloud
	c.Ttl = ttl
	c.Image = image
	c.Rating = rating
	c.SkipHours = skiphours
	c.SkipDays = skipdays
	c.Items = items

	return c
}

type testChannel struct {
	XML             string
	ExpectedError   utils.ParserError
	ExpectedChannel *Channel
}

func testChannelValidator(actual utils.Visitor, expected utils.Visitor) error {
	c1 := actual.(*Channel)
	c2 := expected.(*Channel)

	if c1.Title.Content.Value != c2.Title.Content.Value {
		return fmt.Errorf("Title is invalid '%s' (expected) vs '%s'", c2.Title.Content.Value, c1.Title.Content.Value)
	}

	if c1.Link.Content.Value != c2.Link.Content.Value {
		return fmt.Errorf("Link is invalid '%s' (expected) vs '%s'", c2.Link.Content.Value, c1.Link.Content.Value)
	}

	if c1.Language.String() != c2.Language.String() {
		return fmt.Errorf("Language is invalid '%s' (expected) vs '%s'", c2.Language.String(), c1.Language.String())
	}

	if c1.Copyright.String() != c2.Copyright.String() {
		return fmt.Errorf("Copyright is invalid '%s' (expected) vs '%s'", c2.Copyright.String(), c1.Copyright.String())
	}

	if c1.ManagingEditor.String() != c2.ManagingEditor.String() {
		return fmt.Errorf("ManagingEditor is invalid '%s' (expected) vs '%s'", c2.ManagingEditor.String(), c1.ManagingEditor.String())
	}

	if c1.Webmaster.String() != c2.Webmaster.String() {
		return fmt.Errorf("Webmaster is invalid '%s' (expected) vs '%s'", c2.Webmaster.String(), c1.Webmaster.String())
	}

	if err := testDateValidator(c1.PubDate, c2.PubDate); err != nil {
		return err
	}

	if err := testDateValidator(c1.LastBuildDate, c2.LastBuildDate); err != nil {
		return err
	}

	if len(c1.Categories) != len(c2.Categories) {
		return fmt.Errorf("Channel does not contain the right count of Categories %v (expected) vs %v", len(c2.Categories), len(c1.Categories))
	}

	for i, _ := range c1.Categories {
		if err := testCategoryValidator(c1.Categories[i], c2.Categories[i]); err != nil {
			return err
		}

	}

	if c1.Generator.String() != c2.Generator.String() {
		return fmt.Errorf("Generator is invalid '%s' (expected) vs '%s'", c2.Generator.String(), c1.Generator.String())
	}

	if c1.Docs.String() != c2.Docs.String() {
		return fmt.Errorf("Docs is invalid '%s' (expected) vs '%s'", c2.Docs.String(), c1.Docs.String())
	}

	if err := testCloudValidator(c1.Cloud, c2.Cloud); err != nil {
		return err
	}

	if c1.Ttl.String() != c2.Ttl.String() {
		return fmt.Errorf("Ttl is invalid '%s' (expected) vs '%s'", c2.Ttl.String(), c1.Ttl.String())
	}

	if err := testImageValidator(c1.Image, c2.Image); err != nil {
		return err
	}

	if c1.Rating.String() != c2.Rating.String() {
		return fmt.Errorf("Rating is invalid '%s' (expected) vs '%s'", c2.Rating.String(), c1.Rating.String())
	}

	if c1.SkipHours.String() != c2.SkipHours.String() {
		return fmt.Errorf("SkipHours is invalid '%s' (expected) vs '%s'", c2.SkipHours.String(), c1.SkipHours.String())
	}

	if c1.SkipDays.String() != c2.SkipDays.String() {
		return fmt.Errorf("SkipDays is invalid '%s' (expected) vs '%s'", c2.SkipDays.String(), c1.SkipDays.String())
	}

	if len(c1.Items) != len(c2.Items) {
		return fmt.Errorf("Channel does not contain the right count of Items %v (expected) vs %v", len(c2.Items), len(c1.Items))
	}

	for i, _ := range c1.Items {
		if err := testItemValidator(c1.Items[i], c2.Items[i]); err != nil {
			return err
		}

	}
	return nil
}

func testChannelConstructor() utils.Visitor {
	return NewChannel()
}

func _TestChannelToTestVisitor(t testChannel) utils.TestVisitor {
	testVisitor := utils.TestVisitor{
		XML:                t.XML,
		ExpectedError:      nil,
		ExpectedVisitor:    t.ExpectedChannel,
		VisitorConstructor: testChannelConstructor,
		Validator:          testChannelValidator,
	}

	if t.ExpectedError != nil {
		testVisitor.ExpectedError = t.ExpectedError
	}

	return testVisitor
}

func TestChannelBasic(t *testing.T) {

	var testdata = []testChannel{
		{`
                 <channel>
                   <title>Liftoff News</title>
                   <link>http://liftoff.msfc.nasa.gov/</link>
                   <description>Liftoff to Space Exploration.<a>yo'.com</a></description>
                   <language>en-us</language>
                   <pubDate>Tue, 10 Jun 2003 04:00:00 GMT</pubDate>
                   <lastBuildDate>Tue, 10 Jun 2003 09:41:01 GMT</lastBuildDate>
                   <docs>http://blogs.law.harvard.edu/tech/rss</docs>
                   <generator>Weblog Editor 2.0</generator>
                   <managingEditor>editor@example.com</managingEditor>
                   <webMaster>webmaster@example.com</webMaster>
		 </channel>`,
			nil,
			NewTestChannel(
				NewTestBasicElement("Liftoff News"),
				NewTestBasicElement("http://liftoff.msfc.nasa.gov/"),
				NewTestUnescapedContent("Liftoff to Space Exploration.<a>yo'.com</a>"),
				NewTestBasicElement("en-us"),
				NewTestBasicElement(""),
				NewTestBasicElement("editor@example.com"),
				NewTestBasicElement("webmaster@example.com"),
				NewTestDate("Tue, 10 Jun 2003 04:00:00 GMT"),
				NewTestDate("Tue, 10 Jun 2003 09:41:01 GMT"),
				nil,
				NewTestBasicElement("Weblog Editor 2.0"),
				NewTestBasicElement("http://blogs.law.harvard.edu/tech/rss"),
				NewCloud(),
				NewBasicElement(),
				NewImage(),
				NewBasicElement(),
				NewBasicElement(),
				NewBasicElement(),
				nil,
			),
		},
		{`
                 <channel>
                   <title>Liftoff News</title>
                   <link>http://liftoff.msfc.nasa.gov/</link>
                   <description>Liftoff to Space Exploration.<a>yo'.com</a></description>
	         </channel>`,
			nil,
			NewTestChannel(
				NewTestBasicElement("Liftoff News"),
				NewTestBasicElement("http://liftoff.msfc.nasa.gov/"),
				NewTestUnescapedContent("Liftoff to Space Exploration.<a>yo'.com</a>"),
				NewBasicElement(),
				NewBasicElement(),
				NewBasicElement(),
				NewBasicElement(),
				NewDate(),
				NewDate(),
				nil,
				NewBasicElement(),
				NewBasicElement(),
				NewCloud(),
				NewBasicElement(),
				NewImage(),
				NewBasicElement(),
				NewBasicElement(),
				NewBasicElement(),
				nil,
			),
		},
		{`
                 <channel>
                   <link>http://liftoff.msfc.nasa.gov/</link>
                   <description>Liftoff to Space Exploration.<a>yo'.com</a></description>
	         </channel>`,
			utils.NewError(MissingAttribute, ""),
			NewTestChannel(
				NewBasicElement(),
				NewTestBasicElement("http://liftoff.msfc.nasa.gov/"),
				NewTestUnescapedContent("Liftoff to Space Exploration.<a>yo'.com</a>"),
				NewBasicElement(),
				NewBasicElement(),
				NewBasicElement(),
				NewBasicElement(),
				NewDate(),
				NewDate(),
				nil,
				NewBasicElement(),
				NewBasicElement(),
				NewCloud(),
				NewBasicElement(),
				NewImage(),
				NewBasicElement(),
				NewBasicElement(),
				NewBasicElement(),
				nil,
			),
		}, {`
                 <channel>
                   <title>Liftoff News</title>
                   <description>Liftoff to Space Exploration.<a>yo'.com</a></description>
	         </channel>`,
			utils.NewError(MissingAttribute, ""),
			NewTestChannel(
				NewTestBasicElement("Liftoff News"),
				NewBasicElement(),
				NewTestUnescapedContent("Liftoff to Space Exploration.<a>yo'.com</a>"),
				NewBasicElement(),
				NewBasicElement(),
				NewBasicElement(),
				NewBasicElement(),
				NewDate(),
				NewDate(),
				nil,
				NewBasicElement(),
				NewBasicElement(),
				NewCloud(),
				NewBasicElement(),
				NewImage(),
				NewBasicElement(),
				NewBasicElement(),
				NewBasicElement(),
				nil,
			),
		},
		{`
                 <channel>
                   <title>Liftoff News</title>
                   <link>http://liftoff.msfc.nasa.gov/</link>
	         </channel>`,
			utils.NewError(MissingAttribute, ""),
			NewTestChannel(
				NewTestBasicElement("Liftoff News"),
				NewTestBasicElement("http://liftoff.msfc.nasa.gov/"),
				NewUnescapedContent(),
				NewBasicElement(),
				NewBasicElement(),
				NewBasicElement(),
				NewBasicElement(),
				NewDate(),
				NewDate(),
				nil,
				NewBasicElement(),
				NewBasicElement(),
				NewCloud(),
				NewBasicElement(),
				NewImage(),
				NewBasicElement(),
				NewBasicElement(),
				NewBasicElement(),
				nil,
			),
		},
		{`
                 <channel>
                   <title>Liftoff News</title>
                   <link>http://liftoff.msfc.nasa.gov/</link>
		   <item>
		     <title>Star City</title>
                   </item>
	         </channel>`,
			utils.NewError(MissingAttribute, ""),
			NewTestChannel(
				NewTestBasicElement("Liftoff News"),
				NewTestBasicElement("http://liftoff.msfc.nasa.gov/"),
				NewUnescapedContent(),
				NewBasicElement(),
				NewBasicElement(),
				NewBasicElement(),
				NewBasicElement(),
				NewDate(),
				NewDate(),
				nil,
				NewBasicElement(),
				NewBasicElement(),
				NewCloud(),
				NewBasicElement(),
				NewImage(),
				NewBasicElement(),
				NewBasicElement(),
				NewBasicElement(),
				[]*Item{
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
			),
		},
	}

	nbErrors := 0
	len := len(testdata)
	for _, testchannel := range testdata {
		testcase := _TestChannelToTestVisitor(testchannel)

		if err := testcase.CheckTestCase(); err != nil {
			t.Errorf("FAIL\n%s\nXML:\n %s\n", err, testcase.XML)
			nbErrors++
		}
	}

	t.Logf("PASS RATIO = %v/%v\n", len-nbErrors, len)
}
