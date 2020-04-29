package javaio

import "bufio"
import "unicode/utf8"

func ReadString(stream *bufio.Reader, maxRuneCount int) (result string, err error) {
	maxStrLength := maxRuneCount * 4
	strLength, err := ReadVarInt(stream)

	if err != nil {
		return
	}

	if int(strLength) > maxStrLength {
		err = &MalformedPacketError { "String exceeded max rune count" } //*
		return
	}

	buf := make([]byte, strLength)
	n, _ := stream.Read(buf)
	
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

func WriteString(value string, stream *bufio.Writer) {
	// TODO: int32 cast potentially unsafe?
	WriteVarInt(int32(len(value)), stream)
	stream.WriteString(value)
}

func WriteUTF16(value []uint16, stream *bufio.Writer) {
	// Little-endian
	// Needs proper testing
	for char := range value {
		stream.WriteByte(byte(char))
		stream.WriteByte(byte(char >> 8))
	}
}
