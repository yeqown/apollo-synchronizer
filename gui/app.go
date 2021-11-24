package main

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (b *App) startup(ctx context.Context) {
	// Perform your setup here
	b.ctx = ctx

	// register a menu item
	menu := menu.NewMenuFromItems(
		menu.SubMenu("App", menu.NewMenuFromItems(
			menu.Text("About", nil, b.about),
			menu.Separator(),
			menu.Text("Quit", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
				runtime.Quit(ctx)
			}),
		)),
		menu.SubMenu("Settings", menu.NewMenuFromItems(
			menu.Text("ConfigPath", nil, b.setConfigPath),
			menu.Separator(),
			menu.Text("Preferences", keys.CmdOrCtrl("p"), b.openPreferences),
		)),
	)
	runtime.MenuSetApplicationMenu(ctx, menu)
}

// domReady is called after the front-end dom has been loaded
func (b *App) domReady(ctx context.Context) {
	// Add your action here
}

// shutdown is called at application termination
func (b *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

// Greet returns a greeting for the given name
func (b *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
