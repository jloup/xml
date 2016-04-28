package youtube

import (
	"github.com/jloup/xml/feed/atom"
	"github.com/jloup/xml/feed/extension"
	xmlutils "github.com/jloup/xml/utils"
)

const NS = "http://www.youtube.com/xml/schemas/2015"

func AddToManager(manager *extension.Manager) {
	manager.AddElementExtension("entry", _videoId, newVideoIdElement, xmlutils.UniqueValidator(atom.AttributeDuplicated))
	manager.AddElementExtension("entry", _channelId, newChannelIdElement, xmlutils.UniqueValidator(atom.AttributeDuplicated))
	manager.AddElementExtension("feed", _channelId, newChannelIdElement, xmlutils.UniqueValidator(atom.AttributeDuplicated))
	manager.AddElementExtension("feed", _playlistId, newPlaylistIdElement, xmlutils.UniqueValidator(atom.AttributeDuplicated))
}

func GetVideoId(e *atom.Entry) (*atom.BasicElement, bool) {
	itf, ok := e.Extension.Store.GetItf(_videoId)
	if !ok {
		return nil, false
	}
	i, ok := itf.(*atom.BasicElement)
	return i, ok
}

func GetEntryChannelId(e *atom.Entry) (*atom.BasicElement, bool) {
	itf, ok := e.Extension.Store.GetItf(_channelId)
	if !ok {
		return nil, false
	}
	i, ok := itf.(*atom.BasicElement)
	return i, ok
}

func GetFeedChannelId(f *atom.Feed) (*atom.BasicElement, bool) {
	itf, ok := f.Extension.Store.GetItf(_channelId)
	if !ok {
		return nil, false
	}
	i, ok := itf.(*atom.BasicElement)
	return i, ok
}

func GetFeedPlaylistId(f *atom.Feed) (*atom.BasicElement, bool) {
	itf, ok := f.Extension.Store.GetItf(_playlistId)
	if !ok {
		return nil, false
	}
	i, ok := itf.(*atom.BasicElement)
	return i, ok
}
