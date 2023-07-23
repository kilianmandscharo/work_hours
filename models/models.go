package models

import (
	"github.com/kilianmandscharo/work_hours/datetime"
)

type Block struct {
	Id         int     `json:"id" binding:"required"`
	Start      string  `json:"start" binding:"required"`
	End        string  `json:"end" binding:"required"`
	Homeoffice bool    `json:"homeoffice"`
	Pauses     []Pause `json:"pauses"`
}

func (b *Block) Valid() bool {
	for _, pause := range b.Pauses {
		if !pause.Valid() {
			return false
		}
	}

	return datetime.IsValidRFC3339(b.Start) && datetime.IsValidRFC3339(b.End)
}

type BlockCreate struct {
	Start      string                `json:"start" binding:"required"`
	End        string                `json:"end" binding:"required"`
	Homeoffice bool                  `json:"homeoffice"`
	Pauses     []PauseWithoutBlockID `json:"pauses"`
}

func (b *BlockCreate) Valid() bool {
	for _, pause := range b.Pauses {
		if !pause.Valid() {
			return false
		}
	}

	return datetime.IsValidRFC3339(b.Start) && datetime.IsValidRFC3339(b.End)
}

type Pause struct {
	Id      int    `json:"id" binding:"required"`
	Start   string `json:"start" binding:"required"`
	End     string `json:"end" binding:"required"`
	BlockID int    `json:"blockID" binding:"required"`
}

func (p *Pause) Valid() bool {
	return datetime.IsValidRFC3339(p.Start) && datetime.IsValidRFC3339(p.End)
}

type PauseCreate struct {
	Start   string `json:"start" binding:"required"`
	End     string `json:"end" binding:"required"`
	BlockID int    `json:"blockID" binding:"required"`
}

func (p *PauseCreate) Valid() bool {
	return datetime.IsValidRFC3339(p.Start) && datetime.IsValidRFC3339(p.End)
}

type PauseWithoutBlockID struct {
	Start string `json:"start" binding:"required"`
	End   string `json:"end" binding:"required"`
}

func (p *PauseWithoutBlockID) Valid() bool {
	return datetime.IsValidRFC3339(p.Start) && datetime.IsValidRFC3339(p.End)
}

type BodyStart struct {
	Start string `json:"start" binding:"required"`
}

func (b *BodyStart) Valid() bool {
	return datetime.IsValidRFC3339(b.Start)
}

type BodyEnd struct {
	End string `json:"end" binding:"required"`
}

func (b *BodyEnd) Valid() bool {
	return datetime.IsValidRFC3339(b.End)
}

type BodyHomeoffice struct {
	Homeoffice bool `json:"homeoffice" binding:"required"`
}
