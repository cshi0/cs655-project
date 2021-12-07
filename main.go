package main

import (
	"log"
	"sync"

	"github.com/gin-gonic/gin"
)

const INTERFACE_NAME = "eth0"
const PORT = ":8080"

var availability map[string]bool = make(map[string]bool)
var hostIP string
var boardcastIP string

var crackTasks map[string]string = make(map[string]string)

var mapTasksMutex sync.Mutex
var mapTasks []MapTaskReq = make([]MapTaskReq, 50)

var serverTasksMutex sync.Mutex
var serverTasks map[string]MapTaskReq = make(map[string]MapTaskReq)

var enumerated3CharsList []string = make([]string, 0, 157464)

func main() {
	getHostIP()
	getBoardcastAddr()

	log.Printf("%v", availability)
	log.Printf("%v", boardcastIP)

	enumerated3CharsList = enumerate3Chars()

	router := gin.Default()

	router.POST("/newServer", handleNewServer)
	router.POST("/mapTask", handleMapTask)
	router.POST("/taskResult", handleTaskResult)
	router.POST("/crackTask", handleCrackTask)

	router.Run(":8080")

	go boardcastHost()
}