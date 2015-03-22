package feed

import (
	"github.com/JLoup/xml/feed/atom"
	"github.com/JLoup/xml/feed/rss"
)

type UserFeed interface {
	PopulateFromAtomFeed(f *atom.Feed)
	PopulateFromAtomEntry(e *atom.Entry)
	PopulateFromRssChannel(c *rss.Channel)
	PopulateFromRssItem(i *rss.Item)
}
