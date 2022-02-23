package backend

import (
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (b *App) about(_ *menu.CallbackData) {
	option := runtime.MessageDialogOptions{
		Title: "Apollo Synchronizer",
		Message: "Help developer to sync between local file and remote apollo " +
			"portal web since portal web is so messy to use.",
	}

	if _, err := runtime.MessageDialog(b.ctx, option); err != nil {
		b.errorf("about failed: %v", err)
	}
}

func (b *App) openPreferences(_ *menu.CallbackData) {
	fmt.Println("openPreferences")
}

func (b *App) setConfigPath(_ *menu.CallbackData) {
	fmt.Println("setConfigPath")
}
