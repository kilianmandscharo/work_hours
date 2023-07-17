package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/kilianmandscharo/work_hours/auth"
	"github.com/kilianmandscharo/work_hours/database"
	"github.com/kilianmandscharo/work_hours/models"
	"github.com/kilianmandscharo/work_hours/utils"
)

type RequestHandler struct {
	db *database.DB
}

func newRequestHandler(db *database.DB) RequestHandler {
	return RequestHandler{db: db}
}

func (r *RequestHandler) handleAddBlock(c *gin.Context) {
	var block models.BlockCreate
	if err := c.BindJSON(&block); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read body"})
		return
	}

	if !block.Valid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid datetime found"})
		return
	}

	if newBlock, err := r.db.AddBlock(block); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not add block"})
	} else {
		c.JSON(http.StatusOK, newBlock)
	}
}

func (r *RequestHandler) handleUpdateBlock(c *gin.Context) {
	var block models.Block
	if err := c.BindJSON(&block); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read body"})
		return
	}

	if !block.Valid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid datetime found"})
		return
	}

	if rowsAffected, err := r.db.UpdateBlock(block); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update block"})
	} else {
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "block not found"})
		} else {
			c.Status(http.StatusOK)
		}
	}
}

func (r *RequestHandler) handleUpdateBlockStart(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read query parameter"})
		return
	}

	var body models.BodyStart
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read body"})
		return
	}

	if !body.Valid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid datetime found"})
		return
	}

	if rowsAffected, err := r.db.UpdateBlockStart(id, body.Start); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update block"})
	} else {
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "block not found"})
		} else {
			c.Status(http.StatusOK)
		}
	}
}

func (r *RequestHandler) handleUpdateBlockEnd(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read query parameter"})
		return
	}

	var body models.BodyEnd
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read body"})
		return
	}

	if !body.Valid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid datetime found"})
		return
	}

	if rowsAffected, err := r.db.UpdateBlockEnd(id, body.End); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update block"})
	} else {
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "block not found"})
		} else {
			c.Status(http.StatusOK)
		}
	}
}

func (r *RequestHandler) handleUpdateBlockHomeoffice(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read query parameter"})
		return
	}

	var body models.BodyHomeoffice
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read body"})
		return
	}

	if rowsAffected, err := r.db.UpdateBlockHomeoffice(id, body.Homeoffice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update block"})
	} else {
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "block not found"})
		} else {
			c.Status(http.StatusOK)
		}
	}
}

func (r *RequestHandler) handleDeleteBlock(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read query parameter"})
		return
	}

	if rowsAffected, err := r.db.DeleteBlock(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete block"})
	} else {
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "block not found"})
		} else {
			c.Status(http.StatusOK)
		}
	}
}

func (r *RequestHandler) handleGetBlockByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read query parameter"})
		return
	}

	if block, err := r.db.GetBlockByID(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get block"})
	} else {
		c.JSON(http.StatusOK, block)
	}
}

func (r *RequestHandler) handleGetAllBlocks(c *gin.Context) {
	if blocks, err := r.db.GetAllBlocks(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get blocks"})
	} else if len(blocks) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no blocks available"})
	} else {
		c.JSON(http.StatusOK, blocks)
	}
}

func (r *RequestHandler) handleAddPause(c *gin.Context) {
	var pause models.PauseCreate
	if err := c.BindJSON(&pause); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read body"})
	}

	if !pause.Valid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid datetime found"})
		return
	}

	if newPause, err := r.db.AddPause(pause); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not add pause"})
	} else {
		c.JSON(http.StatusOK, newPause)
	}
}

func (r *RequestHandler) handleUpdatePause(c *gin.Context) {
	var pause models.Pause
	if err := c.BindJSON(&pause); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read body"})
		return
	}

	if !pause.Valid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid datetime found"})
		return
	}

	if rowsAffected, err := r.db.UpdatePause(pause); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update pause"})
	} else {
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "pause not found"})
		} else {
			c.Status(http.StatusOK)
		}
	}
}

func (r *RequestHandler) handleUpdatePauseStart(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read query parameter"})
		return
	}

	var body models.BodyStart
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read body"})
		return
	}

	if !body.Valid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid datetime found"})
		return
	}

	if rowsAffected, err := r.db.UpdatePauseStart(id, body.Start); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update pause"})
	} else {
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "pause not found"})
		} else {
			c.Status(http.StatusOK)
		}
	}
}

func (r *RequestHandler) handleUpdatePauseEnd(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read query parameter"})
		return
	}

	var body models.BodyEnd
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read body"})
		return
	}

	if !body.Valid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid datetime found"})
		return
	}

	if rowsAffected, err := r.db.UpdatePauseEnd(id, body.End); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update pause"})
	} else {
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "pause not found"})
		} else {
			c.Status(http.StatusOK)
		}
	}
}

func (r *RequestHandler) handleDeletePause(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read query parameter"})
		return
	}

	if rowsAffected, err := r.db.DeletePause(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete pause"})
	} else {
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "pause not found"})
		} else {
			c.Status(http.StatusOK)
		}
	}
}

func (r *RequestHandler) handleStartBlock(c *gin.Context) {
	homeoffice, err := strconv.ParseBool(c.Query("homeoffice"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read query parameter"})
		return
	}

	if block, err := r.db.StartBlock(homeoffice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not start block"})
	} else {
		c.JSON(http.StatusOK, block)
	}
}

func (r *RequestHandler) handleEndBlock(c *gin.Context) {
	if block, err := r.db.EndBlock(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not end block"})
	} else {
		c.JSON(http.StatusOK, block)
	}
}

func (r *RequestHandler) handleStartPause(c *gin.Context) {
	if pause, err := r.db.StartPause(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not start pause"})
	} else {
		c.JSON(http.StatusOK, pause)
	}
}

func (r *RequestHandler) handleEndPause(c *gin.Context) {
	if pause, err := r.db.EndPause(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete pause"})
	} else {
		c.JSON(http.StatusOK, pause)
	}
}

func (r *RequestHandler) handleGetCurrentBlock(c *gin.Context) {
	if block, err := r.db.GetCurrentBlock(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get current block"})
	} else {
		c.JSON(http.StatusOK, block)
	}
}

func (r *RequestHandler) handleLogin(c *gin.Context) {
	env, err := utils.EnvVariables()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load .env file"})
		return
	}
	envTest, err := utils.EnvTestVariables()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load .env file"})
		return
	}

	var email string
	if gin.Mode() == gin.TestMode {
		email = envTest.Email
	} else {
		email = env.Email
	}

	var hash string
	if gin.Mode() == gin.TestMode {
		hash = envTest.Hash
	} else {
		hash = env.Hash
	}

	var login auth.Login
	if err := c.BindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read body"})
		return
	}

	if login.Email != email {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email"})
		return
	}

	if !auth.ValidatePassword(login.Password, hash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
		return
	}

	token, err := auth.CreateToken(login.Email, env.TokenKey)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.String(http.StatusOK, token)
}

func (r *RequestHandler) handleRefresh(c *gin.Context) {
	tokenString, err := auth.ExtractBearerToken(c.GetHeader("Authorization"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not extract token"})
		return
	}

	env, err := utils.EnvVariables()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load .env file"})
	}

	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(env.TokenKey), nil
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

	newToken, err := auth.CreateToken(env.Email, env.TokenKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create new token"})
		return
	}

	c.String(http.StatusOK, newToken)
}
