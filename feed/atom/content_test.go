package atom

import (
	"fmt"
	"testing"

	"github.com/JLoup/xml/helper"
)

func NewTestContent(contentType, xhtml, plaintext, other, outsource string) *Content {
	t := NewContent()

	t.Type.Value = contentType
	t.PlainText.Content = plaintext
	t.XHTML.Content.WriteString(xhtml)
	t.InlineContent.Content.WriteString(other)
	t.OutOfLineContent.Src.Value = outsource

	return t
}

func ContentWithBaseLang(c *Content, lang, base string) *Content {
	c.Lang.Value = lang
	c.Base.Value = base

	return c
}

type testContent struct {
	XML             string
	ExpectedError   helper.ParserError
	ExpectedContent *Content
}

func testContentValidator(actual helper.Visitor, expected helper.Visitor) error {
	c1 := actual.(*Content)
	c2 := expected.(*Content)

	if c1.Type.Value != c2.Type.Value {
		return fmt.Errorf("Content type is invalid '%s' (expected) vs '%s'", c2.Type.Value, c1.Type.Value)
	}

	if err := ValidateBaseLang("Content", c1.Base.Value, c1.Lang.Value, c2.Base.Value, c2.Lang.Value); err != nil {
		return err
	}

	if c1.PlainText.String() != c2.PlainText.String() {
		return fmt.Errorf("Text content is invalid '%s' (expected) vs '%s'", c2.PlainText.String(), c1.PlainText.String())
	}

	if c1.XHTML.String() != c2.XHTML.String() {
		return fmt.Errorf("XHTML content is invalid '%s' (expected) vs '%s'", c2.XHTML.String(), c1.XHTML.String())
	}

	if c1.InlineContent.String() != c2.InlineContent.String() {
		return fmt.Errorf("Inline Content content is invalid '%s' (expected) vs '%s'", c2.InlineContent.String(), c1.InlineContent.String())
	}

	if c1.OutOfLineContent.Src.Value != c2.OutOfLineContent.Src.Value {
		return fmt.Errorf("Out of Line source is invalid '%s' (expected) vs '%s'", c2.OutOfLineContent.Src.Value, c1.OutOfLineContent.Src.Value)
	}

	return nil
}

func testContentConstructor() helper.Visitor {
	return NewContent()
}

func _TestContentToTestVisitor(t testContent) helper.TestVisitor {
	testVisitor := helper.TestVisitor{
		XML:                t.XML,
		ExpectedError:      nil,
		ExpectedVisitor:    t.ExpectedContent,
		VisitorConstructor: testContentConstructor,
		Validator:          testContentValidator,
	}

	if t.ExpectedError != nil {
		testVisitor.ExpectedError = t.ExpectedError
	}

	return testVisitor
}

