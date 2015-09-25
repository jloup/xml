package feed

import (
	"encoding/xml"

	"github.com/jloup/xml/feed/atom"
	"github.com/jloup/xml/feed/extension"
	"github.com/jloup/xml/feed/rss"
	xmlutils "github.com/jloup/xml/utils"
)

type wrapper struct {
	AtomFeed   *atom.Feed
	AtomEntry  *atom.Entry
	RssChannel *rss.Channel

	Extensions extension.VisitorExtension
}

func newWrapper() *wrapper {
	return &wrapper{}
}

func newWrapperExt(manager extension.Manager) *wrapper {
	w := newWrapper()

	w.Extensions = extension.InitExtension("__", manager)

	return w
}

func (w *wrapper) ProcessStartElement(el xmlutils.StartElement) (xmlutils.Visitor, xmlutils.ParserError) {
	switch el.Name.Local {
	case "feed":
		w.AtomFeed = atom.NewFeedExt(w.Extensions.Manager)
		w.AtomFeed.Parent = w
		return w.AtomFeed.ProcessStartElement(el)

	case "entry":
		w.AtomEntry = atom.NewEntryExt(w.Extensions.Manager)
		w.AtomEntry.Parent = w
		return w.AtomEntry.ProcessStartElement(el)

	case "channel":
		w.RssChannel = rss.NewChannelExt(w.Extensions.Manager)
		w.RssChannel.Parent = w
		return w.RssChannel.ProcessStartElement(el)
	}

	return w, nil

}

func (w *wrapper) ProcessEndElement(el xml.EndElement) (xmlutils.Visitor, xmlutils.ParserError) {
	return w, nil
}

func (w *wrapper) ProcessCharData(el xml.CharData) (xmlutils.Visitor, xmlutils.ParserError) {
	return w, nil
}

func (w *wrapper) Populate(u UserFeed) {
	if w.AtomFeed != nil {
		u.PopulateFromAtomFeed(w.AtomFeed)

		for _, entry := range w.AtomFeed.Entries {
			u.PopulateFromAtomEntry(entry)
		}

	} else if w.AtomEntry != nil {
		u.PopulateFromAtomEntry(w.AtomEntry)

	} else if w.RssChannel != nil {
		u.PopulateFromRssChannel(w.RssChannel)

		for _, item := range w.RssChannel.Items {
			u.PopulateFromRssItem(item)
		}

	}
}
