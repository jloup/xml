// Package dc implements dc:creator extension (http://purl.org/dc/elements/1.1/) for RSS feed
package dc

import (
	"github.com/JLoup/xml/feed/extension"
	"github.com/JLoup/xml/feed/rss"
	"github.com/JLoup/xml/helper"
)

const NS = "http://purl.org/dc/elements/1.1/"

func AddToManager(manager *extension.Manager) {
	manager.AddElementExtension("item", CREATOR, NewCreatorElement, helper.UniqueValidator(rss.AttributeDuplicated))

}

func GetCreator(item *rss.Item) (*rss.BasicElement, bool) {
	itf, ok := item.Extension.Store.GetItf(CREATOR)
	if !ok {
		return nil, false
	}
	i, ok := itf.(*rss.BasicElement)
	return i, ok
}
