package javaio

import "bufio"

/**  Serverbound packet emission is not implemented.  **/

// Clientbound

func EmitLoginSuccess(loginSuccess LoginSuccess, result *bufio.Writer) (err error) {
	if len(loginSuccess.Username) > 16 {
		err := &PacketTooLargeError { "Username of LoginSuccess is too long (must not be over 16 runes)" }
		panic(err)
	}

	err = EmitString(loginSuccess.Uuid.String(), result)
	if err != nil {
		return
	}

	err = EmitString(loginSuccess.Username, result)
	if err != nil {
		return
	}

	return
}
