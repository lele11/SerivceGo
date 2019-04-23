package game

type Module interface {
	Key() string
	Init()
	Load()
	Save()
}

const (
	ModuleBaseInfo = 1
)

const (
	SourceChangeInit    = 9
	SourceChangeOffline = 10
)
