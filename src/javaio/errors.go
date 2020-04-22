package javaio

import "fmt"

type MalformedPacketError struct {
	details string
}

func (err *MalformedPacketError) Error() string {
	return fmt.Sprintf("Malformed packet: %s", err.details)
}

type UnsupportedPayloadError struct {
	details string
}

func (err *UnsupportedPayloadError) Error() string {
	return fmt.Sprintf("Unsupported payload: %s", err.details)	
}
