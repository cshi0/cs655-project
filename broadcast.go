package main

import (
	"log"
	"net"
	"time"
)

const BROADCAST_SEND_PORT = ":8081"
const BROADCAST_RECEIVE_PORT = ":8082"

var broadcastIP string

func broadcastHost() {
	pc, err := net.ListenPacket("udp4", BROADCAST_SEND_PORT)
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	addr, err := net.ResolveUDPAddr("udp4", broadcastIP+BROADCAST_RECEIVE_PORT)
	if err != nil {
		panic(err)
	}

	for {
		_, err = pc.WriteTo([]byte(hostIP), addr)
		if err != nil {
			panic(err)
		}
		time.Sleep(5 * time.Second)
	}
}

func receiveBroadcast() {
	pc, err := net.ListenPacket("udp4", BROADCAST_RECEIVE_PORT)
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	for {
		buf := make([]byte, 1024)
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			panic(err)
		}

		log.Printf("%s sent this: %s\n", addr, buf[:n])
		if _, ok := availability[string(buf[:n])]; !ok {
			availability[string(buf[:n])] = true
		}
		log.Printf("Availability updated: %v", availability)
	}
}
