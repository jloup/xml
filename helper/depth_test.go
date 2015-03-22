package helper

import (
	"testing"
)

func depthStateToString(lvl DepthLevelState) string {
	switch lvl {
	case MaxDepthReached:
		return "MaxDepthReached"
	case RootLevel:
		return "RootLevel"
	case ParentLevel:
		return "ParentLevel"
	case SubRootLevel:
		return "SubRootLevel"
	}
	return "unknown"
}

func TestBasic(t *testing.T) {
	d := NewDepthWatcher()

	d.SetMaxDepth(1)

	t.Log(depthStateToString(d.Down()))
	t.Log(depthStateToString(d.Up()))
	t.Log(depthStateToString(d.Down()))
	t.Log(depthStateToString(d.Down()))
	t.Log(depthStateToString(d.Up()))
	t.Log(depthStateToString(d.Up()))
	t.Log(depthStateToString(d.Up()))

}
