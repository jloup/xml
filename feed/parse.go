//Package feed implements a flexible RSS/Atom parser
package feed

import (
	"io"

	"github.com/jloup/flag"
	"github.com/jloup/xml/feed/extension"
	"github.com/jloup/xml/utils"
)

var NoFeedFound = flag.InitFlag(&utils.ErrorFlagCounter, "NoFeedFound")

// ParseOptions is passed to Parse functions to customize their behaviors
type ParseOptions struct {
	ExtensionManager extension.Manager
	ErrorFlags       utils.FlagChecker
}

// DefaultOptions set options in order to have:
// - no specification checking
// - no extension
var DefaultOptions ParseOptions

func init() {
	errorFlags := utils.NewErrorChecker(utils.DisableAllError)
	DefaultOptions = ParseOptions{
		extension.Manager{},
		&errorFlags,
	}
}

// ParseCustom parse bytes from a io.Reader into a UserFeed object
func ParseCustom(r io.Reader, feed UserFeed, options ParseOptions) error {
	w := newWrapperExt(options.ExtensionManager)

	err := utils.Walk(r, w, options.ErrorFlags)

	if err != nil {
		return err
	}

	if w.AtomFeed == nil && w.RssChannel == nil && w.AtomEntry == nil {
		return utils.NewError(NoFeedFound, "no feed has been found")
	}

	w.Populate(feed)

	return nil

}

// Parse is a subset a ParseCustom with BasicFeed passed as UserFeed
func Parse(r io.Reader, options ParseOptions) (BasicFeed, error) {
	feed := BasicFeed{}

	return feed, ParseCustom(r, &feed, options)
}
