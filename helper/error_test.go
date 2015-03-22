package helper

import (
	"testing"
)

func TestBasicError(t *testing.T) {
	u := NewErrorChecker(DisableAllError)

	s := Error{flag: XMLTokenError, msg: "YO"}

	u.EnableErrorChecking("ALL", XMLError)

	t.Log(u.CheckFlag("ALL", &s))
	t.Log(u.CheckFlag("ALL", &s))

	u.DisableErrorChecking("content", XMLError)

	t.Log(u.CheckFlag("content", &s))
	t.Log(u.CheckFlag("content", &s))

}
