package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func newRouter(db *DB) *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())
	r.Use(authorizer())

	h := newRequestHandler(db)

	r.POST("/block", h.handleAddBlock)
	r.PUT("/block", h.handleUpdateBlock)
	r.DELETE("/block/:id", h.handleDeleteBlock)
	r.GET("/block/:id", h.handleGetBlockByID)
	r.GET("/block", h.handleGetAllBlocks)
	r.POST("/pause", h.handleAddPause)
	r.PUT("/pause", h.handleUpdatePause)
	r.DELETE("/pause/:id", h.handleDeletePause)

	r.POST("/block_start", h.handleStartBlock)
	r.POST("/block_end", h.handleEndBlock)
	r.GET("/block_current", h.handleGetCurrentBlock)
	r.POST("/pause_start", h.handleStartPause)
	r.POST("/pause_end", h.handleEndPause)

	r.POST("/login", h.handleLogin)
	r.POST("/refresh", h.handleRefresh)

	return r
}
