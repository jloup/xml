// Package dc implements dc:creator extension (http://purl.org/dc/elements/1.1/) for RSS feed
package dc

import (
	"github.com/jloup/xml/feed/extension"
	"github.com/jloup/xml/feed/rss"
	"github.com/jloup/xml/utils"
)

const NS = "http://purl.org/dc/elements/1.1/"

func AddToManager(manager *extension.Manager) {
	manager.AddElementExtension("item", CREATOR, NewCreatorElement, utils.UniqueValidator(rss.AttributeDuplicated))

}

func GetCreator(item *rss.Item) (*rss.BasicElement, bool) {
	itf, ok := item.Extension.Store.GetItf(CREATOR)
	if !ok {
		return nil, false
	}
	i, ok := itf.(*rss.BasicElement)
	return i, ok
}
