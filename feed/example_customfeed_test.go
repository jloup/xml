package feed_test

import (
	"fmt"
	"os"

	"github.com/JLoup/xml/feed"
	"github.com/JLoup/xml/feed/atom"
	"github.com/JLoup/xml/feed/rss"
)

type MyFeed struct {
	feed.BasicFeedBlock
	Generator string
	Entries   []feed.BasicEntryBlock
}

func (m *MyFeed) PopulateFromAtomFeed(f *atom.Feed) {
	m.BasicFeedBlock.PopulateFromAtomFeed(f)

	m.Generator = fmt.Sprintf("%s V%s", f.Generator.Uri.String(), f.Generator.Version.String())
}

func (m *MyFeed) PopulateFromRssChannel(c *rss.Channel) {
	m.BasicFeedBlock.PopulateFromRssChannel(c)

	m.Generator = c.Generator.String()
}

func (m *MyFeed) PopulateFromAtomEntry(e *atom.Entry) {
	newEntry := feed.BasicEntryBlock{}
	newEntry.PopulateFromAtomEntry(e)
	m.Entries = append(m.Entries, newEntry)
}

func (m *MyFeed) PopulateFromRssItem(i *rss.Item) {
	newEntry := feed.BasicEntryBlock{}
	newEntry.PopulateFromRssItem(i)
	m.Entries = append(m.Entries, newEntry)

}

func ExampleParseCustom() {
	f, err := os.Open("testdata/atom.xml")

	if err != nil {
		return
	}

	myfeed := &MyFeed{}

	err = feed.ParseCustom(f, myfeed, feed.DefaultOptions)

	if err != nil {
		fmt.Printf("Cannot parse feed: %s\n", err)
		return
	}

	fmt.Printf("FEED '%s' generated with %s\n", myfeed.Title, myfeed.Generator)
	for i, entry := range myfeed.Entries {
		fmt.Printf("\t#%v '%s' (%s)\n", i, entry.Title, entry.Link)
	}

	// Output:
	//FEED 'Me, Myself and I' generated with http://www.atomgenerator.com/ V1.0
	//	#0 'Breakfast' (http://example.org/2005/04/02/breakfast)
	//	#1 'Dinner' (http://example.org/2005/04/02/dinner)
}
