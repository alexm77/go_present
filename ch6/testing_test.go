package main

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestMul(t *testing.T) {
	var tests = []struct {
		c1r, c1i, c2r, c2i float32
		c3r, c3i           float32
	}{
		{1, -1, -1, 1, 0, 2},
		{-1, 1, 1, -1, 0, 2},
	}

	for _, test := range tests {
		tName := fmt.Sprintf("(%v,%vi)*(%v,%vi)", test.c1r, test.c1i, test.c2r, test.c2i)
		t.Run(tName, func(t *testing.T) {
			c1 := ComplexNo{test.c1r, test.c1i}
			c2 := ComplexNo{test.c2r, test.c2i}
			c3 := c1.mul(c2)

			if c3.real != test.c3r {
				t.Errorf("Real is not %v: %v", test.c3r, c3.real)
			}
			if c3.imaginary != test.c3i {
				t.Errorf("Imaginary is not %v: %v", test.c3i, c3.imaginary)
			}
		})
	}
}

func BenchmarkMul(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ComplexNo{float32(i) * rand.Float32(), float32(i) * rand.Float32()}.mul(ComplexNo{float32(i) * rand.Float32(), float32(i) * rand.Float32()})
	}
}

func FuzzMul(f *testing.F) {
	var tests = []struct {
		c1r, c1i, c2r, c2i float32
	}{
		{1, -1, -1, 1},
		{-1, 1, 1, -1},
		{10, 10, 10, 10},
		{-100, -100, -100, -100},
	}

	for _, test := range tests {
		f.Add(test.c1r, test.c1i, test.c2r, test.c2i)
	}
	f.Fuzz(func(t *testing.T, c1r float32, c1i float32, c2r float32, c2i float32) {
		c1 := ComplexNo{c1r, c1i}
		c2 := ComplexNo{c2r, c2i}
		c3 := c1.mul(c2)
		c4 := c2.mul(c1)

		if c3.real != c4.real || c3.imaginary != c4.imaginary {
			t.Errorf("Multiplied numbers don't match: (%v, %vi) != {%v. %vi)", c3.real, c3.imaginary, c4.real, c4.imaginary)
		}
	})
}

//go test -v -fuzz=Fuzz ch6/*
