package main

import "fmt"

type point3d_v2 struct {
	*point2d
	z int
}

func (p point3d_v2) move3d(x int, y int, z int) {
	p.move2d(x, y)
	p.z = p.z + z
}

func (p point3d_v2) String() string {
	return fmt.Sprintf("{%v %v}", *p.point2d, p.z)
}
