// Package thr implements threading extension (http://purl.org/syndication/thread/1.0) for atom feed
package thr

import (
	"github.com/jloup/xml/feed/atom"
	"github.com/jloup/xml/feed/extension"
	"github.com/jloup/xml/utils"
)

const NS = "http://purl.org/syndication/thread/1.0"

func AddToManager(manager *extension.Manager) {
	manager.AddElementExtension("entry", _inreplyto, newInReplyToElement, utils.UniqueValidator(atom.AttributeDuplicated))
	manager.AddElementExtension("entry", _total, newTotalElement, utils.UniqueValidator(atom.AttributeDuplicated))
	manager.AddAttrExtension("link", _count, newCountAttr, atom.AttributeDuplicated)
	manager.AddAttrExtension("link", _updated, newUpdatedAttr, atom.AttributeDuplicated)

}

func GetInReplyTo(e *atom.Entry) (*InReplyTo, bool) {

	itf, ok := e.Extension.Store.GetItf(_inreplyto)
	if !ok {
		return nil, false
	}
	i, ok := itf.(*InReplyTo)
	return i, ok
}

func GetTotal(e *atom.Entry) (*atom.BasicElement, bool) {

	itf, ok := e.Extension.Store.GetItf(_total)
	if !ok {
		return nil, false
	}
	i, ok := itf.(*atom.BasicElement)
	return i, ok
}

func GetCount(l *atom.Link) (*Count, bool) {
	itf, ok := l.Extension.Store.GetItf(_count)
	if !ok {
		return nil, false
	}

	i, ok := itf.(*Count)

	return i, ok

}

func GetUpdated(l *atom.Link) (*Updated, bool) {
	itf, ok := l.Extension.Store.GetItf(_updated)
	if !ok {
		return nil, false
	}
	i, ok := itf.(*Updated)

	return i, ok

}
