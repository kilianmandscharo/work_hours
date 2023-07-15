package server

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kilianmandscharo/work_hours/auth"
	"github.com/kilianmandscharo/work_hours/database"
	"github.com/kilianmandscharo/work_hours/models"
	"github.com/kilianmandscharo/work_hours/utils"
	"github.com/stretchr/testify/assert"
)

var token string

func init() {
	env, err := utils.EnvVariables()
	if err != nil {
		log.Fatal("could not load env file")
	}
	token, err = auth.CreateToken(env.Email, env.TokenKey)
	if err != nil {
		log.Fatal("could not create token")
	}
}

func TestAddBlockRoute(t *testing.T) {
	db := database.GetNewTestDatabase()
	defer db.Close()
	r := NewRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("no body", func(t *testing.T) {
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodPost,
			"/block",
			http.StatusBadRequest)
	})

	t.Run("invalid body", func(t *testing.T) {
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPost,
			"/block",
			struct{ Invalid string }{Invalid: "test"},
			http.StatusBadRequest)
	})

	t.Run("invalid block start time", func(t *testing.T) {
		block := utils.TestBlockCreate()
		block.Start = "invalid start"
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPost,
			"/block",
			block,
			http.StatusBadRequest)
	})

	t.Run("invalid pause start time", func(t *testing.T) {
		block := utils.TestBlockCreate()
		block.Pauses[0].Start = "invalid start"
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPost,
			"/block",
			block,
			http.StatusBadRequest)
	})

	t.Run("valid body", func(t *testing.T) {
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPost,
			"/block",
			utils.TestBlockCreate(),
			http.StatusOK)
	})
}

func TestUpdateBlockRoute(t *testing.T) {
	db := database.GetNewTestDatabase()
	defer db.Close()
	r := NewRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("no body", func(t *testing.T) {
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodPut,
			"/block",
			http.StatusBadRequest)
	})

	t.Run("invalid body", func(t *testing.T) {
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			"/block",
			struct{ Invalid string }{Invalid: "test"},
			http.StatusBadRequest)
	})

	t.Run("invalid block start time", func(t *testing.T) {
		block := utils.TestBlockCreate()
		block.Start = "invalid start"
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPost,
			"/block",
			block,
			http.StatusBadRequest)
	})

	t.Run("invalid pause start time", func(t *testing.T) {
		block := utils.TestBlockCreate()
		block.Pauses[0].Start = "invalid start"
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPost,
			"/block",
			block,
			http.StatusBadRequest)
	})

	t.Run("not found", func(t *testing.T) {
		block := utils.TestBlockUpdated()
		block.Id = 12
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			"/block",
			block,
			http.StatusNotFound)
	})

	t.Run("valid body", func(t *testing.T) {
		db.AddBlock(utils.TestBlockCreate())
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			"/block",
			utils.TestBlockUpdated(),
			http.StatusOK)
	})
}

func TestUpdateBlockStartRoute(t *testing.T) {
	db := database.GetNewTestDatabase()
	defer db.Close()
	r := NewRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("bad query param", func(t *testing.T) {
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodPut,
			"/block_start/a",
			http.StatusBadRequest)
	})

	t.Run("no body", func(t *testing.T) {
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodPut,
			"/block_start/1",
			http.StatusBadRequest)
	})

	t.Run("invalid body", func(t *testing.T) {
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			"/block_start/1",
			struct{ Invalid string }{Invalid: "test"},
			http.StatusBadRequest)
	})

	t.Run("invalid start time", func(t *testing.T) {
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			fmt.Sprintf("/block_start/%d", utils.BID),
			models.BodyStart{Start: "invalid start"},
			http.StatusBadRequest)
	})

	t.Run("not found", func(t *testing.T) {
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			"/block_start/12",
			models.BodyStart{Start: utils.BStartUpdated},
			http.StatusNotFound)
	})

	t.Run("valid body", func(t *testing.T) {
		db.AddBlock(utils.TestBlockCreate())
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			fmt.Sprintf("/block_start/%d", utils.BID),
			models.BodyStart{Start: utils.BStartUpdated},
			http.StatusOK)
	})
}

