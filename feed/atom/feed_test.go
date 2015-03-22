package atom

import (
	"fmt"
	"testing"

	"github.com/JLoup/xml/helper"
)

func NewTestFeed(
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
	entries []*Entry,
) *Feed {

	f := NewFeed()

	f.Authors = authors
	f.Categories = cat
	f.Contributors = contributors
	f.Generator = generator
	f.Icon = icon
	f.Id = id
	f.Links = links
	f.Logo = logo
	f.Rights = rights
	f.Subtitle = sub
	f.Title = title
	f.Updated = updated
	f.Entries = entries

	return f
}

func FeedWithBaseLang(f *Feed, lang, base string) *Feed {
	f.Lang.Value = lang
	f.Base.Value = base

	return f
}

type testFeed struct {
	XML           string
	ExpectedError helper.ParserError
	ExpectedFeed  *Feed
}

func testFeedValidator(actual helper.Visitor, expected helper.Visitor) error {
	f1 := actual.(*Feed)
	f2 := expected.(*Feed)

	if len(f1.Authors) != len(f2.Authors) {
		return fmt.Errorf("Feed does not contain the right count of Authors %v (expected) vs %v", len(f2.Authors), len(f1.Authors))
	}

	for i, _ := range f1.Authors {
		if err := testPersonValidator(f1.Authors[i], f2.Authors[i]); err != nil {
			return err
		}

	}

	if len(f1.Categories) != len(f2.Categories) {
		return fmt.Errorf("Feed does not contain the right count of Categories %v (expected) vs %v", len(f2.Categories), len(f1.Categories))
	}

	for i, _ := range f1.Categories {
		if err := testCategoryValidator(f1.Categories[i], f2.Categories[i]); err != nil {
			return err
		}

	}

	if len(f1.Contributors) != len(f2.Contributors) {
		return fmt.Errorf("Feed does not contain the right count of Contributors %v (expected) vs %v", len(f2.Contributors), len(f1.Contributors))
	}

	for i, _ := range f1.Contributors {
		if err := testPersonValidator(f1.Contributors[i], f2.Contributors[i]); err != nil {
			return err
		}

	}

	if err := testGeneratorValidator(f1.Generator, f2.Generator); err != nil {
		return err
	}

	if err := testIconValidator(f1.Icon, f2.Icon); err != nil {
		return err
	}

	if err := testIdValidator(f1.Id, f2.Id); err != nil {
		return err
	}

	if len(f1.Links) != len(f2.Links) {
		return fmt.Errorf("Feed does not contain the right count of Links %v (expected) vs %v", len(f2.Links), len(f1.Links))
	}

	for i, _ := range f1.Links {
		if err := testLinkValidator(f1.Links[i], f2.Links[i]); err != nil {
			return err
		}

	}

	if err := testLogoValidator(f1.Logo, f2.Logo); err != nil {
		return err
	}

	if err := testTextConstructValidator(f1.Rights, f2.Rights); err != nil {
		return err
	}

	if err := testTextConstructValidator(f1.Subtitle, f2.Subtitle); err != nil {
		return err
	}

	if err := testTextConstructValidator(f1.Title, f2.Title); err != nil {
		return err
	}

	if err := testDateValidator(f1.Updated, f2.Updated); err != nil {
		return err
	}

	if len(f1.Entries) != len(f2.Entries) {
		return fmt.Errorf("Feed does not contain the right count of Entries %v (expected) vs %v", len(f2.Entries), len(f1.Entries))
	}

	for i, _ := range f1.Entries {
		if err := testEntryValidator(f1.Entries[i], f2.Entries[i]); err != nil {
			return err
		}

	}

	if err := ValidateBaseLang("Feed", f1.Base.Value, f1.Lang.Value, f2.Base.Value, f2.Lang.Value); err != nil {
		return err
	}

	return nil
}

func testFeedConstructor() helper.Visitor {
	return NewFeed()
}

func _TestFeedToTestVisitor(t testFeed) helper.TestVisitor {
	testVisitor := helper.TestVisitor{
		XML:                t.XML,
		ExpectedError:      nil,
		ExpectedVisitor:    t.ExpectedFeed,
		VisitorConstructor: testFeedConstructor,
		Validator:          testFeedValidator,
	}

	if t.ExpectedError != nil {
		testVisitor.ExpectedError = t.ExpectedError
	}

	return testVisitor
}

