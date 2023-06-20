package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
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

func (r *RequestHandler) handleStartBlock(c *gin.Context) {
	if block, err := r.db.startBlock(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not start block"})
	} else {
		c.JSON(http.StatusOK, block)
	}
}

func (r *RequestHandler) handleEndBlock(c *gin.Context) {
	if block, err := r.db.endBlock(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not end block"})
	} else {
		c.JSON(http.StatusOK, block)
	}
}

func (r *RequestHandler) handleStartPause(c *gin.Context) {
	if pause, err := r.db.startPause(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not start pause"})
	} else {
		c.JSON(http.StatusOK, pause)
	}
}

func (r *RequestHandler) handleEndPause(c *gin.Context) {
	if pause, err := r.db.endPause(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete pause"})
	} else {
		c.JSON(http.StatusOK, pause)
	}
}

func (r *RequestHandler) handleGetCurrentBlock(c *gin.Context) {
	if block, err := r.db.getCurrentBlock(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get current block"})
	} else {
		c.JSON(http.StatusOK, block)
	}
}

func (r *RequestHandler) handleLogin(c *gin.Context) {
	env, err := envVariables()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load .env file"})
	}
	envTest, err := envTestVariables()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load .env file"})
	}

	var email string
	if gin.Mode() == gin.TestMode {
		email = envTest.email
	} else {
		email = env.email
	}

	var hash string
	if gin.Mode() == gin.TestMode {
		hash = envTest.hash
	} else {
		hash = env.hash
	}

	var login Login
	if err := c.BindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read body"})
		return
	}

	if login.Email != email {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email"})
		return
	}

	if !validatePassword(login.Password, hash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
		return
	}

	c.String(http.StatusOK, "authorized")

	token, err := createToken(login.Email, env.tokenKey)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (r *RequestHandler) handleRefresh(c *gin.Context) {
	tokenString, err := extractBearerToken(c.GetHeader("Authorization"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not extract token"})
		return
	}

	env, err := envVariables()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load .env file"})
	}

	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(env.tokenKey), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not parse token"})
		return
	}
	if !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if time.Until(time.UnixMilli(claims.ExpiresAt)) > 30*time.Second {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token still valid"})
		return
	}

	newToken, err := createToken(env.email, env.tokenKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create new token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": newToken})
}
