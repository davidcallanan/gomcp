package javaio

import "bufio"

func WriteBlockPos(pos BlockPosition, stream *bufio.Writer) {
	var encoded int64 = ((int64(pos.X) & 0x3FFFFFF) << 38) | ((int64(pos.Z) & 0x3FFFFFF) << 12) | (int64(pos.Y) & 0xFFF)
	WriteLong(encoded, stream)
}
