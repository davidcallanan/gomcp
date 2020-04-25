package javaio

type BlockPosition struct {
	X int
	Y int
	Z int
}

func ConstrainBlockPosition(pos BlockPosition) BlockPosition {
	// This function isn't really neccessary
	
	x := pos.X
	if x > 33554431 {
		x = 33554431
	} else if x < -33554432 {
		x = -33554432
	}

	y := pos.Y
	if y > 2047 {
		y = 2047
	} else if y < -2048 {
		y = -2048
	}

	z := pos.Z
	if z > 33554431 {
		z = 33554431
	} else if z < -33554432 {
		z = -33554432
	}

	return BlockPosition {
		X: x,
		Y: y,
		Z: z,
	}
}
