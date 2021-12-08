package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

func getHostIP() {
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	var ip net.IP
	for _, i := range ifaces {

		if i.Name == INTERFACE_NAME {
			addrs, err := i.Addrs()
			if err != nil {
				panic(err)
			}

			for _, addr := range addrs {
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				}

				if strings.Contains(ip.String(), ".") { // if ipv4
					availability[ip.String()] = true
					hostIP = ip.String()
				}

			}
		}
	}
}

func getBoardcastAddr() {
	addrInts := strings.Split(hostIP, ".")
	addrInts[3] = "255"
	broadcastIP = strings.Join(addrInts, ".")
}

func enumerate3Chars() []string {
	var result []string
	var asciiLst []int
	for i := 65; i < 91; i++ {
		asciiLst = append(asciiLst, i)
	}
	for i := 97; i < 123; i++ {
		asciiLst = append(asciiLst, i)
	}

	for _, firstAscii := range asciiLst {
		for _, secondAscii := range asciiLst {
			for _, thirdAscii := range asciiLst {
				result = append(result, string(firstAscii)+string(secondAscii)+string(thirdAscii))
			}
		}
	}

	return result
}

func startServerTimer(server string, task MapTaskReq) {
	time.Sleep(5 * time.Second)
	serverTasksMutex.Lock()
	if serverTasks[server] == task {
		log.Printf("Server %v timeout", server)
		delete(availability, server)
		mapTasksMutex.Lock()
		mapTasks = append([]MapTaskReq{task}, mapTasks...)
		mapTasksMutex.Unlock()
	}
	serverTasksMutex.Unlock()
}

func dispatchMapTasks() {
	for server, available := range availability {
		if len(mapTasks) == 0 {
			return
		}

		if !available {
			continue
		}

		mapTasksMutex.Lock()
		for len(mapTasks) != 0 {
			t := mapTasks[0]
			mapTasks = mapTasks[1:]
			if result, ok := crackTasks[t.UUID]; ok && result == "" {
				buf, err := json.Marshal(&t)
				if err != nil {
					log.Fatalf("%v", err)
					continue
				}
				reader := bytes.NewReader(buf)
				http.Post("http://"+server+PORT+"/mapTask", "application/json", reader)
				availability[server] = false
				serverTasksMutex.Lock()
				serverTasks[server] = t
				serverTasksMutex.Unlock()
				go startServerTimer(server, t)
				break
			}
		}
		mapTasksMutex.Unlock()

	}
}
