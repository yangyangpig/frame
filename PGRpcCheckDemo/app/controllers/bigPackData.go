package controllers

import (
	_ "encoding/json"
	"fmt"
)

type BigPack struct {
	Route     int
	Container map[string]string
}

func (bigPack *BigPack) CreateBigPack() map[string]string {

	for r := 0; r < bigPack.Route; r++ {
		s := fmt.Sprintf("student- %d  age is :  %d", r, r)
		key := fmt.Sprintf("student-%d", r)
		bigPack.Container[key] = s
	}
	return bigPack.Container
}
