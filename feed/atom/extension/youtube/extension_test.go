package youtube

import (
	"fmt"
	"testing"

	"github.com/jloup/xml/feed/atom"
	"github.com/jloup/xml/feed/extension"
	xmlutils "github.com/jloup/xml/utils"
)

func NewTestYoutubeFeed(channelId string, entries []*atom.Entry) *atom.Feed {
	f := atom.NewFeed()

	f.Entries = entries
	if channelId != "" {
		c := atom.NewBasicElement(nil)
		c.Content.Value = channelId
		f.Extension.Store.Add(_channelId, c)
	}

	return f
}

func NewTestYoutubeEntry(videoId, channelId string) *atom.Entry {
	e := atom.NewEntry()

	if videoId != "" {
		v := atom.NewBasicElement(nil)
		v.Content.Value = videoId
		e.Extension.Store.Add(_videoId, v)
	}

	if channelId != "" {
		c := atom.NewBasicElement(nil)
		c.Content.Value = channelId
		e.Extension.Store.Add(_channelId, c)
	}

	return e
}

type testYoutubeFeed struct {
	XML                 string
	ExpectedError       xmlutils.ParserError
	ExpectedYoutubeFeed *atom.Feed
}

func testYoutubeFeedValidator(actual xmlutils.Visitor, expected xmlutils.Visitor) error {
	f1 := actual.(*atom.Feed)
	f2 := expected.(*atom.Feed)

	channelId1, ok1 := GetFeedChannelId(f1)
	channelId2, ok2 := GetFeedChannelId(f2)

	if ok1 != ok2 {
		return fmt.Errorf("youtube channelId presence does not match %v (expected) vs %v", ok2, ok1)
	}

	if ok1 {
		if channelId1.String() != channelId2.String() {
			return fmt.Errorf("youtube channelId do not match '%s' (expected) vs '%s'", channelId2.String(), channelId1.String())
		}
	}

	if len(f1.Entries) != len(f2.Entries) {
		return fmt.Errorf("Feed does not contain the right count of Entries %v (expected) vs %v", len(f2.Entries), len(f1.Entries))
	}

	for i, _ := range f1.Entries {
		if err := testYoutubeEntryValidator(f1.Entries[i], f2.Entries[i]); err != nil {
			return err
		}

	}

	return nil
}

func testYoutubeEntryValidator(actual xmlutils.Visitor, expected xmlutils.Visitor) error {
	e1 := actual.(*atom.Entry)
	e2 := expected.(*atom.Entry)

	videoId1, ok1 := GetVideoId(e1)
	videoId2, ok2 := GetVideoId(e2)

	if ok1 != ok2 {
		return fmt.Errorf("youtube videoId presence does not match %v (expected) vs %v", ok2, ok1)
	}

	if ok1 {
		if videoId1.String() != videoId2.String() {
			return fmt.Errorf("youtube videoId do not match '%s' (expected) vs '%s'", videoId2.String(), videoId1.String())
		}
	}

	channelId1, ok1 := GetEntryChannelId(e1)
	channelId2, ok2 := GetEntryChannelId(e2)

	if ok1 != ok2 {
		return fmt.Errorf("youtube channelId presence does not match %v (expected) vs %v", ok2, ok1)
	}

	if ok1 {
		if channelId1.String() != channelId2.String() {
			return fmt.Errorf("youtube channelId do not match '%s' (expected) vs '%s'", channelId2.String(), channelId1.String())
		}
	}

	return nil
}

func testYoutubeFeedConstructor() xmlutils.Visitor {
	manager := extension.Manager{}
	AddToManager(&manager)

	return atom.NewFeedExt(manager)
}

func _TestYoutubeFeedToTestVisitor(t testYoutubeFeed) xmlutils.TestVisitor {
	customError := xmlutils.NewErrorChecker(xmlutils.DisableAllError)

	customError.EnableErrorChecking("entry", atom.AttributeDuplicated)

	testVisitor := xmlutils.TestVisitor{
		XML:                t.XML,
		ExpectedError:      nil,
		ExpectedVisitor:    t.ExpectedYoutubeFeed,
		VisitorConstructor: testYoutubeFeedConstructor,
		Validator:          testYoutubeFeedValidator,
		CustomError:        &customError,
	}

	if t.ExpectedError != nil {
		testVisitor.ExpectedError = t.ExpectedError
	}

	return testVisitor
}

func TestYoutubeFeedBasic(t *testing.T) {

	var testdata = []testYoutubeFeed{
		{`
		<feed xmlns:yt="http://www.youtube.com/xml/schemas/2015">
		  <yt:channelId>UCwg2zBWt4C55xcRKG-ThmXQ</yt:channelId>
		  <entry>
		    <id>yt:video:aF4JE5XmkfY</id>
		    <yt:videoId>aF4JE5XmkfY</yt:videoId>
		    <yt:channelId>UCwg2zBWt4C55xcRKG-ThmXQ</yt:channelId>
        </entry>
		  <entry>
		    <id>yt:video:badaboum</id>
		    <yt:videoId>badaboum</yt:videoId>
		    <yt:channelId>UCwg2zBWt4C55xcRKG-ThmXQ</yt:channelId>
        </entry>
		</feed>`,
			nil,
			NewTestYoutubeFeed("UCwg2zBWt4C55xcRKG-ThmXQ", []*atom.Entry{
				NewTestYoutubeEntry(
					"aF4JE5XmkfY",
					"UCwg2zBWt4C55xcRKG-ThmXQ",
				),
				NewTestYoutubeEntry(
					"badaboum",
					"UCwg2zBWt4C55xcRKG-ThmXQ",
				),
			}),
		},
		{`
		<entry xmlns:yt="http://www.youtube.com/xml/schemas/2015">
		  <id>yt:video:aF4JE5XmkfY</id>
		  <yt:videoId>aF4JE5XmkfY</yt:videoId>
		  <yt:videoId>aF4JE5XmkfY</yt:videoId>
		  <yt:channelId>UCwg2zBWt4C55xcRKG-ThmXQ</yt:channelId>
		</entry>
		 `,
			xmlutils.NewError(atom.AttributeDuplicated, ""),
			NewTestYoutubeFeed("", []*atom.Entry{
				NewTestYoutubeEntry(
					"aF4JE5XmkfY",
					"UCwg2zBWt4C55xcRKG-ThmXQ",
				)}),
		},
	}

	nbErrors := 0
	len := len(testdata)
	for _, testentry := range testdata {
		testcase := _TestYoutubeFeedToTestVisitor(testentry)

		if err := testcase.CheckTestCase(); err != nil {
			t.Errorf("FAIL\n%s\nXML:\n %s\n", err, testcase.XML)
			nbErrors++
		}
	}
	t.Logf("PASS RATIO = %v/%v\n", len-nbErrors, len)
}