func TestUpdateBlockEndRoute(t *testing.T) {
	db := database.GetNewTestDatabase()
	defer db.Close()
	r := NewRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("bad query param", func(t *testing.T) {
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodPut,
			"/block_end/a",
			http.StatusBadRequest)
	})

	t.Run("no body", func(t *testing.T) {
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodPut,
			"/block_end/1",
			http.StatusBadRequest)
	})

	t.Run("invalid body", func(t *testing.T) {
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			"/block_end/1",
			struct{ Invalid string }{Invalid: "test"},
			http.StatusBadRequest)
	})

	t.Run("invalid end time", func(t *testing.T) {
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			fmt.Sprintf("/block_start/%d", utils.BID),
			models.BodyEnd{End: "invalid end"},
			http.StatusBadRequest)
	})

	t.Run("not found", func(t *testing.T) {
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			"/block_end/12",
			models.BodyEnd{End: utils.BEndUpdated},
			http.StatusNotFound)
	})

	t.Run("valid body", func(t *testing.T) {
		db.AddBlock(utils.TestBlockCreate())
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			fmt.Sprintf("/block_end/%d", utils.BID),
			models.BodyEnd{End: utils.BEndUpdated},
			http.StatusOK)
	})
}

func TestUpdateBlockHomeofficeRoute(t *testing.T) {
	db := database.GetNewTestDatabase()
	defer db.Close()
	r := NewRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("bad query param", func(t *testing.T) {
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodPut,
			"/block_homeoffice/a",
			http.StatusBadRequest)
	})

	t.Run("no body", func(t *testing.T) {
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodPut,
			"/block_homeoffice/1",
			http.StatusBadRequest)
	})

	t.Run("invalid body", func(t *testing.T) {
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			"/block_homeoffice/1",
			struct{ Invalid string }{Invalid: "test"},
			http.StatusBadRequest)
	})

	t.Run("not found", func(t *testing.T) {
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			"/block_homeoffice/12",
			models.BodyHomeoffice{Homeoffice: true},
			http.StatusNotFound)
	})

	t.Run("valid body", func(t *testing.T) {
		db.AddBlock(utils.TestBlockCreate())
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			fmt.Sprintf("/block_homeoffice/%d", utils.BID),
			models.BodyHomeoffice{Homeoffice: true},
			http.StatusOK)
	})
}

func TestUpdatePauseStartRoute(t *testing.T) {
	db := database.GetNewTestDatabase()
	defer db.Close()
	r := NewRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("bad query param", func(t *testing.T) {
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodPut,
			"/pause_start/a",
			http.StatusBadRequest)
	})

	t.Run("no body", func(t *testing.T) {
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodPut,
			"/pause_start/1",
			http.StatusBadRequest)
	})

	t.Run("invalid body", func(t *testing.T) {
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			"/pause_start/1",
			struct{ Invalid string }{Invalid: "test"},
			http.StatusBadRequest)
	})

	t.Run("invalid start time", func(t *testing.T) {
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			fmt.Sprintf("/pause_start/%d", utils.PID),
			models.BodyStart{Start: "invalid start"},
			http.StatusBadRequest)
	})

	t.Run("not found", func(t *testing.T) {
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			"/pause_start/12",
			models.BodyStart{Start: utils.PStartUpdated},
			http.StatusNotFound)
	})

	t.Run("valid body", func(t *testing.T) {
		db.AddBlock(utils.TestBlockCreate())
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			fmt.Sprintf("/pause_start/%d", utils.BID),
			models.BodyStart{Start: utils.PStartUpdated},
			http.StatusOK)
	})
}

func TestUpdatePauseEndRoute(t *testing.T) {
	db := database.GetNewTestDatabase()
	defer db.Close()
	r := NewRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("bad query param", func(t *testing.T) {
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodPut,
			"/pause_end/a",
			http.StatusBadRequest)
	})

	t.Run("no body", func(t *testing.T) {
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodPut,
			"/pause_end/1",
			http.StatusBadRequest)
	})

	t.Run("invalid body", func(t *testing.T) {
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			"/pause_end/1",
			struct{ Invalid string }{Invalid: "test"},
			http.StatusBadRequest)
	})

	t.Run("invalid end time", func(t *testing.T) {
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			fmt.Sprintf("/pause_end/%d", utils.PID),
			models.BodyEnd{End: "invalid end"},
			http.StatusBadRequest)
	})

	t.Run("not found", func(t *testing.T) {
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			"/pause_end/12",
			models.BodyEnd{End: utils.PEndUpdated},
			http.StatusNotFound)
	})

	t.Run("valid body", func(t *testing.T) {
		db.AddBlock(utils.TestBlockCreate())
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			fmt.Sprintf("/pause_end/%d", utils.BID),
			models.BodyEnd{End: utils.PEndUpdated},
			http.StatusOK)
	})
}

