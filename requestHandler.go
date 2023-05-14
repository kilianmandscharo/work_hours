package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RequestHandler struct {
	db *DB
}

func newRequestHandler(db *DB) RequestHandler {
	return RequestHandler{db: db}
}

func (r *RequestHandler) handleAddBlock(c *gin.Context) {
	var block BlockCreate
	if err := c.BindJSON(&block); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read body"})
		return
	}
	if newBlock, err := r.db.addBlock(block); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not add block"})
	} else {
		c.JSON(http.StatusOK, newBlock)
	}
}

func (r *RequestHandler) handleUpdateBlock(c *gin.Context) {
	var block Block
	if err := c.BindJSON(&block); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read body"})
		return
	}
	if err := r.db.updateBlock(block); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update block"})
	} else {
		c.Status(http.StatusOK)
	}
}

func (r *RequestHandler) handleDeleteBlock(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read query parameter"})
		return
	}
	if err := r.db.deleteBlock(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete block"})
	} else {
		c.Status(http.StatusOK)
	}
}

func (r *RequestHandler) handleGetBlockByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read query parameter"})
		return
	}
	if block, err := r.db.getBlockByID(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get block"})
	} else {
		c.JSON(http.StatusOK, block)
	}
}

func (r *RequestHandler) handleGetAllBlocks(c *gin.Context) {
	if blocks, err := r.db.getAllBlocks(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get blocks"})
	} else {
		c.JSON(http.StatusOK, blocks)
	}
}

func (r *RequestHandler) handleAddPause(c *gin.Context) {
	var pause PauseCreate
	if err := c.BindJSON(&pause); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read body"})
	}
	if newPause, err := r.db.addPause(pause); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not add pause"})
	} else {
		c.JSON(http.StatusOK, newPause)
	}
}

func (r *RequestHandler) handleUpdatePause(c *gin.Context) {
	var pause Pause
	if err := c.BindJSON(&pause); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read body"})
		return
	}
	if err := r.db.updatePause(pause); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update pause"})
	} else {
		c.Status(http.StatusOK)
	}
}

func (r *RequestHandler) handleDeletePause(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read query parameter"})
		return
	}
	if err := r.db.deletePause(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete pause"})
	} else {
		c.Status(http.StatusOK)
	}
}
