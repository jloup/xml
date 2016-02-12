## Feed Parser (RSS, Atom)
[![GoDoc](https://godoc.org/github.com/jloup/xml/feed?status.svg)](https://godoc.org/github.com/jloup/xml/feed)
[![Travis Build Status](https://travis-ci.org/jloup/xml.svg?branch=master)](https://travis-ci.org/jloup/xml)

Package feed implements a flexible and efficient RSS/Atom parser. 

If you just want some bytes to be quickly parsed into an object without care about underlying feed type, you can start with this: [Simple Use](#simple)

If you want to take a deeper dive into how you can customize the parser behavior:
- [Extending BasicFeed](#userfeed)
- [Parse with specification compliancy checking](#spec)
- [RSS and Atom extensions](#extension)

#### Installation & Use

Get the pkg
```
go get github.com/jloup/xml
```

Use it in code
```
import "github.com/jloup/xml/feed"
```

#### <a name="simple"></a>Simple Use : feed.Parse(io.Reader, feed.DefaultOptions)
Example:
```go
f, err := os.Open("feed.txt")

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
    fmt.Printf("\t#%v '%s' (%s)\n\t\t%s\n\n", i, entry.Title,
                                                 entry.Link,
                                                 entry.Summary)
}
```
Output:
```
FEED 'Me, Myself and I'
	#0 'Breakfast' (http://example.org/2005/04/02/breakfast)
		eggs and bacon, yup !

	#1 'Dinner' (http://example.org/2005/04/02/dinner)
		got soap delivered !
```
feed.Parse returns a BasicFeed which fields are : 
```go
// Rss channel or Atom feed
type BasicFeed struct {
  Title   string
  Id      string // Atom:feed:id | RSS:channel:link 
  Date    time.Time
  Entries []BasicEntryBlock
}

type BasicEntryBlock struct {
	Title   string
	Link    string
	Date    time.Time // Atom:entry:updated | RSS:item:pubDate
	Id      string // Atom:entry:id | RSS:item:guid
	Summary string
}
```

#### <a name="userfeed"></a>Extending BasicFeed
BasicFeed is really basic struct implementing **feed.UserFeed** interface. You may want to access more values extracted from feeds. For this purpose you can pass your own implementation of feed.UserFeed to **feed.ParseCustom**.
```go

type UserFeed interface {
    PopulateFromAtomFeed(f *atom.Feed) // see github.com/jloup/xml/feed/atom
    PopulateFromAtomEntry(e *atom.Entry)
    PopulateFromRssChannel(c *rss.Channel) // see github.com/jloup/xml/feed/rss
    PopulateFromRssItem(i *rss.Item)
}

func ParseCustom(r io.Reader, feed UserFeed, options ParseOptions) error
```
To avoid starting from scratch, you can embed feed.BasicEntryBlock and feed.BasicFeedBlock in your structs

Example:
```go
type MyFeed struct {
	feed.BasicFeedBlock
	Generator string
	Entries   []feed.BasicEntryBlock
}

func (m *MyFeed) PopulateFromAtomFeed(f *atom.Feed) {
	m.BasicFeedBlock.PopulateFromAtomFeed(f)

	m.Generator = fmt.Sprintf("%s V%s", f.Generator.Uri.String(), 
	                                    f.Generator.Version.String())
}

func (m *MyFeed) PopulateFromRssChannel(c *rss.Channel) {
	m.BasicFeedBlock.PopulateFromRssChannel(c)

	m.Generator = c.Generator.String()
}

func (m *MyFeed) PopulateFromAtomEntry(e *atom.Entry) {
	newEntry := feed.BasicEntryBlock{}
	newEntry.PopulateFromAtomEntry(e)
	m.Entries = append(m.Entries, newEntry)
}

func (m *MyFeed) PopulateFromRssItem(i *rss.Item) {
	newEntry := feed.BasicEntryBlock{}
	newEntry.PopulateFromRssItem(i)
	m.Entries = append(m.Entries, newEntry)

}

func main() {
    f, err := os.Open("feed.txt")

    if err != nil {
        return
    }

    myfeed := &MyFeed{}

    err = feed.ParseCustom(f, myfeed, feed.DefaultOptions)

    if err != nil {
        fmt.Printf("Cannot parse feed: %s\n", err)
        return
    }

    fmt.Printf("FEED '%s' generated with %s\n", myfeed.Title, myfeed.Generator)
    for i, entry := range myfeed.Entries {
        fmt.Printf("\t#%v '%s' (%s)\n", i, entry.Title, entry.Link)
    }
}
```
Output:
```
FEED 'Me, Myself and I' generated with http://www.atomgenerator.com/ V1.0
	#0 'Breakfast' (http://example.org/2005/04/02/breakfast)
	#1 'Dinner' (http://example.org/2005/04/02/dinner)
```

#### <a name="spec"></a>Parse with specification compliancy checking
RSS and Atom feeds should conform to a specification (which is complex for Atom). The common behavior of Parse functions is to not be too restrictive about input feeds. To validate feeds, you can pass a custom FlagChecker to ParseOptions. If you really know what you are doing you can enable/disable only some spec checks.

Error flags can be found for each standard in packages documentation:
- RSS : github.com/jloup/xml/feed/rss
- Atom : github.com/jloup/xml/feed/atom

Example:
```go
// the input feed is not compliant to spec
f, err := os.Open("feed.txt")
if err != nil {
    return
}

// the input feed should be 100% compliant to spec...
flags := xmlutils.NewErrorChecker(xmlutils.EnableAllError)

//... but it is OK if Atom entry does not have <updated> field
flags.DisableErrorChecking("entry", atom.MissingDate)

options := feed.ParseOptions{extension.Manager{}, &flags}

myfeed, err := feed.Parse(f, options)

if err != nil {
    fmt.Printf("Cannot parse feed:\n%s\n", err)
    return
}

fmt.Printf("FEED '%s'\n", myfeed.Title)
```
Output:
```
Cannot parse feed:
in 'feed':
[MissingId]
	feed's id should exist
```

#### <a name="extension"></a>Rss and Atom extensions
Both formats allow to add third party extensions. Some extensions have been implemented for the example e.g. RSS dc:creator (github.com/jloup/xml/feed/rss/extension/dc)

Example:
```go
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

func main() {
     f, err := os.Open("rss.txt")

    if err != nil {
        return
    }

    //Manager is in github.com/jloup/xml/feed/extension
    manager := extension.Manager{}
    // we add the dc extension to it
    // dc extension is in "github.com/jloup/xml/feed/rss/extension/dc"
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
        fmt.Printf("\t#%v '%s' by %s (%s)\n", i, entry.Title,
                                                 entry.Creator,
                                                 entry.Link)
    }
}
```
Output:
```
FEED 'Me, Myself and I'
	#0 'Breakfast' by Peter J. (http://example.org/2005/04/02/breakfast)
	#1 'Dinner' by Peter J. (http://example.org/2005/04/02/dinner)
```
