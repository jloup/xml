package rss

import (
	xmlutils "github.com/jloup/xml/utils"
	"fmt"
	"testing"
)

func NewTestImage(url, title, link, width, height, description string) *Image {
	i := NewImage()

	i.Url.Content.Value = url
	i.Title.Content.Value = title
	i.Link.Content.Value = link
	i.Width.Content.Value = width
	i.Height.Content.Value = height
	i.Description.Content.Value = description

	return i
}

type testImage struct {
	XML           string
	ExpectedError xmlutils.ParserError
	ExpectedImage *Image
}

func testImageValidator(actual xmlutils.Visitor, expected xmlutils.Visitor) error {
	i1 := actual.(*Image)
	i2 := expected.(*Image)

	if i1.Url.Content.Value != i2.Url.Content.Value {
		return fmt.Errorf("Url is invalid '%s' (expected) vs '%s'", i2.Url.Content.Value, i1.Url.Content.Value)
	}

	if i1.Title.Content.Value != i2.Title.Content.Value {
		return fmt.Errorf("Title is invalid '%s' (expected) vs '%s'", i2.Title.Content.Value, i1.Title.Content.Value)
	}

	if i1.Link.Content.Value != i2.Link.Content.Value {
		return fmt.Errorf("Link is invalid '%s' (expected) vs '%s'", i2.Link.Content.Value, i1.Link.Content.Value)
	}

	if i1.Width.Content.Value != i2.Width.Content.Value {
		return fmt.Errorf("Width is invalid '%s' (expected) vs '%s'", i2.Width.Content.Value, i1.Width.Content.Value)
	}

	if i1.Height.Content.Value != i2.Height.Content.Value {
		return fmt.Errorf("Height is invalid '%s' (expected) vs '%s'", i2.Height.Content.Value, i1.Height.Content.Value)
	}

	if i1.Description.Content.Value != i2.Description.Content.Value {
		return fmt.Errorf("Description is invalid '%s' (expected) vs '%s'", i2.Description.Content.Value, i1.Description.Content.Value)
	}

	return nil
}

func testImageConstructor() xmlutils.Visitor {
	return NewImage()
}

func _TestImageToTestVisitor(t testImage) xmlutils.TestVisitor {
	testVisitor := xmlutils.TestVisitor{
		XML:                t.XML,
		ExpectedError:      nil,
		ExpectedVisitor:    t.ExpectedImage,
		VisitorConstructor: testImageConstructor,
		Validator:          testImageValidator,
	}

	if t.ExpectedError != nil {
		testVisitor.ExpectedError = t.ExpectedError
	}

	return testVisitor
}

var testdata = []testImage{
	{`
         <image>
           <title>WriteTheWeb</title>
           <url>http://writetheweb.com/images/mynetscape88.gif</url>
           <link>http://writetheweb.com</link>
           <width>88</width>
           <height>31</height>
           <description>News for web users that write back</description>
          </image>`,
		nil,
		NewTestImage("http://writetheweb.com/images/mynetscape88.gif", "WriteTheWeb", "http://writetheweb.com", "88", "31", "News for web users that write back"),
	},
	{`
         <image>
           <url>http://writetheweb.com/images/mynetscape88.gif</url>
           <link>http://writetheweb.com</link>
           <width>88</width>
           <height>31</height>
           <description>News for web users that write back</description>
          </image>`,
		xmlutils.NewError(MissingAttribute, ""),
		NewTestImage("http://writetheweb.com/images/mynetscape88.gif", "", "http://writetheweb.com", "88", "31", "News for web users that write back"),
	},
	{`
         <image>
           <title>WriteTheWeb</title>
           <link>http://writetheweb.com</link>
           <width>88</width>
           <height>31</height>
           <description>News for web users that write back</description>
          </image>`,
		xmlutils.NewError(MissingAttribute, ""),
		NewTestImage("", "WriteTheWeb", "http://writetheweb.com", "88", "31", "News for web users that write back"),
	},
	{`
         <image>
           <title>WriteTheWeb</title>
           <url>http://writetheweb.com/images/mynetscape88.gif</url>
           <width>88</width>
           <height>31</height>
           <description>News for web users that write back</description>
          </image>`,
		xmlutils.NewError(MissingAttribute, ""),
		NewTestImage("http://writetheweb.com/images/mynetscape88.gif", "WriteTheWeb", "", "88", "31", "News for web users that write back"),
	},
	{`
         <image>
           <title>WriteTheWeb</title>
           <url>http://writetheweb.com/images/mynetscape88.gif</url>
           <link>http://writetheweb.com</link>
          </image>`,
		nil,
		NewTestImage("http://writetheweb.com/images/mynetscape88.gif", "WriteTheWeb", "http://writetheweb.com", "", "", ""),
	},
}

func TestImageBasic(t *testing.T) {
	nbErrors := 0
	len := len(testdata)
	for _, testimage := range testdata {
		testcase := _TestImageToTestVisitor(testimage)

		if err := testcase.CheckTestCase(); err != nil {
			t.Errorf("FAIL\n%s\nXML:\n %s\n", err, testcase.XML)
			nbErrors++
		}
	}

	t.Logf("PASS RATIO = %v/%v\n", len-nbErrors, len)
}
