package main

type Block struct {
	Id         int     `json:"id" binding:"required"`
	Start      string  `json:"start" binding:"required"`
	End        string  `json:"end" binding:"required"`
	Homeoffice bool    `json:"homeoffice"`
	Pauses     []Pause `json:"pauses"`
}

type BlockCreate struct {
	Start      string                `json:"start" binding:"required"`
	End        string                `json:"end" binding:"required"`
	Homeoffice bool                  `json:"homeoffice"`
	Pauses     []PauseWithoutBlockID `json:"pauses"`
}

type Pause struct {
	Id      int    `json:"id" binding:"required"`
	Start   string `json:"start" binding:"required"`
	End     string `json:"end" binding:"required"`
	BlockID int    `json:"blockID" binding:"required"`
}

type PauseCreate struct {
	Start   string `json:"start" binding:"required"`
	End     string `json:"end" binding:"required"`
	BlockID int    `json:"blockID" binding:"required"`
}

type PauseWithoutBlockID struct {
	Start string `json:"start" binding:"required"`
	End   string `json:"end" binding:"required"`
}

type BodyStart struct {
	Start string `json:"start" binding:"required"`
}

type BodyEnd struct {
	End string `json:"end" binding:"required"`
}

type BodyHomeoffice struct {
	Homeoffice bool `json:"homeoffice" binding:"required"`
}
