package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const BROADCAST_SEND_PORT = ":8081"
const BROADCAST_RECEIVE_PORT = ":8082"

const MAX_SERVER_NUM = 100

var broadcastIP string

func scanServer() {
	startAddr := []string{"10", "10", "0", "1"}
	for i := 1; i <= MAX_SERVER_NUM; i++ {
		startAddr[2] = strconv.Itoa(i)
		ipStr := strings.Join(startAddr, ".")
		go func() {
			for {
				buf, err := json.Marshal(&NewServerReq{Addr: hostIP})
				if err != nil {
					log.Printf("%v", err)
				}
				reader := bytes.NewReader(buf)
				_, err = http.Post("http://"+ipStr+PORT+"/newServer", "application/json", reader)
				// if err != nil && strings.Contains(err.Error(), "timeout") {
				// 	log.Printf("%v", err)
				// }
				time.Sleep(5 * time.Second)
			}
		}()
	}
}

func broadcastHost() {
	pc, err := net.ListenPacket("udp4", BROADCAST_SEND_PORT)
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	startAddr := []string{"10", "10", "0", "255"}
	broadcastAddrs := make([]*net.UDPAddr, 0, MAX_SERVER_NUM)

	for i := 1; i <= MAX_SERVER_NUM; i++ {
		startAddr[2] = strconv.Itoa(i)
		ipStr := strings.Join(startAddr, ".")
		addr, err := net.ResolveUDPAddr("udp4", ipStr+BROADCAST_RECEIVE_PORT)
		broadcastAddrs = append(broadcastAddrs, addr)
		if err != nil {
			panic(err)
		}
	}

	for {
		for _, addr := range broadcastAddrs {
			_, err = pc.WriteTo([]byte(hostIP), addr)
			if err != nil {
				panic(err)
			}
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
