package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

func crackWorker(id int, prefix string, toUnhash string, quit chan bool) string {
	// log.Printf("worker %v created", id)
	for _, suffix := range enumerated3CharsList {
		select {
		case <-quit:
			// log.Printf("worker %v quited by signal", id)
			return ""
		default:
			str := prefix + suffix
			data := []byte(str)
			hash := md5.Sum(data)
			if hex.EncodeToString(hash[:]) == toUnhash {
				// log.Printf("worker %v founded and quited", id)
				return str
			}
		}
	}
	// log.Printf("worker %v quited", id)
	return ""
}

func crackPassword(prefix string, toUnhash string, masterAddr string, uuid string) {
	quit := make(chan bool)

	var wg sync.WaitGroup
	var result string
	for i := 65; i < 91; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			workerResult := crackWorker(i-64, prefix+string(i), toUnhash, quit)
			if workerResult != "" {
				quit <- true
				result = workerResult
			}
		}(i)
	}
	for i := 97; i < 123; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			workerResult := crackWorker(i-64, prefix+string(i), toUnhash, quit)
			if workerResult != "" {
				quit <- true
				result = workerResult
			}
		}(i)
	}

	wg.Wait()

	log.Printf("RESULT: %v, %v, %v", masterAddr, toUnhash, result)

	req := TaskResultReq{IP: hostIP, UUID: uuid, ToUnhash: toUnhash, Result: result, Success: result != ""}
	buf, err := json.Marshal(&req)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	reader := bytes.NewReader(buf)
	http.Post("http://"+boardcastIP+PORT+"/mapTask", "application/json", reader)
}
