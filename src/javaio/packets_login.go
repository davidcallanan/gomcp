package javaio

import "github.com/google/uuid"

// Clientbound

type LoginSuccess struct {
	Uuid uuid.UUID
	Username string
}

// Serverbound

type LoginStart struct {
	ClientsideUsername string
}
