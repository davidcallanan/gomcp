package javaio

import "bufio"

/**  Clientbound parsing will not be available for the foreseeable future unless contributed by others.  **/

// Serverbound

func ParseLoginStart(data *bufio.Reader) (result LoginStart, err error) {
	username, err := ReadString(data, 16)

	if err != nil {
		return
	}

	result = LoginStart {
		ClientsideUsername: username,
	}
	return
}
