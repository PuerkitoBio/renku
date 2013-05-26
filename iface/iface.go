package iface

type FileEvent interface {
	IsCreate() bool
	IsDelete() bool
	IsModify() bool
	IsRename() bool
	Name() string
}

type Watcher interface {
	Start()
	Stop()
	Dir() string
	Event() <-chan FileEvent
}
