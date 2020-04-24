package javaio

import "bufio"
import "unicode/utf8"

/**  Not all types are implemented at this time.  **/

func ParseVarInt(data *bufio.Reader) (result int32, err error) {
	maxLength := 5
	idx := 0

	for {
		byte_, readErr := data.ReadByte()

		if readErr != nil {
			err = &MalformedPacketError { "VarInt ended abruptly" }
			return
		} else if idx >= maxLength {
			err = &MalformedPacketError { "VarInt exceeded max length "}
			return
		}

		value := byte_ & 0b01111111
		result |= int32(uint(value) << uint(7 * idx))
		idx++

		if byte_ & 0b10000000 == 0 {
			break
		}
	}

	return
}

func ParseLong(data *bufio.Reader) (result int64, err error) {
	const size = 8

	var buf [size]byte
	n, _ := data.Read(buf[:])

	if n < size {
		err = &MalformedPacketError { "Long ended abruptly" }
		return
	}

	for exp := 0; exp < size; exp++ {
		idx := size - exp - 1
		result += int64(buf[idx]) * 1 << exp
	}

	return
}

func ParseUnsignedShort(data *bufio.Reader) (result uint16, err error) {
	const size = 2

	var buf [size]byte
	n, _ := data.Read(buf[:])

	if n < size {
		err = &MalformedPacketError { "Unsigned short ended abruptly" }
		return
	}

	result = uint16(buf[1]) + 256 * uint16(buf[0])
	return
}

func ParseString(data *bufio.Reader, maxRuneCount int) (result string, err error) {
	maxStrLength := maxRuneCount * 4
	strLength, err := ParseVarInt(data)

	if err != nil {
		return
	}

	if int(strLength) > maxStrLength {
		err = &MalformedPacketError { "String exceeded max rune count" } //*
		return
	}

	buf := make([]byte, strLength)
	n, _ := data.Read(buf)
	
	if (n != int(strLength)) {
		err = &MalformedPacketError { "String ended abruptly" }
		return
	}

	result = string(buf)

	if utf8.RuneCountInString(result) > maxRuneCount {
		err = &MalformedPacketError { "String exceeded max rune count" }
		return
	}

	return
}
