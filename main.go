package main

import (
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const INTERFACE_NAME = "eth0"
const PORT = ":8080"

var availability map[string]bool = make(map[string]bool)
var hostIP string
var hostMask string

var crackTasks map[string]string = make(map[string]string)

var mapTasksMutex sync.Mutex
var mapTasks []MapTaskReq = make([]MapTaskReq, 0, 50)

var serverTasksMutex sync.Mutex
var serverTasks map[string]MapTaskReq = make(map[string]MapTaskReq)

var enumerated3CharsList []string = make([]string, 0, 157464)

var taskTimeMutex sync.Mutex
var avgLatency = 0.
var numTaskDone = 0.
var taskStartingTime map[string]time.Time = make(map[string]time.Time)

var startTime time.Time
var avgThrouput = 0.

func main() {
	getHostIP()
	getBroadcastAddr()

	log.Printf("%v", availability)
	log.Printf("%v", broadcastIP)

	go receiveBroadcast()
	go broadcastHost()

	enumerated3CharsList = enumerate3Chars()

	router := gin.Default()

	router.POST("/newServer", handleNewServer)
	router.POST("/mapTask", handleMapTask)
	router.POST("/taskResult", handleTaskResult)
	router.POST("/crackTask", handleCrackTask)

	startTime = time.Now()

	router.Run(":8080")
}
