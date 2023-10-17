package main

import "fmt"

type movable2d interface {
	move2d(x int, y int)
	getx() int
	gety() int
}

type point2d struct {
	x int
	y int
}

type point3d struct {
	p2d movable2d
	z   int
}

func (p *point2d) move2d(x int, y int) {
	p.x = p.x + x
	p.y = p.y + y
}

func (p *point2d) getx() int {
	return p.x
}

func (p *point2d) gety() int {
	return p.y
}

func (p point3d) move3d(x int, y int, z int) {
	p.p2d.move2d(x, y)
	p.z = p.z + z
}

func (p point3d) String() string {
	return fmt.Sprintf("{%v %v}", p.p2d, p.z)
}

func main() {
	p1 := point2d{}
	p1.move2d(5, -5)
	fmt.Println(p1)

	p2 := point3d{&point2d{10, 20}, 30}
	p2.move3d(-15, -20, 20)
	fmt.Println(p2)

	p3 := point3d{&point2d{y: 20}, 30}
	p3.move3d(-15, -20, 20)
	fmt.Println(p3)
}
