package rss

import (
	"github.com/JLoup/xml/helper"
	"fmt"
	"testing"
)

func NewTestCloud(domain, port, path, registerProcedure, protocol string) *Cloud {
	c := NewCloud()

	c.Domain.Value = domain
	c.Port.Value = port
	c.Path.Value = path
	c.RegisterProcedure.Value = registerProcedure
	c.Protocol.Value = protocol

	return c
}

type testCloud struct {
	XML           string
	ExpectedError helper.ParserError
	ExpectedCloud *Cloud
}

func testCloudValidator(actual helper.Visitor, expected helper.Visitor) error {
	c1 := actual.(*Cloud)
	c2 := expected.(*Cloud)

	if c1.Domain.Value != c2.Domain.Value {
		return fmt.Errorf("Domain is invalid '%s' (expected) vs '%s'", c2.Domain.Value, c1.Domain.Value)
	}

	if c1.Port.Value != c2.Port.Value {
		return fmt.Errorf("Port is invalid '%s' (expected) vs '%s'", c2.Port.Value, c1.Port.Value)
	}

	if c1.Path.Value != c2.Path.Value {
		return fmt.Errorf("Path is invalid '%s' (expected) vs '%s'", c2.Path.Value, c1.Path.Value)
	}

	if c1.RegisterProcedure.Value != c2.RegisterProcedure.Value {
		return fmt.Errorf("RegisterProcedure is invalid '%s' (expected) vs '%s'", c2.RegisterProcedure.Value, c1.RegisterProcedure.Value)
	}

	if c1.Protocol.Value != c2.Protocol.Value {
		return fmt.Errorf("Protocol is invalid '%s' (expected) vs '%s'", c2.Protocol.Value, c1.Protocol.Value)
	}

	return nil
}

func testCloudConstructor() helper.Visitor {
	return NewCloud()
}

func _TestCloudToTestVisitor(t testCloud) helper.TestVisitor {
	testVisitor := helper.TestVisitor{
		XML:                t.XML,
		ExpectedError:      nil,
		ExpectedVisitor:    t.ExpectedCloud,
		VisitorConstructor: testCloudConstructor,
		Validator:          testCloudValidator,
	}

	if t.ExpectedError != nil {
		testVisitor.ExpectedError = t.ExpectedError
	}

	return testVisitor
}
func TestCloudBasic(t *testing.T) {

	var testdata = []testCloud{
		{`<cloud domain="rpc.sys.com" port="80" path="/RPC2" registerProcedure="myCloud.rssPleaseNotify" protocol="xml-rpc" />`,
			nil,
			NewTestCloud("rpc.sys.com", "80", "/RPC2", "myCloud.rssPleaseNotify", "xml-rpc"),
		},
		{`<cloud port="80" path="/RPC2" registerProcedure="myCloud.rssPleaseNotify" protocol="xml-rpc" />`,
			helper.NewError(MissingAttribute, ""),
			NewTestCloud("", "80", "/RPC2", "myCloud.rssPleaseNotify", "xml-rpc"),
		},
		{`<cloud domain="rpc.sys.com" path="/RPC2" registerProcedure="myCloud.rssPleaseNotify" protocol="xml-rpc" />`,
			helper.NewError(MissingAttribute, ""),
			NewTestCloud("rpc.sys.com", "", "/RPC2", "myCloud.rssPleaseNotify", "xml-rpc"),
		},
		{`<cloud domain="rpc.sys.com" port="80" registerProcedure="myCloud.rssPleaseNotify" protocol="xml-rpc" />`,
			helper.NewError(MissingAttribute, ""),
			NewTestCloud("rpc.sys.com", "80", "", "myCloud.rssPleaseNotify", "xml-rpc"),
		},
		{`<cloud domain="rpc.sys.com" port="80" path="/RPC2" protocol="xml-rpc" />`,
			helper.NewError(MissingAttribute, ""),
			NewTestCloud("rpc.sys.com", "80", "/RPC2", "", "xml-rpc"),
		},
		{`<cloud domain="rpc.sys.com" port="80" path="/RPC2" registerProcedure="myCloud.rssPleaseNotify" />`,
			helper.NewError(MissingAttribute, ""),
			NewTestCloud("rpc.sys.com", "80", "/RPC2", "myCloud.rssPleaseNotify", ""),
		},
	}

	nbErrors := 0
	len := len(testdata)
	for _, testcloud := range testdata {
		testcase := _TestCloudToTestVisitor(testcloud)

		if err := testcase.CheckTestCase(); err != nil {
			t.Errorf("FAIL\n%s\nXML:\n %s\n", err, testcase.XML)
			nbErrors++
		}
	}

	t.Logf("PASS RATIO = %v/%v\n", len-nbErrors, len)
}