func TestFeedBasic(t *testing.T) {

	var testdata = []testFeed{
		{`
  <feed xml:lang="en-us" xml:base="http://yo.com" xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <subtitle type="html">A &lt;em&gt;lot&lt;/em&gt; of effort went into making this effortless</subtitle>
    <updated>2005-07-31T12:29:29Z</updated>
    <id>tag:example.org,2003:3</id>
    <link rel="alternate" type="text/html"
     hreflang="en" href="http://example.org/"/>
    <link rel="self" type="application/atom+xml"
     href="http://example.org/feed.atom"/>
    <rights>Copyright (c) 2003, Mark Pilgrim</rights>
    <generator uri="http://www.example.com/" version="1.0">Example Toolkit</generator>
    <entry>
      <title>Atom draft-07 snapshot</title>
      <link rel="alternate" type="text/html"
       href="http://example.org/2005/04/02/atom"/>
      <link rel="enclosure" type="audio/mpeg" length="1337"
       href="http://example.org/audio/ph34r_my_podcast.mp3"/>
      <id>tag:example.org,2003:3.2397</id>
      <updated>2005-07-31T12:29:29Z</updated>
      <published>2003-12-13T08:29:29-04:00</published>
      <author>
        <name>Mark Pilgrim</name>
        <uri>http://example.org/</uri>
        <email>f8dy@example.com</email>
      </author>
      <contributor>
        <name>Sam Ruby</name>
      </contributor>
      <contributor>
        <name>Joe Gregorio</name>
      </contributor>
      <content type="xhtml" xml:lang="en"
       xml:base="http://diveintomark.org/">
        <div xmlns="http://www.w3.org/1999/xhtml">
          <p><i>[Update: The Atom draft is finished.]</i></p>
        </div>
      </content>
    </entry>
  </feed>`,
			nil,
			FeedWithBaseLang(NewTestFeed(
				nil,
				nil,
				nil,
				NewTestGenerator("1.0", "http://www.example.com/", "Example Toolkit"),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				[]*Link{
					NewTestLink("http://example.org/", "alternate", "text/html", "en", "", ""),
					NewTestLink("http://example.org/feed.atom", "self", "application/atom+xml", "", "", ""),
				},
				NewLogo(),
				NewTestTextConstruct("", "Copyright (c) 2003, Mark Pilgrim"),
				NewTestTextConstruct("", "A <em>lot</em> of effort went into making this effortless"),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
				[]*Entry{
					NewTestEntry(
						[]*Person{
							NewTestPerson("Mark Pilgrim", "http://example.org/", "f8dy@example.com"),
						},
						nil,
						ContentWithBaseLang(NewTestContent("xhtml", `<div><p><i>[Update: The Atom draft is finished.]</i></p></div>`, "", "", ""), "en", "http://diveintomark.org/"),
						[]*Person{
							NewTestPerson("Sam Ruby", "", ""),
							NewTestPerson("Joe Gregorio", "", ""),
						},
						NewTestId("tag:example.org,2003:3.2397"),
						[]*Link{
							NewTestLink("http://example.org/2005/04/02/atom", "alternate", "text/html", "", "", ""),
							NewTestLink("http://example.org/audio/ph34r_my_podcast.mp3", "enclosure", "audio/mpeg", "", "", "1337"),
						},
						NewTestDate("2003-12-13T08:29:29-04:00"),
						NewTextConstruct(),
						NewTextConstruct(),
						NewTestTextConstruct("", "Atom draft-07 snapshot"),
						NewTestDate("2005-07-31T12:29:29Z"),
						NewSource(),
					),
				},
			), "en-us", "http://yo.com"),
		},
		{`
  <feed xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <link rel="self" type="application/atom+xml" href="http://example.org/feed.atom"/>
    <id>tag:example.org,2003:3</id>
    <author>
      <name>Mark Pilgrim</name>
    </author>
    <entry>
       <title>Atom-Powered Robots Run Amok</title>
       <author>
         <name>Jean</name>  
       </author>
       <link href="http://example.org/2003/12/13/atom03"/>
       <id>urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a</id>
       <updated>2003-12-13T18:30:02Z</updated>
       <summary>Some text.</summary>
    </entry>
  </feed>`,
			nil,
			NewTestFeed(
				[]*Person{
					NewTestPerson("Mark Pilgrim", "", ""),
				},
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				[]*Link{
					NewTestLink("http://example.org/feed.atom", "self", "application/atom+xml", "", "", ""),
				},
				NewLogo(),
				NewTextConstruct(),
				NewTextConstruct(),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
				[]*Entry{
					NewTestEntry(
						[]*Person{
							NewTestPerson("Jean", "", ""),
						},
						nil,
						NewContent(),
						nil,
						NewTestId("urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a"),
						[]*Link{
							NewTestLink("http://example.org/2003/12/13/atom03", "alternate", "", "", "", ""),
						},
						NewDate(),
						NewTextConstruct(),
						NewTestTextConstruct("", "Some text."),
						NewTestTextConstruct("", "Atom-Powered Robots Run Amok"),
						NewTestDate("2003-12-13T18:30:02Z"),
						NewSource(),
					),
				},
			),
		},
		{`
  <feed xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <author>
      <name>Mark Pilgrim</name>
    </author>
    <link rel="self" type="application/atom+xml" href="http://example.org/feed.atom"/>
    <id>tag:example.org,2003:3</id>
    <entry>
       <title>Atom-Powered Robots Run Amok</title>
       <link href="http://example.org/2003/12/13/atom03"/>
       <id>urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a</id>
       <updated>2003-12-13T18:30:02Z</updated>
       <summary>Some text.</summary>
    </entry>
  </feed>`,
			nil,
			NewTestFeed(
				[]*Person{
					NewTestPerson("Mark Pilgrim", "", ""),
				},
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				[]*Link{
					NewTestLink("http://example.org/feed.atom", "self", "application/atom+xml", "", "", ""),
				},
				NewLogo(),
				NewTextConstruct(),
				NewTextConstruct(),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
				[]*Entry{
					NewTestEntry(
						nil,
						nil,
						NewContent(),
						nil,
						NewTestId("urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a"),
						[]*Link{
							NewTestLink("http://example.org/2003/12/13/atom03", "alternate", "", "", "", ""),
						},
						NewDate(),
						NewTextConstruct(),
						NewTestTextConstruct("", "Some text."),
						NewTestTextConstruct("", "Atom-Powered Robots Run Amok"),
						NewTestDate("2003-12-13T18:30:02Z"),
						NewSource(),
					),
				},
			),
		},

		{`
  <feed xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <link rel="self" type="application/atom+xml" href="http://example.org/feed.atom"/>
    <id>tag:example.org,2003:3</id>
    <entry>
       <title>Atom-Powered Robots Run Amok</title>
       <link href="http://example.org/2003/12/13/atom03"/>
       <id>urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a</id>
       <updated>2003-12-13T18:30:02Z</updated>
       <summary>Some text.</summary>
    </entry>
  </feed>`,
			helper.NewError(MissingAuthor, ""),
			NewTestFeed(
				nil,
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				[]*Link{
					NewTestLink("http://example.org/feed.atom", "self", "application/atom+xml", "", "", ""),
				},
				NewLogo(),
				NewTextConstruct(),
				NewTextConstruct(),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
				[]*Entry{
					NewTestEntry(
						nil,
						nil,
						NewContent(),
						nil,
						NewTestId("urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a"),
						[]*Link{
							NewTestLink("http://example.org/2003/12/13/atom03", "alternate", "", "", "", ""),
						},
						NewDate(),
						NewTextConstruct(),
						NewTestTextConstruct("", "Some text."),
						NewTestTextConstruct("", "Atom-Powered Robots Run Amok"),
						NewTestDate("2003-12-13T18:30:02Z"),
						NewSource(),
					),
				},
			),
		},
		{`
  <feed xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <author>
      <name>Mark Pilgrim</name>
    </author>
    <link rel="self" type="application/atom+xml" href="http://example.org/feed.atom"/>
    <id>tag:example.org,2003:3</id>
    <generator uri="http://www.example.com/" version="1.0">Example Toolkit</generator>
    <generator uri="http://www.example2.com/" version="1.0">Example Toolkit</generator>
  </feed>`,
			helper.NewError(AttributeDuplicated, ""),
			NewTestFeed(
				[]*Person{
					NewTestPerson("Mark Pilgrim", "", ""),
				},
				nil,
				nil,
				NewTestGenerator("1.0", "http://www.example2.com/", "Example Toolkit"),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				[]*Link{
					NewTestLink("http://example.org/feed.atom", "self", "application/atom+xml", "", "", ""),
				},
				NewLogo(),
				NewTextConstruct(),
				NewTextConstruct(),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
				nil,
			),
		},
		{`
  <feed xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <author>
      <name>Mark Pilgrim</name>
    </author>
    <id>tag:example.org,2003:3</id>
    <generator uri="http://www.example2.com/" version="1.0">Example Toolkit</generator>
  </feed>`,
			helper.NewError(MissingSelfLink, ""),
			NewTestFeed(
				[]*Person{
					NewTestPerson("Mark Pilgrim", "", ""),
				},
				nil,
				nil,
				NewTestGenerator("1.0", "http://www.example2.com/", "Example Toolkit"),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				nil,
				NewLogo(),
				NewTextConstruct(),
				NewTextConstruct(),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
				nil,
			),
		},

		{`
  <feed xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <author>
      <name>Mark Pilgrim</name>
    </author>
    <link rel="self" type="application/atom+xml" href="http://example.org/feed.atom"/>
    <id>tag:example.org,2003:3</id>
    <icon>https://www.yo.com</icon>
    <icon>https://www.yo2.com</icon>
  </feed>`,
			helper.NewError(AttributeDuplicated, ""),
			NewTestFeed(
				[]*Person{
					NewTestPerson("Mark Pilgrim", "", ""),
				},
				nil,
				nil,
				NewGenerator(),
				NewTestIcon("https://www.yo2.com"),
				NewTestId("tag:example.org,2003:3"),
				[]*Link{
					NewTestLink("http://example.org/feed.atom", "self", "application/atom+xml", "", "", ""),
				},
				NewLogo(),
				NewTextConstruct(),
				NewTextConstruct(),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
				nil,
			),
		},
		{`
  <feed xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <author>
      <name>Mark Pilgrim</name>
    </author>
    <link rel="self" type="application/atom+xml" href="http://example.org/feed.atom"/>
  </feed>`,
			helper.NewError(MissingId, ""),
			NewTestFeed(
				[]*Person{
					NewTestPerson("Mark Pilgrim", "", ""),
				},
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewId(),
				[]*Link{
					NewTestLink("http://example.org/feed.atom", "self", "application/atom+xml", "", "", ""),
				},
				NewLogo(),
				NewTextConstruct(),
				NewTextConstruct(),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
				nil,
			),
		},
		{`
  <feed xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <author>
      <name>Mark Pilgrim</name>
    </author>
    <link rel="self" type="application/atom+xml" href="http://example.org/feed.atom"/>
    <id>tag:example.org,2003:3</id>
    <id>tag:example.org,2010:3</id>
  </feed>`,
			helper.NewError(IdDuplicated, ""),
			NewTestFeed(
				[]*Person{
					NewTestPerson("Mark Pilgrim", "", ""),
				},
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2010:3"),
				[]*Link{
					NewTestLink("http://example.org/feed.atom", "self", "application/atom+xml", "", "", ""),
				},
				NewLogo(),
				NewTextConstruct(),
				NewTextConstruct(),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
				nil,
			),
		},
		{`
  <feed xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <author>
      <name>Mark Pilgrim</name>
    </author>
    <link rel="self" type="application/atom+xml" href="http://example.org/feed.atom"/>
    <link type="application/atom+xml" hreflang="fr" href="http://example.org/content"/> 
    <link type="application/atom+xml" hreflang="fr" href="http://example.org/content"/>
   <id>tag:example.org,2003:3</id>
  </feed>`,
			helper.NewError(LinkAlternateDuplicated, ""),
			NewTestFeed(
				[]*Person{
					NewTestPerson("Mark Pilgrim", "", ""),
				},
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				[]*Link{
					NewTestLink("http://example.org/feed.atom", "self", "application/atom+xml", "", "", ""),
					NewTestLink("http://example.org/content", "alternate", "application/atom+xml", "fr", "", ""),
					NewTestLink("http://example.org/content", "alternate", "application/atom+xml", "fr", "", ""),
				},
				NewLogo(),
				NewTextConstruct(),
				NewTextConstruct(),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
				nil,
			),
		},
		{`
  <feed xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <author>
      <name>Mark Pilgrim</name>
    </author>
    <link rel="self" type="application/atom+xml" href="http://example.org/feed.atom"/>
    <link type="application/atom+xml" hreflang="en" href="http://example.org/content"/> 
    <link type="application/atom+xml" hreflang="fr" href="http://example.org/content"/>
   <id>tag:example.org,2003:3</id>
  </feed>`,
			nil,
			NewTestFeed(
				[]*Person{
					NewTestPerson("Mark Pilgrim", "", ""),
				},
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				[]*Link{
					NewTestLink("http://example.org/feed.atom", "self", "application/atom+xml", "", "", ""),
					NewTestLink("http://example.org/content", "alternate", "application/atom+xml", "en", "", ""),
					NewTestLink("http://example.org/content", "alternate", "application/atom+xml", "fr", "", ""),
				},
				NewLogo(),
				NewTextConstruct(),
				NewTextConstruct(),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
				nil,
			),
		},
		{`
  <feed xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <author>
      <name>Mark Pilgrim</name>
    </author>
    <link rel="self" type="application/atom+xml" href="http://example.org/feed.atom"/>
    <id>tag:example.org,2003:3</id>
    <logo>https://www.yo.com</logo>
    <logo>https://www.yo2.com</logo>
  </feed>`,
			helper.NewError(AttributeDuplicated, ""),
			NewTestFeed(
				[]*Person{
					NewTestPerson("Mark Pilgrim", "", ""),
				},
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				[]*Link{
					NewTestLink("http://example.org/feed.atom", "self", "application/atom+xml", "", "", ""),
				},
				NewTestLogo("https://www.yo2.com"),
				NewTextConstruct(),
				NewTextConstruct(),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
				nil,
			),
		},
		{`
  <feed xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <author>
      <name>Mark Pilgrim</name>
    </author>
    <id>tag:example.org,2003:3</id>
    <link rel="self" type="application/atom+xml"
     href="http://example.org/feed.atom"/>
    <rights>Copyright (c) 2003, Mark Pilgrim</rights>
    <rights>Copyright (c) 2010, Mark Pilgrim</rights>
  </feed>`,
			helper.NewError(AttributeDuplicated, ""),
			NewTestFeed(
				[]*Person{
					NewTestPerson("Mark Pilgrim", "", ""),
				},
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				[]*Link{
					NewTestLink("http://example.org/feed.atom", "self", "application/atom+xml", "", "", ""),
				},
				NewLogo(),
				NewTestTextConstruct("", "Copyright (c) 2010, Mark Pilgrim"),
				NewTextConstruct(),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
				nil,
			),
		},
		{`
  <feed xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <id>tag:example.org,2003:3</id>
    <author>
      <name>Mark Pilgrim</name>
    </author>
    <link rel="self" type="application/atom+xml"
     href="http://example.org/feed.atom"/>
    <rights>Copyright (c) 2003, Mark Pilgrim</rights>
    <rights>Copyright (c) 2010, Mark Pilgrim</rights>
  </feed>`,
			helper.NewError(AttributeDuplicated, ""),
			NewTestFeed(
				[]*Person{
					NewTestPerson("Mark Pilgrim", "", ""),
				},
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				[]*Link{
					NewTestLink("http://example.org/feed.atom", "self", "application/atom+xml", "", "", ""),
				},
				NewLogo(),
				NewTestTextConstruct("", "Copyright (c) 2010, Mark Pilgrim"),
				NewTextConstruct(),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
				nil,
			),
		},
		{`
  <feed xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <id>tag:example.org,2003:3</id>
    <author>
      <name>Mark Pilgrim</name>
    </author>
    <link rel="self" type="application/atom+xml"
     href="http://example.org/feed.atom"/>
    <rights>Copyright (c) 2003, Mark Pilgrim</rights>
    <rights>Copyright (c) 2010, Mark Pilgrim</rights>
  </feed>`,
			helper.NewError(AttributeDuplicated, ""),
			NewTestFeed(
				[]*Person{
					NewTestPerson("Mark Pilgrim", "", ""),
				},
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				[]*Link{
					NewTestLink("http://example.org/feed.atom", "self", "application/atom+xml", "", "", ""),
				},
				NewLogo(),
				NewTestTextConstruct("", "Copyright (c) 2010, Mark Pilgrim"),
				NewTextConstruct(),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
				nil,
			),
		},
		{`
  <feed xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <id>tag:example.org,2003:3</id>
    <author>
      <name>Mark Pilgrim</name>
    </author>
    <link rel="self" type="application/atom+xml"
     href="http://example.org/feed.atom"/>
    <rights>Copyright (c) 2003, Mark Pilgrim</rights>
    <rights>Copyright (c) 2010, Mark Pilgrim</rights>
  </feed>`,
			helper.NewError(AttributeDuplicated, ""),
			NewTestFeed(
				[]*Person{
					NewTestPerson("Mark Pilgrim", "", ""),
				},
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				[]*Link{
					NewTestLink("http://example.org/feed.atom", "self", "application/atom+xml", "", "", ""),
				},
				NewLogo(),
				NewTestTextConstruct("", "Copyright (c) 2010, Mark Pilgrim"),
				NewTextConstruct(),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
				nil,
			),
		},
		{`
  <feed xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <id>tag:example.org,2003:3</id>
    <author>
      <name>Mark Pilgrim</name>
    </author>
    <link rel="self" type="application/atom+xml"
     href="http://example.org/feed.atom"/>
    <rights>Copyright (c) 2003, Mark Pilgrim</rights>
    <rights>Copyright (c) 2010, Mark Pilgrim</rights>
  </feed>`,
			helper.NewError(AttributeDuplicated, ""),
			NewTestFeed(
				[]*Person{
					NewTestPerson("Mark Pilgrim", "", ""),
				},
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				[]*Link{
					NewTestLink("http://example.org/feed.atom", "self", "application/atom+xml", "", "", ""),
				},
				NewLogo(),
				NewTestTextConstruct("", "Copyright (c) 2010, Mark Pilgrim"),
				NewTextConstruct(),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
				nil,
			),
		},
		{`
  <feed xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <subtitle type="html">A &lt;em&gt;lot&lt;/em&gt; of effort went into making this effortless</subtitle>
    <subtitle type="html">A &lt;em&gt;lot&lt;/em&gt; of effort went into making this effortless again</subtitle>
    <id>tag:example.org,2003:3</id>
    <author>
      <name>Mark Pilgrim</name>
    </author>
    <link rel="self" type="application/atom+xml"
     href="http://example.org/feed.atom"/>
    <rights>Copyright (c) 2003, Mark Pilgrim</rights>
  </feed>`,
			helper.NewError(AttributeDuplicated, ""),
			NewTestFeed(
				[]*Person{
					NewTestPerson("Mark Pilgrim", "", ""),
				},
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				[]*Link{
					NewTestLink("http://example.org/feed.atom", "self", "application/atom+xml", "", "", ""),
				},
				NewLogo(),
				NewTestTextConstruct("", "Copyright (c) 2003, Mark Pilgrim"),
				NewTestTextConstruct("", "A <em>lot</em> of effort went into making this effortless again"),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
				nil,
			),
		},
		{`
  <feed xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <title type="text">dive into mark again</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <author>
      <name>Mark Pilgrim</name>
    </author>
    <subtitle type="html">A &lt;em&gt;lot&lt;/em&gt; of effort went into making this effortless</subtitle>
    <id>tag:example.org,2003:3</id>
    <link rel="self" type="application/atom+xml"
     href="http://example.org/feed.atom"/>
    <rights>Copyright (c) 2003, Mark Pilgrim</rights>
  </feed>`,
			helper.NewError(TitleDuplicated, ""),
			NewTestFeed(
				[]*Person{
					NewTestPerson("Mark Pilgrim", "", ""),
				},
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				[]*Link{
					NewTestLink("http://example.org/feed.atom", "self", "application/atom+xml", "", "", ""),
				},
				NewLogo(),
				NewTestTextConstruct("", "Copyright (c) 2003, Mark Pilgrim"),
				NewTestTextConstruct("", "A <em>lot</em> of effort went into making this effortless"),
				NewTestTextConstruct("", "dive into mark again"),
				NewTestDate("2005-07-31T12:29:29Z"),
				nil,
			),
		},
		{`
  <feed xmlns="http://www.w3.org/2005/Atom">
    <updated>2005-07-31T12:29:29Z</updated>
    <subtitle type="html">A &lt;em&gt;lot&lt;/em&gt; of effort went into making this effortless</subtitle>
    <id>tag:example.org,2003:3</id>
    <author>
      <name>Mark Pilgrim</name>
    </author>
    <link rel="self" type="application/atom+xml"
     href="http://example.org/feed.atom"/>
    <rights>Copyright (c) 2003, Mark Pilgrim</rights>
  </feed>`,
			helper.NewError(MissingTitle, ""),
			NewTestFeed(
				[]*Person{
					NewTestPerson("Mark Pilgrim", "", ""),
				},
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				[]*Link{
					NewTestLink("http://example.org/feed.atom", "self", "application/atom+xml", "", "", ""),
				},
				NewLogo(),
				NewTestTextConstruct("", "Copyright (c) 2003, Mark Pilgrim"),
				NewTestTextConstruct("", "A <em>lot</em> of effort went into making this effortless"),
				NewTextConstruct(),
				NewTestDate("2005-07-31T12:29:29Z"),
				nil,
			),
		},
		{`
  <feed xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <subtitle type="html">A &lt;em&gt;lot&lt;/em&gt; of effort went into making this effortless</subtitle>
    <id>tag:example.org,2003:3</id>
    <author>
      <name>Mark Pilgrim</name>
    </author>
    <link rel="self" type="application/atom+xml"
     href="http://example.org/feed.atom"/>
    <rights>Copyright (c) 2003, Mark Pilgrim</rights>
  </feed>`,
			helper.NewError(MissingDate, ""),
			NewTestFeed(
				[]*Person{
					NewTestPerson("Mark Pilgrim", "", ""),
				},
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				[]*Link{
					NewTestLink("http://example.org/feed.atom", "self", "application/atom+xml", "", "", ""),
				},
				NewLogo(),
				NewTestTextConstruct("", "Copyright (c) 2003, Mark Pilgrim"),
				NewTestTextConstruct("", "A <em>lot</em> of effort went into making this effortless"),
				NewTestTextConstruct("", "dive into mark"),
				NewDate(),
				nil,
			),
		},
		{`
  <feed xmlns="http://www.w3.org/2005/Atom">
    <updated>2005-07-31T12:29:29Z</updated>
    <updated>2010-07-31T12:29:29Z</updated>
    <title type="text">dive into mark</title>
    <subtitle type="html">A &lt;em&gt;lot&lt;/em&gt; of effort went into making this effortless</subtitle>
    <id>tag:example.org,2003:3</id>
    <author>
      <name>Mark Pilgrim</name>
    </author>
    <link rel="self" type="application/atom+xml"
     href="http://example.org/feed.atom"/>
    <rights>Copyright (c) 2003, Mark Pilgrim</rights>
  </feed>`,
			helper.NewError(AttributeDuplicated, ""),
			NewTestFeed(
				[]*Person{
					NewTestPerson("Mark Pilgrim", "", ""),
				},
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				[]*Link{
					NewTestLink("http://example.org/feed.atom", "self", "application/atom+xml", "", "", ""),
				},
				NewLogo(),
				NewTestTextConstruct("", "Copyright (c) 2003, Mark Pilgrim"),
				NewTestTextConstruct("", "A <em>lot</em> of effort went into making this effortless"),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2010-07-31T12:29:29Z"),
				nil,
			),
		},
		{`
  <feed xmlns="http://www.w3.org/2005/Atom">
    <updated>2005-07-31T12:29:29Z</updated>
    <title type="text">dive into mark</title>
    <subtitle type="html">A &lt;em&gt;lot&lt;/em&gt; of effort went into making this effortless</subtitle>
    <id>tag:example.org,2003:3</id>
    <link rel="self" type="application/atom+xml"
     href="http://example.org/feed.atom"/>
    <rights>Copyright (c) 2003, Mark Pilgrim</rights>
  </feed>`,
			helper.NewError(MissingAuthor, ""),
			NewTestFeed(
				nil,
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				[]*Link{
					NewTestLink("http://example.org/feed.atom", "self", "application/atom+xml", "", "", ""),
				},
				NewLogo(),
				NewTestTextConstruct("", "Copyright (c) 2003, Mark Pilgrim"),
				NewTestTextConstruct("", "A <em>lot</em> of effort went into making this effortless"),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
				nil,
			),
		},
		{`
  <feed xmlns="http://www.w3.org/2005/Atom">
    <title type="text">dive into mark</title>
    <updated>2005-07-31T12:29:29Z</updated>
    <link rel="self" type="application/atom+xml" href="http://example.org/feed.atom"/>
    <id>tag:example.org,2003:3</id>
    <author>
      <name>Mark Pilgrim</name>
    </author>
    <entry>
       <title>Atom-Powered Robots Run Amok</title>
       <link href="http://example.org/2003/12/13/atom03"/>
       <id>urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a</id>
       <updated>2003-12-13T18:30:02Z</updated>
       <summary>Some text.</summary>
    </entry>
    <entry>
       <title>Atom-Powered Robots Run Amok</title>
       <link href="http://example.org/2003/12/13/atom03"/>
       <id>urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a</id>
       <updated>2003-12-13T18:30:02Z</updated>
       <summary>Some text.</summary>
    </entry>
  </feed>`,
			helper.NewError(EntryWithIdAndDateDuplicated, ""),
			NewTestFeed(
				[]*Person{
					NewTestPerson("Mark Pilgrim", "", ""),
				},
				nil,
				nil,
				NewGenerator(),
				NewIcon(),
				NewTestId("tag:example.org,2003:3"),
				[]*Link{
					NewTestLink("http://example.org/feed.atom", "self", "application/atom+xml", "", "", ""),
				},
				NewLogo(),
				NewTextConstruct(),
				NewTextConstruct(),
				NewTestTextConstruct("", "dive into mark"),
				NewTestDate("2005-07-31T12:29:29Z"),
				[]*Entry{
					NewTestEntry(
						nil,
						nil,
						NewContent(),
						nil,
						NewTestId("urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a"),
						[]*Link{
							NewTestLink("http://example.org/2003/12/13/atom03", "alternate", "", "", "", ""),
						},
						NewDate(),
						NewTextConstruct(),
						NewTestTextConstruct("", "Some text."),
						NewTestTextConstruct("", "Atom-Powered Robots Run Amok"),
						NewTestDate("2003-12-13T18:30:02Z"),
						NewSource(),
					),
					NewTestEntry(
						nil,
						nil,
						NewContent(),
						nil,
						NewTestId("urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a"),
						[]*Link{
							NewTestLink("http://example.org/2003/12/13/atom03", "alternate", "", "", "", ""),
						},
						NewDate(),
						NewTextConstruct(),
						NewTestTextConstruct("", "Some text."),
						NewTestTextConstruct("", "Atom-Powered Robots Run Amok"),
						NewTestDate("2003-12-13T18:30:02Z"),
						NewSource(),
					),
				},
			),
		},
	}

	nbErrors := 0
	len := len(testdata)
	for _, testfeed := range testdata {
		testcase := _TestFeedToTestVisitor(testfeed)

		if err := testcase.CheckTestCase(); err != nil {
			t.Errorf("FAIL\n%s\nXML:\n %s\n", err, testcase.XML)
			nbErrors++
		}
	}

	t.Logf("PASS RATIO = %v/%v\n", len-nbErrors, len)
}
