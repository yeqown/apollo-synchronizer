package main

import (
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/menu"
)

func (b *App) about(_ *menu.CallbackData) {
	fmt.Println("about")
}

func (b *App) openPreferences(_ *menu.CallbackData) {
	fmt.Println("openPreferences")
}

func (b *App) setConfigPath(_ *menu.CallbackData) {
	fmt.Println("setConfigPath")
}
