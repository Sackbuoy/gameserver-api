package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sackbuoy/gameserver-api/internal/manager"
)

const (
	servers      = "/servers"
	createServer = servers + "/create"
	updateServer = servers + "/update"
	deleteServer = servers + "/delete"
)

func main() {
	r := gin.Default()

	manager, err := manager.New()
	if err != nil {
		panic(err)
	}

	r.POST(createServer, manager.Create)
	r.GET(servers, manager.Read)
	r.PUT(updateServer, manager.Update)
	r.DELETE(deleteServer, manager.Delete)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
