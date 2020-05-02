package javaio

type ClientContext struct {
	Protocol uint
	State State
}

var InitialClientContext = ClientContext {
	Protocol: 0,
	State: StateDeterminingProtocol,
}
