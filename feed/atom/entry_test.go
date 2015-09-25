package atom

import (
	"fmt"
	"testing"

	xmlutils "github.com/jloup/xml/utils"
)

func NewTestEntry(
	authors []*Person,
	cat []*Category,
	content *Content,
	contributors []*Person,
	id *Id,
	links []*Link,
	pub *Date,
	rights, sum, title *TextConstruct,
	updated *Date,
	source *Source,
) *Entry {

	e := NewEntry()

	e.Authors = authors
	e.Categories = cat
	e.Content = content
	e.Contributors = contributors
	e.Id = id
	e.Links = links
	e.Published = pub
	e.Rights = rights
	e.Summary = sum
	e.Title = title
	e.Updated = updated
	e.Source = source

	return e
}

func EntryWithBaseLang(e *Entry, lang, base string) *Entry {
	e.Lang.Value = lang
	e.Base.Value = base

	return e
}

type testEntry struct {
	XML           string
	ExpectedError xmlutils.ParserError
	ExpectedEntry *Entry
}

func testEntryValidator(actual xmlutils.Visitor, expected xmlutils.Visitor) error {
	e1 := actual.(*Entry)
	e2 := expected.(*Entry)

	if len(e1.Authors) != len(e2.Authors) {
		return fmt.Errorf("Entry does not contain the right count of Authors %v (expected) vs %v", len(e2.Authors), len(e1.Authors))
	}

	for i, _ := range e1.Authors {
		if err := testPersonValidator(e1.Authors[i], e2.Authors[i]); err != nil {
			return err
		}

	}

	if len(e1.Categories) != len(e2.Categories) {
		return fmt.Errorf("Entry does not contain the right count of Categories %v (expected) vs %v", len(e2.Categories), len(e1.Categories))
	}

	for i, _ := range e1.Categories {
		if err := testCategoryValidator(e1.Categories[i], e2.Categories[i]); err != nil {
			return err
		}

	}

	if err := testContentValidator(e1.Content, e2.Content); err != nil {
		return err
	}

	if len(e1.Contributors) != len(e2.Contributors) {
		return fmt.Errorf("Entry does not contain the right count of Contributors %v (expected) vs %v", len(e2.Contributors), len(e1.Contributors))
	}

	for i, _ := range e1.Contributors {
		if err := testPersonValidator(e1.Contributors[i], e2.Contributors[i]); err != nil {
			return err
		}

	}

	if err := testIdValidator(e1.Id, e2.Id); err != nil {
		return err
	}

	if len(e1.Links) != len(e2.Links) {
		return fmt.Errorf("Entry does not contain the right count of Links %v (expected) vs %v", len(e2.Links), len(e1.Links))
	}

	for i, _ := range e1.Links {
		if err := testLinkValidator(e1.Links[i], e2.Links[i]); err != nil {
			return err
		}

	}

	if err := testDateValidator(e1.Published, e2.Published); err != nil {
		return err
	}

	if err := testTextConstructValidator(e1.Rights, e2.Rights); err != nil {
		return err
	}

	if err := testTextConstructValidator(e1.Summary, e2.Summary); err != nil {
		return err
	}

	if err := testTextConstructValidator(e1.Title, e2.Title); err != nil {
		return err
	}

	if err := testDateValidator(e1.Updated, e2.Updated); err != nil {
		return err
	}

	if err := testSourceValidator(e1.Source, e2.Source); err != nil {
		return err
	}

	if err := ValidateBaseLang("Entry", e1.Base.Value, e1.Lang.Value, e2.Base.Value, e2.Lang.Value); err != nil {
		return err
	}

	return nil
}

func testEntryConstructor() xmlutils.Visitor {
	return NewEntry()
}

