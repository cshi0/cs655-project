package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// new ip address
type NewServerReq struct {
	Addr string
}

// add new servers to IPList
func handleNewServer(c *gin.Context) {
	req := new(NewServerReq)
	c.BindJSON(req)
	availability[req.Addr] = true
	log.Printf("%v", availability)
	c.Status(http.StatusOK)
}

type MapTaskReq struct {
	MasterAddr string
	ToUnhash   string
	Prefix     string
	UUID       string
}

func handleMapTask(c *gin.Context) {
	req := new(MapTaskReq)
	c.BindJSON(req)
	if req.MasterAddr != "" && req.ToUnhash != "" && req.Prefix != "" && req.UUID != "" {
		go crackPassword(req.Prefix, req.ToUnhash, req.MasterAddr, req.UUID)
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusBadRequest)
	}
}

type TaskResultReq struct {
	IP       string
	UUID     string
	ToUnhash string
	Result   string
	Success  bool
}

func handleTaskResult(c *gin.Context) {
	req := new(TaskResultReq)
	c.BindJSON(req)
	if req.Success {
		crackTasks[req.UUID] = req.Result
	}
	availability[req.IP] = true
	serverTasksMutex.Lock()
	delete(serverTasks, req.IP)
	serverTasksMutex.Unlock()
	dispatchMapTasks()
}

type CrackTaskReq struct {
	ToUnhash string
}

func handleCrackTask(c *gin.Context) {
	req := new(CrackTaskReq)
	c.BindJSON(req)
	uuid, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("%v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	for i := 65; i < 91; i++ {
		mapTasks = append(mapTasks, MapTaskReq{MasterAddr: hostIP, ToUnhash: req.ToUnhash, Prefix: string(i), UUID: uuid.String()})
	}
	for i := 97; i < 123; i++ {
		mapTasks = append(mapTasks, MapTaskReq{MasterAddr: hostIP, ToUnhash: req.ToUnhash, Prefix: string(i), UUID: uuid.String()})
	}
	dispatchMapTasks()

	crackTasks[uuid.String()] = ""
	c.JSON(http.StatusOK, gin.H{
		"UUID": uuid.String(),
	})
}