func TestDeleteBlockRoute(t *testing.T) {
	db := database.GetNewTestDatabase()
	defer db.Close()
	r := NewRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("invalid query param", func(t *testing.T) {
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodDelete,
			"/block/a",
			http.StatusBadRequest)
	})

	t.Run("block not found", func(t *testing.T) {
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodDelete,
			"/block/12",
			http.StatusNotFound)
	})

	t.Run("valid request", func(t *testing.T) {
		db.AddBlock(utils.TestBlockCreate())
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodDelete,
			fmt.Sprintf("/block/%d", utils.BID),
			http.StatusOK)
	})
}

func TestGetBlockByIDRoute(t *testing.T) {
	db := database.GetNewTestDatabase()
	defer db.Close()
	r := NewRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("invalid query param", func(t *testing.T) {
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodGet,
			"/block/a",
			http.StatusBadRequest)
	})

	t.Run("not found", func(t *testing.T) {
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodGet,
			"/block/12",
			http.StatusInternalServerError)
	})

	t.Run("valid request", func(t *testing.T) {
		db.AddBlock(utils.TestBlockCreate())
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodGet,
			fmt.Sprintf("/block/%d", utils.BID),
			http.StatusOK)
	})
}

func TestGetAllBlocksRoute(t *testing.T) {
	db := database.GetNewTestDatabase()
	defer db.Close()
	r := NewRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("no blocks available", func(t *testing.T) {
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodGet,
			"/block",
			http.StatusNotFound)
	})

	t.Run("blocks found", func(t *testing.T) {
		db.AddBlock(utils.TestBlockCreate())
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodGet,
			"/block",
			http.StatusOK)
	})
}

func TestAddPauseRoute(t *testing.T) {
	db := database.GetNewTestDatabase()
	defer db.Close()
	r := NewRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("no body", func(t *testing.T) {
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodPost,
			"/pause",
			http.StatusBadRequest)
	})

	t.Run("no block with blockID", func(t *testing.T) {
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPost,
			"/pause",
			utils.TestPauseCreate(),
			http.StatusInternalServerError)
	})

	t.Run("valid body", func(t *testing.T) {
		db.AddBlock(utils.TestBlockCreateWithoutPause())
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPost,
			"/pause",
			utils.TestPauseCreate(),
			http.StatusOK)
	})
}

func TestUpdatePauseRoute(t *testing.T) {
	db := database.GetNewTestDatabase()
	defer db.Close()
	r := NewRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("no body", func(t *testing.T) {
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodPut,
			"/pause",
			http.StatusBadRequest)
	})

	t.Run("not found", func(t *testing.T) {
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			"/pause",
			utils.TestPauseUpdated(),
			http.StatusNotFound)
	})

	t.Run("valid body", func(t *testing.T) {
		db.AddBlock(utils.TestBlockCreate())
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPut,
			"/pause",
			utils.TestPauseUpdated(),
			http.StatusOK)
	})
}

func TestDeletePauseRoute(t *testing.T) {
	db := database.GetNewTestDatabase()
	defer db.Close()
	r := NewRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("invalid query param", func(t *testing.T) {
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodDelete,
			"/pause/a",
			http.StatusBadRequest)
	})

	t.Run("not found", func(t *testing.T) {
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodDelete,
			"/pause/12",
			http.StatusNotFound)
	})

	t.Run("valid request", func(t *testing.T) {
		db.AddBlock(utils.TestBlockCreate())
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodDelete,
			fmt.Sprintf("/pause/%d", utils.PID),
			http.StatusOK)
	})
}

