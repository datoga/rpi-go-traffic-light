package rpigpiowrapper

type RPIGPIOPhy struct {
	*RPIGPIOMock
}

func NewRPIGPIOPhy() *RPIGPIOPhy {
	return &RPIGPIOPhy{RPIGPIOMock: NewRPIGPIOMock()}
}
