package main

type Block struct {
	Id         int     `json:"id"`
	Start      string  `json:"start"`
	End        string  `json:"end"`
	Homeoffice bool    `json:"homeoffice"`
	Pauses     []Pause `json:"pauses"`
}

type BlockCreate struct {
	Start      string                `json:"start"`
	End        string                `json:"end"`
	Homeoffice bool                  `json:"homeoffice"`
	Pauses     []PauseWithoutBlockID `json:"pauses"`
}

type Pause struct {
	Id      int    `json:"id"`
	Start   string `json:"start"`
	End     string `json:"end"`
	BlockID int    `json:"blockID"`
}

type PauseCreate struct {
	Start   string `json:"start"`
	End     string `json:"end"`
	BlockID int    `json:"blockID"`
}

type PauseWithoutBlockID struct {
	Start string `json:"start"`
	End   string `json:"end"`
}
