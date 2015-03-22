package feed

import (
	"time"

	"github.com/JLoup/xml/feed/atom"
	"github.com/JLoup/xml/feed/rss"
)

// BasicEntryBlock is a common brick to build UserFeed
type BasicEntryBlock struct {
	Title   string
	Link    string
	Date    time.Time
	Id      string
	Summary string
}

// BasicFeedBlock is a common brick to build UserFeed
type BasicFeedBlock struct {
	Title string
	Id    string
	Date  time.Time
}

// BasicFeed is a basic UserFeed
type BasicFeed struct {
	BasicFeedBlock
	Entries []BasicEntryBlock
}

func (b *BasicFeed) PopulateFromRssItem(i *rss.Item) {
	newEntry := BasicEntryBlock{}
	newEntry.PopulateFromRssItem(i)

	b.Entries = append(b.Entries, newEntry)
}

func (b *BasicFeed) PopulateFromAtomEntry(e *atom.Entry) {
	newEntry := BasicEntryBlock{}
	newEntry.PopulateFromAtomEntry(e)

	b.Entries = append(b.Entries, newEntry)
}

func (b *BasicFeedBlock) PopulateFromAtomFeed(f *atom.Feed) {
	b.Title = f.Title.String()
	b.Date = f.Updated.Time
	b.Id = f.Id.String()
}

func (b *BasicEntryBlock) PopulateFromAtomEntry(e *atom.Entry) {
	b.Title = e.Title.String()
	b.Id = e.Id.String()
	b.Date = e.Updated.Time
	b.Summary = e.Summary.String()

	for _, link := range e.Links {
		if link.Rel.String() == "alternate" {
			b.Link = link.Href.String()
		}
	}
}

func (b *BasicFeedBlock) PopulateFromRssChannel(c *rss.Channel) {
	b.Title = c.Title.String()
	b.Date = c.LastBuildDate.Time
	b.Id = c.Link.String()
}

func (b *BasicEntryBlock) PopulateFromRssItem(item *rss.Item) {
	b.Title = item.Title.String()
	b.Link = item.Link.String()
	b.Id = item.Guid.Content.String()
	b.Date = item.PubDate.Time
	b.Summary = item.Description.String()
}
