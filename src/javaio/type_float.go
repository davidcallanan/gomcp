package javaio

import "math"
import "bufio"

func WriteFloat(value float32, stream *bufio.Writer) {
	n := math.Float32bits(value)
	stream.Write([]byte {
		byte(n >> 24),
		byte(n >> 16),
		byte(n >> 8),
		byte(n),
	})
}
