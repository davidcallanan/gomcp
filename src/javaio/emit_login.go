package javaio

import "bufio"

/**  Serverbound packet emission is not implemented.  **/

// Clientbound

func EmitLoginSuccess(loginSuccess LoginSuccess, result *bufio.Writer) {
	if len(loginSuccess.Username) > 16 {
		panic("Username of LoginSuccess is too long (must not be over 16 runes)")
	}

	WriteString(loginSuccess.Uuid.String(), result)
	WriteString(loginSuccess.Username, result)
}