func TestContentBasic(t *testing.T) {

	var testdata = []testContent{
		{`<content xml:lang="en-us" xml:base="http://yo.com" src="image.png" type="image/png"/>`,
			nil,
			ContentWithBaseLang(NewTestContent("image/png", "", "", "", "image.png"), "en-us", "http://yo.com"),
		},
		{`
<content type="image/png">iVBORw0KGgoAAAANSUhEUgAAAB8AAAAqCAYAAABLGYAnAAAABmJLR0QA/wD/AP+gvaeTAAAACXBIWXMAAAsTAAALEwEAmpwYAAAAB3RJTUUH1QwCBCUlRSCuygAAAetJREFUWMPt1j1IVmEUB/Dfo1mvBg1iSgjhZNGnFIhTDUXQ25x9LE2NbQ1NLe1tbbo1tjQEDk2BUzSITUapUejyvkQthXKfhvvcEs3gvl55Fd4/HLgfz//8zzn3Ps85tBGhBc4oLuIYuvAV77CwW0F24x7mEbex+bSmu0rhEcz9R3SzzSXOjnEWjRLChTUSt2X0Y7EF4cIWk4+WMLUD4cKmtPhHr1cgvp58/RNd2zy/VdFf2518lcJsBVkXNls285GKt2qpE+4nDlUk/gu1Mpk3Ksy8UbbsSxWKL5UVn6lQfGZP7vM9ecK1/Wxva1drez/f1UlmX8xwHXTQQQcd7Ctkl8jGiMeJk8SeDe9uEI+Q1bfyYg/xNrGf7CaxRhzP77esPUy8soFXLwbIU4QaPuM6YY04SDxAOI2MMJGIw8Q0z4c14lg+k4dr6MEAoUkcIvYlzgAO4nKK5CgmknjIUou8mvflrJ6uH+Vt+k/09UR88rc64Xnq4Su4gzXiKMbxgHgBTzGcfEyirxivitH5LeF1yvIkXuLZptKdwQ98Q28Sf58ymsbd1NdPJMKHlPFCWgfnsIxmEo/NvHRxCB/xgng/ZfkFg1glTBPP4w3h+4agHhOW8TAvuVcpuBV8yrmxl7iKVKm4mn/WDtqA3yOQKuHaSApTAAAAAElFTkSuQmCC</content>`,
			nil,
			NewTestContent("image/png", "", "", "iVBORw0KGgoAAAANSUhEUgAAAB8AAAAqCAYAAABLGYAnAAAABmJLR0QA/wD/AP+gvaeTAAAACXBIWXMAAAsTAAALEwEAmpwYAAAAB3RJTUUH1QwCBCUlRSCuygAAAetJREFUWMPt1j1IVmEUB/Dfo1mvBg1iSgjhZNGnFIhTDUXQ25x9LE2NbQ1NLe1tbbo1tjQEDk2BUzSITUapUejyvkQthXKfhvvcEs3gvl55Fd4/HLgfz//8zzn3Ps85tBGhBc4oLuIYuvAV77CwW0F24x7mEbex+bSmu0rhEcz9R3SzzSXOjnEWjRLChTUSt2X0Y7EF4cIWk4+WMLUD4cKmtPhHr1cgvp58/RNd2zy/VdFf2518lcJsBVkXNls285GKt2qpE+4nDlUk/gu1Mpk3Ksy8UbbsSxWKL5UVn6lQfGZP7vM9ecK1/Wxva1drez/f1UlmX8xwHXTQQQcd7Ctkl8jGiMeJk8SeDe9uEI+Q1bfyYg/xNrGf7CaxRhzP77esPUy8soFXLwbIU4QaPuM6YY04SDxAOI2MMJGIw8Q0z4c14lg+k4dr6MEAoUkcIvYlzgAO4nKK5CgmknjIUou8mvflrJ6uH+Vt+k/09UR88rc64Xnq4Su4gzXiKMbxgHgBTzGcfEyirxivitH5LeF1yvIkXuLZptKdwQ98Q28Sf58ymsbd1NdPJMKHlPFCWgfnsIxmEo/NvHRxCB/xgng/ZfkFg1glTBPP4w3h+4agHhOW8TAvuVcpuBV8yrmxl7iKVKm4mn/WDtqA3yOQKuHaSApTAAAAAElFTkSuQmCC", ""),
		},
		{`
    <content type="xhtml"><div xmlns='http://www.w3.org/1999/xhtml'><p>I got<a href="/ongoing/When/201x/2014/03/19/Keybase">interested in Keybase.io</a> the day I left Google in March, and Iâ€™ve been evangelizing it, but even more the idea behind it: Usingauthenticated posts here and there to prove public-key ownership. Also Iâ€™vecontributed Keybase-client code to<a href="http://www.openkeychain.org/">OpenKeychain</a> (letâ€™s just say â€œOKCâ€), a pretty good Android crypto app.  Iâ€™m more or less done now.</p></div></content>`,
			nil,
			NewTestContent("xhtml", "<div><p>I got<a href=\"/ongoing/When/201x/2014/03/19/Keybase\">interested in Keybase.io</a> the day I left Google in March, and Iâ€™ve been evangelizing it, but even more the idea behind it: Usingauthenticated posts here and there to prove public-key ownership. Also Iâ€™vecontributed Keybase-client code to<a href=\"http://www.openkeychain.org/\">OpenKeychain</a> (letâ€™s just say â€œOKCâ€), a pretty good Android crypto app.  Iâ€™m more or less done now.</p></div>", "", "", ""),
		},
		{`<content src="image.png" type="image/png"/>`,
			nil,
			NewTestContent("image/png", "", "", "", "image.png"),
		},
		{`<content src="/iri" type="text"/>`,
			helper.NewError(ContentTypeIsNotValid, ""),
			NewTestContent("text", "", "", "", "/iri"),
		},
		{`<content src="/iri" type="html"/>`,
			helper.NewError(ContentTypeIsNotValid, ""),
			NewTestContent("html", "", "", "", "/iri"),
		},
		{`<content src="/iri" type="xhtml"/>`,
			helper.NewError(ContentTypeIsNotValid, ""),
			NewTestContent("xhtml", "", "", "", "/iri"),
		},
		{`<content src="http://%g.com" type="application/xml"/>`,
			helper.NewError(IriNotValid, ""),
			NewTestContent("application/xml", "", "", "", "http://%g.com"),
		},

		{`<content src="http://g.com" type="application/xml">text</content>`,
			helper.NewError(SourcedContentElementNotEmpty, ""),
			NewTestContent("application/xml", "", "", "", "http://g.com"), // content is ignored
		},
		{`<content src="http://g.com" type="application/xml"><div>CHILD</div></content>`,
			helper.NewError(SourcedContentElementNotEmpty, ""),
			NewTestContent("application/xml", "", "", "", "http://g.com"), // content is ignored
		},
		{`<content src="http://g.com" type="xml"/>`,
			helper.NewError(IsNotMIME, ""),
			NewTestContent("xml", "", "", "", "http://g.com"),
		},
		{`<content type="text/plain">text</content>`,
			nil,
			NewTestContent("text/plain", "", "text", "", ""),
		},
		{`<content type="text/plain"><div>CHILD</div></content>`,
			helper.NewError(LeafElementHasChild, ""),
			NewTestContent("text/plain", "", "CHILD", "", ""),
		},
		{`<content type="text"><div>CHILD</div><p>CHILD2</p></content>`,
			helper.NewError(LeafElementHasChild, ""),
			NewTestContent("text", "", "CHILD2", "", ""),
		},
		{`<content type="html"><div>CHILD</div></content>`,
			helper.NewError(LeafElementHasChild, ""),
			NewTestContent("html", "", "CHILD", "", ""),
		},
		{`<content type="xhtml"><div xmlns="http://www.w3.org/1999/xhtml">XHTML CHILD</div></content>`,
			nil,
			NewTestContent("xhtml", "<div>XHTML CHILD</div>", "", "", ""),
		},
		{`<content type="xhtml">XHTML CHILD</content>`,
			helper.NewError(XHTMLRootNodeNotDiv, ""),
			NewTestContent("xhtml", "XHTML CHILD", "", "", ""),
		},
		{`<content type="multipart/mixed"></content>`,
			helper.NewError(IsNotMIME, ""),
			NewTestContent("multipart/mixed", "", "", "", ""),
		},
	}

	nbErrors := 0
	len := len(testdata)
	for _, testcontent := range testdata {
		testcase := _TestContentToTestVisitor(testcontent)

		if err := testcase.CheckTestCase(); err != nil {
			t.Errorf("FAIL\n%s\nXML:\n %s\n", err, testcase.XML)
			nbErrors++
		}
	}

	t.Logf("PASS RATIO = %v/%v\n", len-nbErrors, len)
}
