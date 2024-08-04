package models

import (
	"fyne.io/fyne/v2/data/binding"
)

type BaseState interface{}

type GlobalState struct {
	BaseState
	UserInfo binding.String
}

func NewGlobalState() *GlobalState {
	return &GlobalState{
		UserInfo: binding.NewString(),
	}
}