func _TestEntryToTestVisitor(t testEntry) xmlutils.TestVisitor {
	testVisitor := xmlutils.TestVisitor{
		XML:                t.XML,
		ExpectedError:      nil,
		ExpectedVisitor:    t.ExpectedEntry,
		VisitorConstructor: testEntryConstructor,
		Validator:          testEntryValidator,
	}

	if t.ExpectedError != nil {
		testVisitor.ExpectedError = t.ExpectedError
	}

	return testVisitor
}

func TestEntryBasic(t *testing.T) {

	var testdata = []testEntry{
		{`
     <entry xml:lang="en-us" xml:base="http://yo.com" xmlns:thr="http://purl.org/syndication/thread/1.0">
       <title>Atom-Powered Robots Run Amok</title>
       <author>
         <name>Jean</name>  
       </author>
       <link href="http://example.org/2003/12/13/atom03"/>
       <id>urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a</id>
       <updated>2003-12-13T18:30:02Z</updated>
       <category scheme='https://www.yo.com' term='music' />
       <category scheme='https://www.yo.com' term='house' />
       <summary>Some text.</summary>
       <thr:in-reply-to
         ref="tag:example.org,2005:1"
         type="application/xhtml+xml"
         href="http://www.example.org/entries/1"/>
     </entry>`,
			nil,
			EntryWithBaseLang(NewTestEntry(
				[]*Person{
					NewTestPerson("Jean", "", ""),
				},
				[]*Category{
					NewTestCategory("https://www.yo.com", "music", ""),
					NewTestCategory("https://www.yo.com", "house", ""),
				},
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
			), "en-us", "http://yo.com"),
		},
		{`
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
`,
			nil,
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
		{`
     <entry>
       <author>
         <name>Jean</name>  
       </author>
       <link href="http://example.org/2003/12/13/atom03"/>
       <id>urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a</id>
       <updated>2003-12-13T18:30:02Z</updated>
       <summary>Some text.</summary>
     </entry>
`,
			xmlutils.NewError(MissingTitle, ""),
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
				NewTextConstruct(),
				NewTestDate("2003-12-13T18:30:02Z"),
				NewSource(),
			),
		},
		{`
     <entry>
      <link href="http://example.com" />
      <title>Atom draft-07 snapshot</title>
       <author>
         <name>Jean</name>  
       </author>
      <id>tag:example.org,2003:3.2397</id>
      <updated>2005-07-31T12:29:29Z</updated>
      <summary>Some text.</summary>
     </entry>
`,
			nil,
			NewTestEntry(
				[]*Person{
					NewTestPerson("Jean", "", ""),
				},
				nil,
				NewContent(),
				nil,
				NewTestId("tag:example.org,2003:3.2397"),
				[]*Link{
					NewTestLink("http://example.com", "alternate", "", "", "", ""),
				},
				NewDate(),
				NewTextConstruct(),
				NewTestTextConstruct("", "Some text."),
				NewTestTextConstruct("", "Atom draft-07 snapshot"),
				NewTestDate("2005-07-31T12:29:29Z"),
				NewSource(),
			),
		},
		{`
     <entry>
      <link href="http://example.com" type="application/xhtml+xml"/>
      <title>Atom draft-07 snapshot</title>
       <author>
         <name>Jean</name>  
       </author>
      <id>tag:example.org,2003:3.2397</id>
      <updated>2005-07-31T12:29:29Z</updated>
      <summary>Some text.</summary>
      <content>FIRST CONTENT</content>
      <content>SECOND CONTENT</content>
     </entry>
`,
			xmlutils.NewError(AttributeDuplicated, ""),
			NewTestEntry(
				[]*Person{
					NewTestPerson("Jean", "", ""),
				},
				nil,
				NewTestContent("text", "", "SECOND CONTENT", "", ""),
				nil,
				NewTestId("tag:example.org,2003:3.2397"),
				[]*Link{
					NewTestLink("http://example.com", "alternate", "application/xhtml+xml", "", "", ""),
				},
				NewDate(),
				NewTextConstruct(),
				NewTestTextConstruct("", "Some text."),
				NewTestTextConstruct("", "Atom draft-07 snapshot"),
				NewTestDate("2005-07-31T12:29:29Z"),
				NewSource(),
			),
		},
		{`
     <entry>
      <link href="http://example.com" />
       <author>
         <name>Jean</name>  
       </author>
      <title>Atom draft-07 snapshot</title>
      <summary>Some text.</summary>
      <updated>2005-07-31T12:29:29Z</updated>
     </entry>
`,
			xmlutils.NewError(MissingId, ""),
			NewTestEntry(
				[]*Person{
					NewTestPerson("Jean", "", ""),
				},
				nil,
				NewContent(),
				nil,
				NewId(),
				[]*Link{
					NewTestLink("http://example.com", "alternate", "", "", "", ""),
				},
				NewDate(),
				NewTextConstruct(),
				NewTestTextConstruct("", "Some text."),
				NewTestTextConstruct("", "Atom draft-07 snapshot"),
				NewTestDate("2005-07-31T12:29:29Z"),
				NewSource(),
			),
		},
		{`
     <entry>
      <link href="http://example.com" />
       <author>
         <name>Jean</name>  
       </author>
      <title>Atom draft-07 snapshot</title>
      <id>tag:example.org,2003:3.2397</id>
      <id>tag:example.org,2003:3.2397</id>
      <updated>2005-07-31T12:29:29Z</updated>
      <summary>Some text.</summary>
     </entry>
`,
			xmlutils.NewError(IdDuplicated, ""),
			NewTestEntry(
				[]*Person{
					NewTestPerson("Jean", "", ""),
				},
				nil,
				NewContent(),
				nil,
				NewTestId("tag:example.org,2003:3.2397"),
				[]*Link{
					NewTestLink("http://example.com", "alternate", "", "", "", ""),
				},
				NewDate(),
				NewTextConstruct(),
				NewTestTextConstruct("", "Some text."),
				NewTestTextConstruct("", "Atom draft-07 snapshot"),
				NewTestDate("2005-07-31T12:29:29Z"),
				NewSource(),
			),
		},
		{`
     <entry>
      <link href="http://example.com" type="application/xhtml+xml"/>
      <link href="http://example.com" type="application/xhtml+xml"/>
      <title>Atom draft-07 snapshot</title>
       <author>
         <name>Jean</name>  
       </author>
      <id>tag:example.org,2003:3.2397</id>
      <updated>2005-07-31T12:29:29Z</updated>
      <summary>Some text.</summary>
     </entry>
`,
			xmlutils.NewError(LinkAlternateDuplicated, ""),
			NewTestEntry(
				[]*Person{
					NewTestPerson("Jean", "", ""),
				},
				nil,
				NewContent(),
				nil,
				NewTestId("tag:example.org,2003:3.2397"),
				[]*Link{
					NewTestLink("http://example.com", "alternate", "application/xhtml+xml", "", "", ""),
					NewTestLink("http://example.com", "alternate", "application/xhtml+xml", "", "", ""),
				},
				NewDate(),
				NewTextConstruct(),
				NewTestTextConstruct("", "Some text."),
				NewTestTextConstruct("", "Atom draft-07 snapshot"),
				NewTestDate("2005-07-31T12:29:29Z"),
				NewSource(),
			),
		},
		{`
     <entry>
      <link href="http://example.com" hreflang="en" type="application/xhtml+xml"/>
      <link href="http://example.com" hreflang="en" type="application/xhtml+xml"/>
      <title>Atom draft-07 snapshot</title>
       <author>
         <name>Jean</name>  
       </author>
      <id>tag:example.org,2003:3.2397</id>
      <updated>2005-07-31T12:29:29Z</updated>
      <summary>Some text.</summary>
     </entry>
`,
			xmlutils.NewError(LinkAlternateDuplicated, ""),
			NewTestEntry(
				[]*Person{
					NewTestPerson("Jean", "", ""),
				},
				nil,
				NewContent(),
				nil,
				NewTestId("tag:example.org,2003:3.2397"),
				[]*Link{
					NewTestLink("http://example.com", "alternate", "application/xhtml+xml", "en", "", ""),
					NewTestLink("http://example.com", "alternate", "application/xhtml+xml", "en", "", ""),
				},
				NewDate(),
				NewTextConstruct(),
				NewTestTextConstruct("", "Some text."),
				NewTestTextConstruct("", "Atom draft-07 snapshot"),
				NewTestDate("2005-07-31T12:29:29Z"),
				NewSource(),
			),
		},
		{`
     <entry>
      <link href="http://example.com" hreflang="en" type="application/xhtml+xml"/>
      <link href="http://example.com" hreflang="fr" type="application/xhtml+xml"/>
      <title>Atom draft-07 snapshot</title>
       <author>
         <name>Jean</name>  
       </author>
      <id>tag:example.org,2003:3.2397</id>
      <updated>2005-07-31T12:29:29Z</updated>
      <summary>Some text.</summary>
     </entry>
`,
			nil,
			NewTestEntry(
				[]*Person{
					NewTestPerson("Jean", "", ""),
				},
				nil,
				NewContent(),
				nil,
				NewTestId("tag:example.org,2003:3.2397"),
				[]*Link{
					NewTestLink("http://example.com", "alternate", "application/xhtml+xml", "en", "", ""),
					NewTestLink("http://example.com", "alternate", "application/xhtml+xml", "fr", "", ""),
				},
				NewDate(),
				NewTextConstruct(),
				NewTestTextConstruct("", "Some text."),
				NewTestTextConstruct("", "Atom draft-07 snapshot"),
				NewTestDate("2005-07-31T12:29:29Z"),
				NewSource(),
			),
		},
		{`
     <entry>
      <link href="http://example.com" hreflang="en" type="application/xhtml+xml"/>
      <title>Atom draft-07 snapshot</title>
       <author>
         <name>Jean</name>  
       </author>
      <id>tag:example.org,2003:3.2397</id>
      <updated>2005-07-31T12:29:29Z</updated>
   </entry>
`,
			xmlutils.NewError(MissingSummary, ""),
			NewTestEntry(
				[]*Person{
					NewTestPerson("Jean", "", ""),
				},
				nil,
				NewContent(),
				nil,
				NewTestId("tag:example.org,2003:3.2397"),
				[]*Link{
					NewTestLink("http://example.com", "alternate", "application/xhtml+xml", "en", "", ""),
				},
				NewDate(),
				NewTextConstruct(),
				NewTextConstruct(),
				NewTestTextConstruct("", "Atom draft-07 snapshot"),
				NewTestDate("2005-07-31T12:29:29Z"),
				NewSource(),
			),
		},
		{`
     <entry>
      <link href="http://example.com" hreflang="en" type="application/xhtml+xml"/>
      <title>Atom draft-07 snapshot</title>
       <author>
         <name>Jean</name>  
       </author>
      <id>tag:example.org,2003:3.2397</id>
      <updated>2005-07-31T12:29:29Z</updated>
      <content src="http://g.com" type="application/xml"></content>  
   </entry>
`,
			xmlutils.NewError(MissingSummary, ""),
			NewTestEntry(
				[]*Person{
					NewTestPerson("Jean", "", ""),
				},
				nil,
				NewTestContent("application/xml", "", "", "", "http://g.com"),
				nil,
				NewTestId("tag:example.org,2003:3.2397"),
				[]*Link{
					NewTestLink("http://example.com", "alternate", "application/xhtml+xml", "en", "", ""),
				},
				NewDate(),
				NewTextConstruct(),
				NewTextConstruct(),
				NewTestTextConstruct("", "Atom draft-07 snapshot"),
				NewTestDate("2005-07-31T12:29:29Z"),
				NewSource(),
			),
		},
		{`
     <entry>
      <link href="http://example.com" hreflang="en" type="application/xhtml+xml"/>
      <title>Atom draft-07 snapshot</title>
       <author>
         <name>Jean</name>  
       </author>
      <id>tag:example.org,2003:3.2397</id>
      <updated>2005-07-31T12:29:29Z</updated>
      <content type="image/png">iVBORw0KGgoAAAANSUhEUgAAAB8AAAAqCAYAAABLGYAnAAAABmJLR0QA/wD/AP+gvaeTAAAACXBIWXMAAAsTAAALEwEAmpwYAAAAB3RJTUUH1QwCBCUlRSCuygAAAetJREFUWMPt1j1IVmEUB/Dfo1mvBg1iSgjhZNGnFIhTDUXQ25x9LE2NbQ1NLe1tbbo1tjQEDk2BUzSITUapUejyvkQthXKfhvvcEs3gvl55Fd4/HLgfz//8zzn3Ps85tBGhBc4oLuIYuvAV77CwW0F24x7mEbex+bSmu0rhEcz9R3SzzSXOjnEWjRLChTUSt2X0Y7EF4cIWk4+WMLUD4cKmtPhHr1cgvp58/RNd2zy/VdFf2518lcJsBVkXNls285GKt2qpE+4nDlUk/gu1Mpk3Ksy8UbbsSxWKL5UVn6lQfGZP7vM9ecK1/Wxva1drez/f1UlmX8xwHXTQQQcd7Ctkl8jGiMeJk8SeDe9uEI+Q1bfyYg/xNrGf7CaxRhzP77esPUy8soFXLwbIU4QaPuM6YY04SDxAOI2MMJGIw8Q0z4c14lg+k4dr6MEAoUkcIvYlzgAO4nKK5CgmknjIUou8mvflrJ6uH+Vt+k/09UR88rc64Xnq4Su4gzXiKMbxgHgBTzGcfEyirxivitH5LeF1yvIkXuLZptKdwQ98Q28Sf58ymsbd1NdPJMKHlPFCWgfnsIxmEo/NvHRxCB/xgng/ZfkFg1glTBPP4w3h+4agHhOW8TAvuVcpuBV8yrmxl7iKVKm4mn/WDtqA3yOQKuHaSApTAAAAAElFTkSuQmCC</content>
   </entry>
`,
			xmlutils.NewError(MissingSummary, ""),
			NewTestEntry(
				[]*Person{
					NewTestPerson("Jean", "", ""),
				},
				nil,
				NewTestContent("image/png", "", "", "iVBORw0KGgoAAAANSUhEUgAAAB8AAAAqCAYAAABLGYAnAAAABmJLR0QA/wD/AP+gvaeTAAAACXBIWXMAAAsTAAALEwEAmpwYAAAAB3RJTUUH1QwCBCUlRSCuygAAAetJREFUWMPt1j1IVmEUB/Dfo1mvBg1iSgjhZNGnFIhTDUXQ25x9LE2NbQ1NLe1tbbo1tjQEDk2BUzSITUapUejyvkQthXKfhvvcEs3gvl55Fd4/HLgfz//8zzn3Ps85tBGhBc4oLuIYuvAV77CwW0F24x7mEbex+bSmu0rhEcz9R3SzzSXOjnEWjRLChTUSt2X0Y7EF4cIWk4+WMLUD4cKmtPhHr1cgvp58/RNd2zy/VdFf2518lcJsBVkXNls285GKt2qpE+4nDlUk/gu1Mpk3Ksy8UbbsSxWKL5UVn6lQfGZP7vM9ecK1/Wxva1drez/f1UlmX8xwHXTQQQcd7Ctkl8jGiMeJk8SeDe9uEI+Q1bfyYg/xNrGf7CaxRhzP77esPUy8soFXLwbIU4QaPuM6YY04SDxAOI2MMJGIw8Q0z4c14lg+k4dr6MEAoUkcIvYlzgAO4nKK5CgmknjIUou8mvflrJ6uH+Vt+k/09UR88rc64Xnq4Su4gzXiKMbxgHgBTzGcfEyirxivitH5LeF1yvIkXuLZptKdwQ98Q28Sf58ymsbd1NdPJMKHlPFCWgfnsIxmEo/NvHRxCB/xgng/ZfkFg1glTBPP4w3h+4agHhOW8TAvuVcpuBV8yrmxl7iKVKm4mn/WDtqA3yOQKuHaSApTAAAAAElFTkSuQmCC", ""),
				nil,
				NewTestId("tag:example.org,2003:3.2397"),
				[]*Link{
					NewTestLink("http://example.com", "alternate", "application/xhtml+xml", "en", "", ""),
				},
				NewDate(),
				NewTextConstruct(),
				NewTextConstruct(),
				NewTestTextConstruct("", "Atom draft-07 snapshot"),
				NewTestDate("2005-07-31T12:29:29Z"),
				NewSource(),
			),
		},
		{`
     <entry>
       <title>Atom-Powered Robots Run Amok</title>
       <link href="http://example.org/2003/12/13/atom03"/>
       <id>urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a</id>
       <updated>2003-12-13T18:30:02Z</updated>
       <category scheme='https://www.yo.com' term='music' />
       <category scheme='https://www.yo.com' term='house' />
       <summary>Some text.</summary>
     </entry>`,
			xmlutils.NewError(MissingAuthor, ""),
			NewTestEntry(
				nil,
				[]*Category{
					NewTestCategory("https://www.yo.com", "music", ""),
					NewTestCategory("https://www.yo.com", "house", ""),
				},
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
		{`
     <entry>
       <title>Atom-Powered Robots Run Amok</title>
       <link href="http://example.org/2003/12/13/atom03"/>
       <id>urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a</id>
       <updated>2003-12-13T18:30:02Z</updated>
       <category scheme='https://www.yo.com' term='music' />
       <category scheme='https://www.yo.com' term='house' />
       <summary>Some text.</summary>
       <source xmlns="http://www.w3.org/2005/Atom">
         <title type="text">dive into mark</title>
         <updated>2005-07-31T12:29:29Z</updated>
         <id>tag:example.org,2003:3</id>  
         <author>
           <name>jean</name>
         </author>
       </source>       
     </entry>`,
			nil,
			NewTestEntry(
				nil,
				[]*Category{
					NewTestCategory("https://www.yo.com", "music", ""),
					NewTestCategory("https://www.yo.com", "house", ""),
				},
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
				NewTestSource(
					[]*Person{
						NewTestPerson("jean", "", ""),
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
			),
		},
		{`
     <entry>
       <title>Atom-Powered Robots Run Amok</title>
       <link href="http://example.org/2003/12/13/atom03"/>
       <id>urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a</id>
       <updated>2003-12-13T18:30:02Z</updated>
       <category scheme='https://www.yo.com' term='music' />
       <category scheme='https://www.yo.com' term='house' />
       <summary>Some text.</summary>
       <source xmlns="http://www.w3.org/2005/Atom">
         <title type="text">dive into mark</title>
         <updated>2005-07-31T12:29:29Z</updated>
         <id>tag:example.org,2003:3</id>  
       </source>       
     </entry>`,
			xmlutils.NewError(MissingAuthor, ""),
			NewTestEntry(
				nil,
				[]*Category{
					NewTestCategory("https://www.yo.com", "music", ""),
					NewTestCategory("https://www.yo.com", "house", ""),
				},
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
			),
		},
	}

	nbErrors := 0
	len := len(testdata)
	for _, testentry := range testdata {
		testcase := _TestEntryToTestVisitor(testentry)

		if err := testcase.CheckTestCase(); err != nil {
			t.Errorf("FAIL\n%s\nXML:\n %s\n", err, testcase.XML)
			nbErrors++
		}
	}

	t.Logf("PASS RATIO = %v/%v\n", len-nbErrors, len)
}
