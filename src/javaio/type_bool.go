package javaio

import "bufio"

func ReadBool(stream *bufio.Reader) (result bool, err error) {
	byte_, err := stream.ReadByte()

	if err != nil {
		err = MalformedPacketError { "Bool ended abruptly" }
	}

	if (byte_ == 0x00) {
		result = false
	} else {
		result = true
	}
	return
}

func WriteBool(value bool, stream *bufio.Writer) {
	if (value == true) {
		WriteUByte(0x01, stream)
	} else {
		WriteUByte(0x00, stream)
	}
}
