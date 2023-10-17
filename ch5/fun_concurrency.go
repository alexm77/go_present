package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

const matrixSize = 1000
const workers = 100

type matrix [matrixSize][matrixSize]float32

func main() {
	start := time.Now()
	mchan := make(chan *matrix)
	sum := matrix{}
	wgWorkers := sync.WaitGroup{}
	wgWorkers.Add(workers)

	for i := 0; i < workers; i++ {
		go matrixMultiply(mchan, &wgWorkers)
	}
	go matrixSum(&sum, mchan)
	wgWorkers.Wait()

	for i := 0; i < matrixSize; i++ {
		fmt.Println(sum[i])
	}
	fmt.Printf("Done in %v\n", time.Now().Sub(start))
}

func matrixMultiply(mchan chan *matrix, wg *sync.WaitGroup) {
	defer wg.Done()
	m1 := newRandomizedMatrix()
	m2 := newRandomizedMatrix()
	mul := matrix{}

	for i := 0; i < matrixSize; i++ {
		for j := 0; j < matrixSize; j++ {
			var sum float32 = 0
			for k := 0; k < matrixSize; k++ {
				sum += m1[i][k] * m2[k][j]
			}
			mul[i][j] = sum
		}
	}
	mchan <- &mul
}

func matrixSum(msum *matrix, mchan chan *matrix) {
	for m := range mchan {
		for i := 0; i < matrixSize; i++ {
			for j := 0; j < matrixSize; j++ {
				msum[i][j] += m[i][j]
			}
		}
	}
}

func newRandomizedMatrix() matrix {
	m := matrix{}
	for i := 0; i < matrixSize; i++ {
		for j := 0; j < matrixSize; j++ {
			m[i][j] = rand.Float32()
		}
	}
	return m
}

func (m matrix) String() string {
	str := ""
	for i := 0; i < matrixSize; i++ {
		str += "Row " + strconv.Itoa(i+1) + "\t"
		for j := 0; j < matrixSize; j++ {
			str = str + strconv.FormatFloat(float64(m[i][j]), 'f', 2, 32) + "\t"
		}
		str += "\n"
	}
	return str
}
