package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	pb "google.golang.org/protobuf/proto"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/alexm77/go_present/ch7/proto"
)

var conn net.Conn

type ResourceByIDHandler struct {
	resourceName string
	b            *Benchmark
}

func (h ResourceByIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, _ := io.ReadAll(r.Body)

	restStart := time.Now()
	runes := []rune(r.URL.Path)
	id := string(runes[len(h.resourceName):])
	var response string
	var duration time.Duration = -1
	switch r.Method {
	case "PUT":
		response, duration = put(id, string(body))
	default:
		response = "LOL@" + r.Method
	}
	fmt.Fprintf(w, response+"\n")
	restEnd := time.Now()
	restDuration := restEnd.Sub(restStart)
	if duration >= 0 {
		bench := *h.b
		bench.AddDurations(restDuration, duration)
	}
}

func (ResourceByIDHandler) Init() {
	c, err := net.Dial("unix", "/tmp/cpm.sock")
	if err != nil {
		log.Fatalln("Oopsy!", err)
	}
	conn = c
}

func (ResourceByIDHandler) Shutdown() {
	conn.Close()
}

func put(id string, oper string) (string, time.Duration) {
	var (
		r = bufio.NewReader(conn)
		w = bufio.NewWriter(conn)
	)

	pbStart := time.Now()
	message := createCmd(id, oper)
	writeMessage(w, message)

	//go readResponse(conn, buf)

	response, err := readResponse(r)
	pbEnd := time.Now()

	pbDuration := pbEnd.Sub(pbStart)

	if err == nil {
		tmp, jsonErr := json.Marshal(response)
		if jsonErr == nil {
			return string(tmp[:]), pbDuration
		}
		return "{}", pbDuration
	}
	return "{}", -1
}

func writeMessage(w *bufio.Writer, message *cpmproto.ChargingCommand) {
	out, err := pb.Marshal(message)
	if err == nil {
		//log.Println("Sending ", out)
		w.Write(out)
		w.Flush()
	} else {
		log.Fatalln("Failed to encode message:", err)
	}
}

var pbReadBuf = make([]byte, 1024)

func readResponse(r io.Reader) (*cpmproto.ChargingSessionEvent, error) {
	n, err := r.Read(pbReadBuf)
	if err != nil {
		log.Fatalln("Failed to parse response:", err)
	}
	//data := string(buf[:n])

	response := &cpmproto.ChargingSessionEvent{}
	err = pb.Unmarshal(pbReadBuf[:n], response)
	//log.Println("Got response:", response)

	return response, err
}

func createCmd(id string, oper string) *cpmproto.ChargingCommand {
	var opType cpmproto.Type

	if oper == "stop" {
		opType = cpmproto.Type_Stop
	} else {
		opType = cpmproto.Type_Start
	}

	message := cpmproto.ChargingCommand{
		Type:                opType,
		WallboxSerialNumber: id,
	}

	return &message
}
