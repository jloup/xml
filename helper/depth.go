package helper

type DepthWatcher struct {
	Level    int
	MaxDepth int
}

type DepthLevelState uint32

const (
	MaxDepthReached DepthLevelState = iota
	RootLevel
	ParentLevel
	SubRootLevel
)

func NewDepthWatcher() DepthWatcher {
	return DepthWatcher{Level: 0, MaxDepth: -1}
}

func (d *DepthWatcher) SetMaxDepth(depth int) {
	d.MaxDepth = depth
}

func (d *DepthWatcher) IsRoot() bool {
	return d.Level == 0
}

func (d *DepthWatcher) Up() DepthLevelState {
	d.Level--
	if d.Level == 0 {
		return RootLevel
	}
	if d.Level < 0 {
		return ParentLevel
	}
	return SubRootLevel
}

func (d *DepthWatcher) Down() DepthLevelState {
	if d.MaxDepth != -1 && d.MaxDepth == d.Level {
		d.Level++
		return MaxDepthReached
	}

	d.Level++
	return SubRootLevel
}

func (d *DepthWatcher) Reset() {
	d.Level = 0
}
