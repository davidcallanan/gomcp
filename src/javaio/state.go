package javaio

type ClientState struct {
	Protocol uint
	State State
}

var InitialClientState = ClientState {
	Protocol: 0,
	State: StateDeterminingProtocol,
}
