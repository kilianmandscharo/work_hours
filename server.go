package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func newRouter(db *DB) *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	h := newRequestHandler(db)

	r.POST("/block", h.handleAddBlock)
	r.PUT("/block", h.handleUpdateBlock)
	r.DELETE("/block/:id", h.handleDeleteBlock)
	r.GET("/block/:id", h.handleGetBlockByID)
	r.GET("/block", h.handleGetAllBlocks)
	r.POST("/pause", h.handleAddPause)
	r.PUT("/pause", h.handleUpdatePause)
	r.DELETE("/pause/:id", h.handleDeletePause)

	return r
}
