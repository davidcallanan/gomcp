package javaio

type State int
const (
	StateInvalid = iota
	StateHandshaking = iota
	StateStatus = iota
	StateLogin = iota
	StatePlay = iota
)

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
