package javaio

type State int
const (
	StateInvalid = iota
	StateDeterminingProtocol = iota
	StateHandshaking = iota
	StateStatus = iota
	StateLogin = iota
	StatePlay = iota
	StatePreNetty = iota
)

// type Protocol int
// const (
// 	ProtocolNetty = iota
// 	ProtocolPreNetty = iota
// )

type Gamemode int
const (
	GamemodeInvalid = iota
	GamemodeSurvival = iota
	GamemodeCreative = iota
	GamemodeAdventure = iota
	GamemodeSpectator = iota
)

type Dimension int
const (
	DimensionInvalid = iota
	DimensionOverworld = iota
	DimensionNether = iota
	DimensionEnd = iota
)
