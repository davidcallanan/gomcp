package javaio

import "math"
import "bufio"

func WriteDouble(value float64, stream *bufio.Writer) {
	n := math.Float64bits(value)
	stream.Write([]byte {
		byte(n >> 56),
		byte(n >> 48),
		byte(n >> 40),
		byte(n >> 32),
		byte(n >> 24),
		byte(n >> 16),
		byte(n >> 8),
		byte(n),
	})
}
