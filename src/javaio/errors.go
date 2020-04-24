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

// This is a validation error and should be switched to a panic.
// A separate validate function would be useful.
type WrongStateError struct {
	details string
}

func (err *WrongStateError) Error() string {
	return fmt.Sprintf("Wrong state: %s", err.details)	
}

type PacketTooLargeError struct {
	details string
}

func (err *PacketTooLargeError) Error() string {
	return fmt.Sprintf("Packet too large: %s", err.details)	
}
