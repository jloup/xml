package feed_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jloup/xml/feed"
	"github.com/jloup/xml/feed/atom"

	"github.com/jloup/xml/feed/extension"
	xmlutils "github.com/jloup/xml/utils"
)

func ExampleParse() {
	f, err := os.Open("testdata/atom.xml")

	if err != nil {
		return
	}

	myfeed, err := feed.Parse(f, feed.DefaultOptions)

	if err != nil {
		fmt.Printf("Cannot parse feed: %s\n", err)
		return
	}

	fmt.Printf("FEED '%s'\n", myfeed.Title)
	for i, entry := range myfeed.Entries {
		fmt.Printf("\t#%v '%s' (%s)\n\t\t%s\n\n", i, entry.Title, entry.Link, entry.Summary)
	}

	// Output:
	//FEED 'Me, Myself and I'
	//	#0 'Breakfast' (http://example.org/2005/04/02/breakfast)
	//		eggs and bacon, yup !
	//
	//	#1 'Dinner' (http://example.org/2005/04/02/dinner)
	//		got soap delivered !
}

func ExampleParseOptions() {
	// the input feed is not compliant to spec
	f, err := os.Open("testdata/atom.xml")
	if err != nil {
		return
	}

	// the input feed should be 100% compliant to spec...
	flags := xmlutils.NewErrorChecker(xmlutils.EnableAllError)

	//... but it is OK if Atom entry does not have <updated> field
	flags.DisableErrorChecking("entry", atom.MissingDate)

	options := feed.ParseOptions{ExtensionManager: extension.Manager{},
		ErrorFlags: &flags}

	myfeed, err := feed.Parse(f, options)

	if err != nil {
		fmt.Printf("Cannot parse feed:\n%s\n", err)
		return
	}

	fmt.Printf("FEED '%s'\n", myfeed.Title)

	// Output:
	//Cannot parse feed:
	//in 'feed':
	//[MissingId]
	//	feed's id should exist
}

func ExampleParse_recover() {
	f, err := os.Open("testdata/invalid_atom.xml")

	if err != nil {
		return
	}

	b, _ := ioutil.ReadAll(f)
	opt := feed.DefaultOptions

	// no retry (default)
	fmt.Println("attempt #1 with XMLTokenErrorRetry=0")
	opt.XMLTokenErrorRetry = 0

	_, err = feed.Parse(bytes.NewReader(b), opt)

	if err != nil {
		fmt.Printf("\t->Cannot parse feed: %s\n", err)
	} else {
		fmt.Println("\t->no error")
	}

	// one retry
	fmt.Println("attempt #2 with XMLTokenErrorRetry=1")
	opt.XMLTokenErrorRetry = 1

	_, err = feed.Parse(bytes.NewReader(b), opt)

	if err != nil {
		fmt.Printf("\t->Cannot parse feed: %s\n", err)
	} else {
		fmt.Println("\t->no error")
	}

	// Output:
	//attempt #1 with XMLTokenErrorRetry=0
	//	->Cannot parse feed: [XMLTokenError] XML syntax error on line 574: illegal character code U+000C
	//attempt #2 with XMLTokenErrorRetry=1
	//	->no error
}
