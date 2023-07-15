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
	r.PUT("/block_start/:id", h.handleUpdateBlockStart)
	r.PUT("/block_end/:id", h.handleUpdateBlockEnd)
	r.PUT("/block_homeoffice/:id", h.handleUpdateBlockHomeoffice)
	r.DELETE("/block/:id", h.handleDeleteBlock)
	r.GET("/block/:id", h.handleGetBlockByID)
	r.GET("/block", h.handleGetAllBlocks)
	r.POST("/pause", h.handleAddPause)
	r.PUT("/pause", h.handleUpdatePause)
	r.PUT("/pause_start/:id", h.handleUpdatePauseStart)
	r.PUT("/pause_end/:id", h.handleUpdatePauseEnd)
	r.DELETE("/pause/:id", h.handleDeletePause)

	r.POST("/current_block_start", h.handleStartBlock)
	r.POST("/current_block_end", h.handleEndBlock)
	r.GET("/block_current", h.handleGetCurrentBlock)
	r.POST("/current_pause_start", h.handleStartPause)
	r.POST("/current_pause_end", h.handleEndPause)

	r.POST("/login", h.handleLogin)
	r.POST("/refresh", h.handleRefresh)

	return r
}
