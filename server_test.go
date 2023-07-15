package main

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var token string

func init() {
	env, err := envVariables()
	if err != nil {
		log.Fatal("could not load env file")
	}
	token, err = createToken(env.email, env.tokenKey)
	if err != nil {
		log.Fatal("could not create token")
	}
}

func TestAddBlockRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	r := newRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("no body", func(t *testing.T) {
		assertRequest(
			t,
			r,
			token,
			http.MethodPost,
			"/block",
			http.StatusBadRequest)
	})

	t.Run("invalid body", func(t *testing.T) {
		assertRequestWithBody(
			t,
			r,
			token,
			http.MethodPost,
			"/block",
			struct{ Invalid string }{Invalid: "test"},
			http.StatusBadRequest)
	})

	t.Run("valid body", func(t *testing.T) {
		assertRequestWithBody(
			t,
			r,
			token,
			http.MethodPost,
			"/block",
			testBlockCreate(),
			http.StatusOK)
	})
}

func TestUpdateBlockRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	r := newRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("no body", func(t *testing.T) {
		assertRequest(
			t,
			r,
			token,
			http.MethodPut,
			"/block",
			http.StatusBadRequest)
	})

	t.Run("not found", func(t *testing.T) {
		block := testBlockUpdated()
		block.Id = 12
		assertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			"/block",
			block,
			http.StatusNotFound)
	})

	t.Run("valid body", func(t *testing.T) {
		db.addBlock(testBlockCreate())
		assertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			"/block",
			testBlockUpdated(),
			http.StatusOK)
	})
}

func TestUpdateBlockStartRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	r := newRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("bad query param", func(t *testing.T) {
		assertRequest(
			t,
			r,
			token,
			http.MethodPut,
			"/block_start/a",
			http.StatusBadRequest)
	})

	t.Run("invalid body", func(t *testing.T) {
		assertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			"/block_start/1",
			BodyEnd{End: "test"},
			http.StatusBadRequest)
	})

	t.Run("not found", func(t *testing.T) {
		assertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			"/block_start/12",
			BodyStart{Start: "test"},
			http.StatusNotFound)
	})

	t.Run("valid body", func(t *testing.T) {
		db.addBlock(testBlockCreate())
		assertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			fmt.Sprintf("/block_start/%d", bID),
			BodyStart{Start: "test"},
			http.StatusOK)
	})
}

func TestUpdateBlockEndRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	r := newRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("bad query param", func(t *testing.T) {
		assertRequest(
			t,
			r,
			token,
			http.MethodPut,
			"/block_end/a",
			http.StatusBadRequest)
	})

	t.Run("invalid body", func(t *testing.T) {
		assertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			"/block_end/1",
			BodyStart{Start: "test"},
			http.StatusBadRequest)
	})

	t.Run("not found", func(t *testing.T) {
		assertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			"/block_end/12",
			BodyEnd{End: "test"},
			http.StatusNotFound)
	})

	t.Run("valid body", func(t *testing.T) {
		db.addBlock(testBlockCreate())
		assertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			fmt.Sprintf("/block_end/%d", bID),
			BodyEnd{End: "test"},
			http.StatusOK)
	})
}

func TestUpdateBlockHomeofficeRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	r := newRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("bad query param", func(t *testing.T) {
		assertRequest(
			t,
			r,
			token,
			http.MethodPut,
			"/block_homeoffice/a",
			http.StatusBadRequest)
	})

	t.Run("invalid body", func(t *testing.T) {
		assertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			"/block_homeoffice/1",
			BodyStart{Start: "test"},
			http.StatusBadRequest)
	})

	t.Run("not found", func(t *testing.T) {
		assertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			"/block_homeoffice/12",
			BodyHomeoffice{Homeoffice: true},
			http.StatusNotFound)
	})

	t.Run("valid body", func(t *testing.T) {
		db.addBlock(testBlockCreate())
		assertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			fmt.Sprintf("/block_homeoffice/%d", bID),
			BodyHomeoffice{Homeoffice: true},
			http.StatusOK)
	})
}

func TestUpdatePauseStartRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	r := newRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("bad query param", func(t *testing.T) {
		assertRequest(
			t,
			r,
			token,
			http.MethodPut,
			"/pause_start/a",
			http.StatusBadRequest)
	})

	t.Run("invalid body", func(t *testing.T) {
		assertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			"/pause_start/1",
			BodyEnd{End: "test"},
			http.StatusBadRequest)
	})

	t.Run("not found", func(t *testing.T) {
		assertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			"/pause_start/12",
			BodyStart{Start: "test"},
			http.StatusNotFound)
	})

	t.Run("valid body", func(t *testing.T) {
		db.addBlock(testBlockCreate())
		assertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			fmt.Sprintf("/pause_start/%d", bID),
			BodyStart{Start: "test"},
			http.StatusOK)
	})
}

func TestUpdatePauseEndRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	r := newRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("bad query param", func(t *testing.T) {
		assertRequest(
			t,
			r,
			token,
			http.MethodPut,
			"/pause_end/a",
			http.StatusBadRequest)
	})

	t.Run("invalid body", func(t *testing.T) {
		assertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			"/pause_end/1",
			BodyStart{Start: "test"},
			http.StatusBadRequest)
	})

	t.Run("not found", func(t *testing.T) {
		assertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			"/pause_end/12",
			BodyEnd{End: "test"},
			http.StatusNotFound)
	})

	t.Run("valid body", func(t *testing.T) {
		db.addBlock(testBlockCreate())
		assertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			fmt.Sprintf("/pause_end/%d", bID),
			BodyEnd{End: "test"},
			http.StatusOK)
	})
}

func TestDeleteBlockRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	r := newRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("invalid query param", func(t *testing.T) {
		assertRequest(
			t,
			r,
			token,
			http.MethodDelete,
			"/block/a",
			http.StatusBadRequest)
	})

	t.Run("block not found", func(t *testing.T) {
		assertRequest(
			t,
			r,
			token,
			http.MethodDelete,
			"/block/12",
			http.StatusNotFound)
	})

	t.Run("valid request", func(t *testing.T) {
		db.addBlock(testBlockCreate())
		assertRequest(
			t,
			r,
			token,
			http.MethodDelete,
			fmt.Sprintf("/block/%d", bID),
			http.StatusOK)
	})
}

func TestGetBlockByIDRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	r := newRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("invalid query param", func(t *testing.T) {
		assertRequest(
			t,
			r,
			token,
			http.MethodGet,
			"/block/a",
			http.StatusBadRequest)
	})

	t.Run("not found", func(t *testing.T) {
		assertRequest(
			t,
			r,
			token,
			http.MethodGet,
			"/block/12",
			http.StatusInternalServerError)
	})

	t.Run("valid request", func(t *testing.T) {
		db.addBlock(testBlockCreate())
		assertRequest(
			t,
			r,
			token,
			http.MethodGet,
			fmt.Sprintf("/block/%d", bID),
			http.StatusOK)
	})
}

func TestGetAllBlocksRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	r := newRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("no blocks available", func(t *testing.T) {
		assertRequest(
			t,
			r,
			token,
			http.MethodGet,
			"/block",
			http.StatusNotFound)
	})

	t.Run("blocks found", func(t *testing.T) {
		db.addBlock(testBlockCreate())
		assertRequest(
			t,
			r,
			token,
			http.MethodGet,
			"/block",
			http.StatusOK)
	})
}

func TestAddPauseRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	r := newRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("no body", func(t *testing.T) {
		assertRequest(
			t,
			r,
			token,
			http.MethodPost,
			"/pause",
			http.StatusBadRequest)
	})

	t.Run("no block with blockID", func(t *testing.T) {
		assertRequestWithBody(
			t,
			r,
			token,
			http.MethodPost,
			"/pause",
			testPauseCreate(),
			http.StatusInternalServerError)
	})

	t.Run("valid body", func(t *testing.T) {
		db.addBlock(testBlockCreateWithoutPause())
		assertRequestWithBody(
			t,
			r,
			token,
			http.MethodPost,
			"/pause",
			testPauseCreate(),
			http.StatusOK)
	})
}

func TestUpdatePauseRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	r := newRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("no body", func(t *testing.T) {
		assertRequest(
			t,
			r,
			token,
			http.MethodPut,
			"/pause",
			http.StatusBadRequest)
	})

	t.Run("not found", func(t *testing.T) {
		assertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			"/pause",
			testPauseUpdated(),
			http.StatusNotFound)
	})

	t.Run("valid body", func(t *testing.T) {
		db.addBlock(testBlockCreate())
		assertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			"/pause",
			testPauseUpdated(),
			http.StatusOK)
	})
}

func TestDeletePauseRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	r := newRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("invalid query param", func(t *testing.T) {
		assertRequest(
			t,
			r,
			token,
			http.MethodDelete,
			"/pause/a",
			http.StatusBadRequest)
	})

	t.Run("not found", func(t *testing.T) {
		assertRequest(
			t,
			r,
			token,
			http.MethodDelete,
			"/pause/12",
			http.StatusNotFound)
	})

	t.Run("valid request", func(t *testing.T) {
		db.addBlock(testBlockCreate())
		assertRequest(
			t,
			r,
			token,
			http.MethodDelete,
			fmt.Sprintf("/pause/%d", pID),
			http.StatusOK)
	})
}

func TestStartBlockRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	r := newRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("invalid query param", func(t *testing.T) {
		assertRequest(
			t,
			r,
			token,
			http.MethodPost,
			"/current_block_start?homeoffice=bad_param",
			http.StatusBadRequest)
	})

	t.Run("valid request", func(t *testing.T) {
		assertRequest(
			t,
			r,
			token,
			http.MethodPost,
			"/current_block_start?homeoffice=false",
			http.StatusOK)
	})

	t.Run("block already active", func(t *testing.T) {
		assertRequest(
			t,
			r,
			token,
			http.MethodPost,
			"/current_block_start?homeoffice=false",
			http.StatusInternalServerError)
	})
}

func TestEndBlockRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	r := newRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("block not started", func(t *testing.T) {
		assertRequest(
			t,
			r,
			token,
			http.MethodPost,
			"/current_block_end",
			http.StatusInternalServerError)
	})

	t.Run("pause still active", func(t *testing.T) {
		db.startBlock(false)
		db.startPause()
		assertRequest(
			t,
			r,
			token,
			http.MethodPost,
			"/current_block_end",
			http.StatusInternalServerError)
	})

	t.Run("valid request", func(t *testing.T) {
		db.endPause()
		assertRequest(
			t,
			r,
			token,
			http.MethodPost,
			"/current_block_end",
			http.StatusOK)
	})
}

func TestStartPauseRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	r := newRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("no block active", func(t *testing.T) {
		assertRequest(
			t,
			r,
			token,
			http.MethodPost,
			"/current_pause_start",
			http.StatusInternalServerError)
	})

	t.Run("valid request", func(t *testing.T) {
		db.startBlock(false)
		assertRequest(
			t,
			r,
			token,
			http.MethodPost,
			"/current_pause_start",
			http.StatusOK)
	})

	t.Run("pause already active", func(t *testing.T) {
		assertRequest(
			t,
			r,
			token,
			http.MethodPost,
			"/current_pause_start",
			http.StatusInternalServerError)
	})
}

func TestEndPauseRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	r := newRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("no block active", func(t *testing.T) {
		assertRequest(
			t,
			r,
			token,
			http.MethodPost,
			"/current_pause_end",
			http.StatusInternalServerError)
	})

	t.Run("no pause active", func(t *testing.T) {
		db.startBlock(false)
		assertRequest(
			t,
			r,
			token,
			http.MethodPost,
			"/current_pause_end",
			http.StatusInternalServerError)
	})

	t.Run("valid request", func(t *testing.T) {
		db.startPause()
		assertRequest(
			t,
			r,
			token,
			http.MethodPost,
			"/current_pause_end",
			http.StatusOK)
	})
}

func TestGetCurrentBlockRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	r := newRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("no block active", func(t *testing.T) {
		assertRequest(
			t,
			r,
			token,
			http.MethodGet,
			"/block_current",
			http.StatusInternalServerError)
	})

	t.Run("valid request", func(t *testing.T) {
		db.startBlock(false)
		assertRequest(
			t,
			r,
			token,
			http.MethodGet,
			"/block_current",
			http.StatusOK)
	})
}

func TestLoginRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	r := newRouter(db)
	gin.SetMode(gin.TestMode)

	envTest, err := envTestVariables()
	assert.NoError(t, err)

	t.Run("no body", func(t *testing.T) {
		assertRequest(
			t,
			r,
			token,
			http.MethodPost,
			"/login",
			http.StatusBadRequest)
	})

	t.Run("invalid email", func(t *testing.T) {
		assertRequestWithBody(
			t,
			r,
			token,
			http.MethodPost,
			"/login",
			Login{Email: "invalid@gmail.com", Password: envTest.password},
			http.StatusUnauthorized)
	})

	t.Run("invalid password", func(t *testing.T) {
		assertRequestWithBody(
			t,
			r,
			token,
			http.MethodPost,
			"/login",
			Login{Email: envTest.email, Password: "987654321"},
			http.StatusUnauthorized)
	})

	t.Run("valid request", func(t *testing.T) {
		assertRequestWithBody(
			t,
			r,
			token,
			http.MethodPost,
			"/login",
			Login{Email: envTest.email, Password: envTest.password},
			http.StatusOK)
	})
}

func TestRefreshRoute(t *testing.T) {
	db := getNewTestDatabase()
	defer db.close()
	r := newRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("token still valid", func(t *testing.T) {
		assertRequest(
			t,
			r,
			token,
			http.MethodPost,
			"/refresh",
			http.StatusBadRequest)
	})
}
