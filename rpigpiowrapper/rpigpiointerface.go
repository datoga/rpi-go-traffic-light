package rpigpiowrapper

type RPIGPIO interface {
	Prepare() error
	SetState(color string) error
	Destroy()
}
