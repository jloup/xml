package feed_test

import (
	"fmt"
	"os"

	"github.com/JLoup/xml/feed"
	"github.com/JLoup/xml/feed/atom"
	"github.com/JLoup/xml/feed/rss"
	"github.com/JLoup/xml/feed/rss/extension/dc"

	"github.com/JLoup/xml/feed/extension"
)

type ExtendedFeed struct {
	feed.BasicFeedBlock
	Entries []ExtendedEntry
}

type ExtendedEntry struct {
	feed.BasicEntryBlock
	Creator string // <dc:creator> only present in RSS feeds
	Entries []feed.BasicEntryBlock
}

func (f *ExtendedFeed) PopulateFromAtomEntry(e *atom.Entry) {
	newEntry := ExtendedEntry{}
	newEntry.PopulateFromAtomEntry(e)
	f.Entries = append(f.Entries, newEntry)
}

func (f *ExtendedFeed) PopulateFromRssItem(i *rss.Item) {
	newEntry := ExtendedEntry{}
	newEntry.PopulateFromRssItem(i)

	creator, ok := dc.GetCreator(i)
	// we must check the item actually has a dc:creator element
	if ok {
		newEntry.Creator = creator.String()
	}
	f.Entries = append(f.Entries, newEntry)

}

func ExampleUserFeed_extension() {
	f, err := os.Open("testdata/rss.xml")

	if err != nil {
		return
	}

	//Manager is in github.com/JLoup/xml/feed/extension
	manager := extension.Manager{}
	// we add the dc extension to it
	// dc extension is in "github.com/JLoup/xml/feed/rss/extension/dc"
	dc.AddToManager(&manager)

	opt := feed.DefaultOptions
	//we pass our custom extension Manager to ParseOptions
	opt.ExtensionManager = manager

	myfeed := &ExtendedFeed{}
	err = feed.ParseCustom(f, myfeed, opt)

	if err != nil {
		fmt.Printf("Cannot parse feed: %s\n", err)
		return
	}

	fmt.Printf("FEED '%s'\n", myfeed.Title)
	for i, entry := range myfeed.Entries {
		fmt.Printf("\t#%v '%s' by %s (%s)\n", i, entry.Title, entry.Creator, entry.Link)
	}

	// Output:
	//FEED 'Me, Myself and I'
	//	#0 'Breakfast' by Peter J. (http://example.org/2005/04/02/breakfast)
	//	#1 'Dinner' by Peter J. (http://example.org/2005/04/02/dinner)
}
