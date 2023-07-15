package main

import "time"

type Block struct {
	Id         int     `json:"id" binding:"required"`
	Start      string  `json:"start" binding:"required"`
	End        string  `json:"end" binding:"required"`
	Homeoffice bool    `json:"homeoffice"`
	Pauses     []Pause `json:"pauses"`
}

func (b *Block) valid() bool {
	for _, pause := range b.Pauses {
		if !pause.valid() {
			return false
		}
	}

	_, startErr := time.Parse(time.RFC3339, b.Start)
	_, endErr := time.Parse(time.RFC3339, b.End)

	return startErr == nil && endErr == nil
}

type BlockCreate struct {
	Start      string                `json:"start" binding:"required"`
	End        string                `json:"end" binding:"required"`
	Homeoffice bool                  `json:"homeoffice"`
	Pauses     []PauseWithoutBlockID `json:"pauses"`
}

func (b *BlockCreate) valid() bool {
	for _, pause := range b.Pauses {
		if !pause.valid() {
			return false
		}
	}
	_, startErr := time.Parse(time.RFC3339, b.Start)
	_, endErr := time.Parse(time.RFC3339, b.End)

	return startErr == nil && endErr == nil
}

type Pause struct {
	Id      int    `json:"id" binding:"required"`
	Start   string `json:"start" binding:"required"`
	End     string `json:"end" binding:"required"`
	BlockID int    `json:"blockID" binding:"required"`
}

func (b *Pause) valid() bool {
	_, startErr := time.Parse(time.RFC3339, b.Start)
	_, endErr := time.Parse(time.RFC3339, b.End)

	return startErr == nil && endErr == nil
}

type PauseCreate struct {
	Start   string `json:"start" binding:"required"`
	End     string `json:"end" binding:"required"`
	BlockID int    `json:"blockID" binding:"required"`
}

func (b *PauseCreate) valid() bool {
	_, startErr := time.Parse(time.RFC3339, b.Start)
	_, endErr := time.Parse(time.RFC3339, b.End)

	return startErr == nil && endErr == nil
}

type PauseWithoutBlockID struct {
	Start string `json:"start" binding:"required"`
	End   string `json:"end" binding:"required"`
}

func (b *PauseWithoutBlockID) valid() bool {
	_, startErr := time.Parse(time.RFC3339, b.Start)
	_, endErr := time.Parse(time.RFC3339, b.End)

	return startErr == nil && endErr == nil
}

type BodyStart struct {
	Start string `json:"start" binding:"required"`
}

func (b *BodyStart) valid() bool {
	_, err := time.Parse(time.RFC3339, b.Start)

	return err == nil
}

type BodyEnd struct {
	End string `json:"end" binding:"required"`
}

func (b *BodyEnd) valid() bool {
	_, err := time.Parse(time.RFC3339, b.End)

	return err == nil
}

type BodyHomeoffice struct {
	Homeoffice bool `json:"homeoffice" binding:"required"`
}
