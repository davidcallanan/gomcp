package javaio

func EncodePostNettyVersion(version uint) int32 {
	return int32(version) - 81
}

func DecodePostNettyVersion(version int32) uint {
	return uint(version) + 81
}
