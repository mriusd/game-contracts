package main

import (
	"errors"
)


type Backpack struct {
	Grid [][]bool
}

func NewBackpack(width, height int) *Backpack {
	grid := make([][]bool, height)
	for i := range grid {
		grid[i] = make([]bool, width)
	}
	return &Backpack{Grid: grid}
}

func (bp *Backpack) AddItem(item ItemAttributes) (int, int, error) {
	for y := 0; y <= len(bp.Grid)-int(item.ItemHeight.Int64()); y++ {
		for x := 0; x <= len(bp.Grid[y])-int(item.ItemWidth.Int64()); x++ {
			if bp.isSpaceAvailable(x, y, int(item.ItemWidth.Int64()), int(item.ItemHeight.Int64())) {
				bp.fillSpace(x, y, int(item.ItemWidth.Int64()), int(item.ItemHeight.Int64()))
				return x, y, nil
			}
		}
	}
	return -1, -1, errors.New("not enough space in backpack")
}

func (bp *Backpack) isSpaceAvailable(x, y, width, height int) bool {
	for row := y; row < y+height; row++ {
		for col := x; col < x+width; col++ {
			if bp.Grid[row][col] {
				return false
			}
		}
	}
	return true
}

func (bp *Backpack) fillSpace(x, y, width, height int) {
	for row := y; row < y+height; row++ {
		for col := x; col < x+width; col++ {
			bp.Grid[row][col] = true
		}
	}
}