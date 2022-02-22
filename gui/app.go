package main

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	*appState

	ctx    context.Context
	logger logger.Logger
}

type clusterSetting struct {
	Title         string   `json:"title"`
	Secret        string   `json:"secret"`
	Clusters      []string `json:"clusters"`
	Env           string   `json:"env"`
	PortalAddress string   `json:"portalAddr"`
	Account       string   `json:"account"`
	LocalDir      string   `json:"fs"`
}

type appState struct {
	Initialized bool             `json:"initialized"`
	Clusters    []clusterSetting `json:"clusters"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		appState: new(appState),
	}
}

func (b *App) loadConfigAndState() error {
	home, _ := os.UserHomeDir()
	asyHome := filepath.Join(home, ".asy")
	c := filepath.Join(asyHome, "asyrc")

	fd, err := os.OpenFile(c, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	bytes, err := io.ReadAll(fd)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(bytes, b.appState); err != nil {
		return err
	}

	b.debugf("loadConfigAndState: %+v", b.appState)

	return nil
}

// startup is called at application startup
func (b *App) startup(ctx context.Context) {
	// Perform your setup here
	b.ctx = ctx
	b.logger = logger.NewDefaultLogger()
	if err := b.loadConfigAndState(); err != nil {
		b.errorf("Failed to load config and state: %v", err)
	}

	// register a menu item
	mymenu := menu.NewMenuFromItems(
		menu.SubMenu("", menu.NewMenuFromItems(
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
	runtime.MenuSetApplicationMenu(ctx, mymenu)
}

// domReady is called after the front-end dom has been loaded
func (b *App) domReady(ctx context.Context) {
	// Add your action here
	// check config path and alert a tip
}

// shutdown is called at application termination
func (b *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}