func TestStartBlockRoute(t *testing.T) {
	db := database.GetNewTestDatabase()
	defer db.Close()
	r := NewRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("invalid query param", func(t *testing.T) {
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodPost,
			"/current_block_start?homeoffice=bad_param",
			http.StatusBadRequest)
	})

	t.Run("valid request", func(t *testing.T) {
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodPost,
			"/current_block_start?homeoffice=false",
			http.StatusOK)
	})

	t.Run("block already active", func(t *testing.T) {
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodPost,
			"/current_block_start?homeoffice=false",
			http.StatusInternalServerError)
	})
}

func TestEndBlockRoute(t *testing.T) {
	db := database.GetNewTestDatabase()
	defer db.Close()
	r := NewRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("block not started", func(t *testing.T) {
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodPost,
			"/current_block_end",
			http.StatusInternalServerError)
	})

	t.Run("pause still active", func(t *testing.T) {
		db.StartBlock(false)
		db.StartPause()
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodPost,
			"/current_block_end",
			http.StatusInternalServerError)
	})

	t.Run("valid request", func(t *testing.T) {
		db.EndPause()
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodPost,
			"/current_block_end",
			http.StatusOK)
	})
}

func TestStartPauseRoute(t *testing.T) {
	db := database.GetNewTestDatabase()
	defer db.Close()
	r := NewRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("no block active", func(t *testing.T) {
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodPost,
			"/current_pause_start",
			http.StatusInternalServerError)
	})

	t.Run("valid request", func(t *testing.T) {
		db.StartBlock(false)
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodPost,
			"/current_pause_start",
			http.StatusOK)
	})

	t.Run("pause already active", func(t *testing.T) {
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodPost,
			"/current_pause_start",
			http.StatusInternalServerError)
	})
}

func TestEndPauseRoute(t *testing.T) {
	db := database.GetNewTestDatabase()
	defer db.Close()
	r := NewRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("no block active", func(t *testing.T) {
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodPost,
			"/current_pause_end",
			http.StatusInternalServerError)
	})

	t.Run("no pause active", func(t *testing.T) {
		db.StartBlock(false)
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodPost,
			"/current_pause_end",
			http.StatusInternalServerError)
	})

	t.Run("valid request", func(t *testing.T) {
		db.StartPause()
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodPost,
			"/current_pause_end",
			http.StatusOK)
	})
}

func TestGetCurrentBlockRoute(t *testing.T) {
	db := database.GetNewTestDatabase()
	defer db.Close()
	r := NewRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("no block active", func(t *testing.T) {
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodGet,
			"/block_current",
			http.StatusInternalServerError)
	})

	t.Run("valid request", func(t *testing.T) {
		db.StartBlock(false)
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodGet,
			"/block_current",
			http.StatusOK)
	})
}

func TestLoginRoute(t *testing.T) {
	db := database.GetNewTestDatabase()
	defer db.Close()
	r := NewRouter(db)
	gin.SetMode(gin.TestMode)

	envTest, err := utils.EnvTestVariables()
	assert.NoError(t, err)

	t.Run("no body", func(t *testing.T) {
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodPost,
			"/login",
			http.StatusBadRequest)
	})

	t.Run("invalid body", func(t *testing.T) {
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPost,
			"/login",
			struct{ Invalid string }{Invalid: "test"},
			http.StatusBadRequest)
	})

	t.Run("invalid email", func(t *testing.T) {
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPost,
			"/login",
			auth.Login{Email: "invalid@gmail.com", Password: envTest.Password},
			http.StatusUnauthorized)
	})

	t.Run("invalid password", func(t *testing.T) {
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPost,
			"/login",
			auth.Login{Email: envTest.Email, Password: "987654321"},
			http.StatusUnauthorized)
	})

	t.Run("valid request", func(t *testing.T) {
		utils.AssertRequestWithBody(
			t,
			r,
			token,
			http.MethodPost,
			"/login",
			auth.Login{Email: envTest.Email, Password: envTest.Password},
			http.StatusOK)
	})
}

func TestRefreshRoute(t *testing.T) {
	db := database.GetNewTestDatabase()
	defer db.Close()
	r := NewRouter(db)
	gin.SetMode(gin.TestMode)

	t.Run("token still valid", func(t *testing.T) {
		utils.AssertRequest(
			t,
			r,
			token,
			http.MethodPost,
			"/refresh",
			http.StatusBadRequest)
	})
}
