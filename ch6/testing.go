package main

import "fmt"

type ComplexNo struct {
	real      float32
	imaginary float32
}

func (c1 ComplexNo) mul(c2 ComplexNo) ComplexNo {
	return ComplexNo{
		c1.real*c2.real - c1.imaginary*c2.imaginary,
		c1.real*c2.imaginary + c1.imaginary*c2.real}
}

func main() {
	fmt.Println(ComplexNo{10, -10}.mul(ComplexNo{-10, 10}))
}
