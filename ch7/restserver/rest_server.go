package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
)

const wallboxPath = "/wallbox/"

func main() {
	benchmarkMode, cpuprofile, memprofile := parseFlags()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	b := setupBenchmark(benchmarkMode)
	handler := ResourceByIDHandler{wallboxPath, &b}
	handler.Init()

	handleSigTerm(handler, b, cpuprofile, memprofile)

	http.Handle(wallboxPath, handler)
	log.Println("Started REST server")
	log.Fatalln(http.ListenAndServe(":8080", nil))
}

func parseFlags() (*bool, *string, *string) {
	benchmarkMode := flag.Bool("benchmark", false, "enable benchmark mode")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile := flag.String("memprofile", "", "write memory profile to this file")
	flag.Parse()
	return benchmarkMode, cpuprofile, memprofile
}

func setupBenchmark(benchmarkMode *bool) Benchmark {
	var data [][2]bu
	var b Benchmark
	if *benchmarkMode {
		log.Println("benchmark mode enabled")
		b = &RealBenchmark{data}
	} else {
		b = &NoBenchmark{}
	}
	return b
}

func handleSigTerm(handler ResourceByIDHandler, b Benchmark, cpuprofile *string, memprofile *string) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		handler.Shutdown()
		log.Println("Stopped REST server")
		b.Summarize()
		if *cpuprofile != "" {
			pprof.StopCPUProfile()
		}
		if *memprofile != "" {
			f, err := os.Create(*memprofile)
			if err != nil {
				log.Fatal(err)
			}
			pprof.WriteHeapProfile(f)
			f.Close()
		}
		os.Exit(0)
	}()
}
