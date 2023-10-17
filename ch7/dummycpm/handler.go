package main

import (
	"bufio"
	"errors"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"math/rand"
	"net"
	"sync"

	"github.com/alexm77/go_present/ch7/proto"
)

type ProtobufListener struct {
	channelMux *sync.Mutex
}

func NewProtobufListener() *ProtobufListener {
	return &ProtobufListener{&sync.Mutex{}}
}

func (p ProtobufListener) Handle(conn net.Conn) {
	defer conn.Close()
	var (
		buf = make([]byte, 1024)
		r   = bufio.NewReader(conn)
		w   = bufio.NewWriter(conn)
	)

	for {
		p.channelMux.Lock()
		message, err := readMessage(r, buf)
		//log.Println("Got message:", message)

		switch err {
		case nil:
			response, err := createResponse(message)

			if err == nil {
				writeResponse(w, response)
			} else {
				log.Fatalln("Failed to write response:", err)
			}
		case io.EOF:
			p.channelMux.Unlock()
			log.Println("Got EOF")
			return
		default:
			log.Fatalln("Failed to parse message:", err)
		}
		p.channelMux.Unlock()
	}
}

func readMessage(r io.Reader, buf []byte) (*cpmproto.ChargingCommand, error) {
	n, err := r.Read(buf)
	//log.Printf("Read %d bytes\n", n)
	//log.Println("Received", buf[:n-1])
	switch err {
	case nil:
		//data := string(buf[:n])

		message := &cpmproto.ChargingCommand{}
		err = proto.Unmarshal(buf[:n], message)
		//err = proto.Unmarshal(buf[:n-1], message)

		return message, err
	default:
		return nil, err
	}
}

func writeResponse(w *bufio.Writer, response *cpmproto.ChargingSessionEvent) {
	out, err := proto.Marshal(response)
	if err == nil {
		w.Write(out)
		w.Flush()
	} else {
		log.Fatalln("Failed to encode response:", err)
	}
}

func createResponse(message *cpmproto.ChargingCommand) (*cpmproto.ChargingSessionEvent, error) {
	var response cpmproto.ChargingSessionEvent
	switch message.GetType() {
	case cpmproto.Type_Start:
		response = cpmproto.ChargingSessionEvent{
			Type:                cpmproto.Type_Start,
			WallboxSerialNumber: message.GetWallboxSerialNumber(),
			MeterReading:        1000 + uint64(rand.Intn(500)),
			StartDate:           1000000 + uint64(rand.Intn(100000)),
		}
	case cpmproto.Type_Stop:
		response = cpmproto.ChargingSessionEvent{
			Type:                cpmproto.Type_Stop,
			WallboxSerialNumber: message.GetWallboxSerialNumber(),
			MeterReading:        2000 + uint64(rand.Intn(500)),
			StartDate:           1200000 + uint64(rand.Intn(100000)),
		}
	default:
		errDesc := fmt.Sprintf("Unsupported message type: %d", message.GetType())
		return nil, errors.New(errDesc)
	}

	return &response, nil
}
