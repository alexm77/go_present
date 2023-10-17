package main

import (
	"log"
	"math"
	"time"
)

type bu float64

type Benchmark interface {
	AddDurations(restCallDuration time.Duration, pbCallDuration time.Duration)
	Summarize()
}

type RealBenchmark struct {
	duration [][2]bu
}

func (b *RealBenchmark) AddDurations(restCallDuration time.Duration, pbCallDuration time.Duration) {
	newPair := [2]bu{bu(restCallDuration), bu(pbCallDuration)}
	b.duration = append(b.duration, newPair)
}

func (b *RealBenchmark) Summarize() {
	restAvg, pbAvg, restOnlyAvg := avg(b.duration)
	restStdDev, pbStdDev, restOnlyStdDev := stdDev(b.duration, restAvg, pbAvg, restOnlyAvg)
	log.Printf("REST      call avg time: %v,\tstd dev: %v", toDuration(restAvg), toDuration(restStdDev))
	log.Printf("Protobuf  call avg time: %v,\tstd dev: %v", toDuration(pbAvg), toDuration(pbStdDev))
	log.Printf("REST only call avg time: %v,\tstd dev: %v", toDuration(restOnlyAvg), toDuration(restOnlyStdDev))
	log.Println("Successful calls", len(b.duration))
}

type NoBenchmark struct{}

func (*NoBenchmark) AddDurations(restCallDuration time.Duration, pbCallDuration time.Duration) {}

func (NoBenchmark) Summarize() {}

func toDuration(d bu) time.Duration {
	return time.Duration(d) * time.Nanosecond
}

func avg(vec [][2]bu) (bu, bu, bu) {
	restAvg := bu(0)
	pbAvg := bu(0)
	restOnlyAvg := bu(0)
	for _, v := range vec {
		restAvg += v[0] / bu(len(vec))
		pbAvg += v[1] / bu(len(vec))
		restOnlyAvg += (v[0] - v[1]) / bu(len(vec))
	}

	return restAvg, pbAvg, restOnlyAvg
}

func stdDev(vec [][2]bu, restAvg bu, pbAvg bu, restOnlyAvg bu) (bu, bu, bu) {
	restStdDev := bu(0)
	pbStdDev := bu(0)
	restOnlyStdDev := bu(0)
	for _, v := range vec {
		restStdDev += bu(math.Pow(float64(v[0]-restAvg), 2)) / bu(len(vec)-1)
		pbStdDev += bu(math.Pow(float64(v[1]-restAvg), 2)) / bu(len(vec)-1)
		restOnlyStdDev += (v[0] - v[1]) / bu(len(vec))
	}

	return bu(math.Sqrt(float64(restStdDev))), bu(math.Sqrt(float64(pbStdDev))), bu(math.Sqrt(float64(restOnlyStdDev)))
}
