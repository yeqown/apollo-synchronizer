package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	config     *config
	statistics *statistics

	ctx    context.Context
	logger logger.Logger
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		config:     new(config),
		statistics: new(statistics),
		ctx:        nil,
		logger:     nil,
	}
}

func appConfigRoot() string {
	home, _ := os.UserHomeDir()
	asyRoot := filepath.Join(home, ".asy")

	return asyRoot
}

var (
	_configFp     string = filepath.Join(appConfigRoot(), "asyrc")
	_statisticsFp string = filepath.Join(appConfigRoot(), "statistics.bat")
)

// load config and statistics for app.
func (b *App) load() error {
	err := read(_configFp, b.config, _ext_json)
	if err != nil {
		b.errorf("Failed to load config: %v", err)
		if !os.IsNotExist(err) {
			return err
		}
		// config file not exist, create a new one.
		save(_configFp, b.config, _ext_json)
	}

	err2 := read(_statisticsFp, b.statistics, _ext_binary)
	if err2 != nil {
		b.errorf("Failed to load statistics: %v", err2)
	}

	b.debugf("load finished, config: %+v, statistics: %+v", b.config, b.statistics)

	return nil
}

// OnStartup is called at application OnStartup
func (b *App) OnStartup(ctx context.Context) {
	b.ctx = ctx
	b.logger = logger.NewDefaultLogger()
	if err := b.load(); err != nil {
		b.errorf("Failed to load config and state: %v", err)
	}

	{
		// statistics
		if b.statistics.LastOpened == 0 {
			b.statistics.FirstOpened = time.Now().Unix()
		}
		b.statistics.LastOpened = time.Now().Unix()
		b.statistics.OpenCount++
	}

	// register a menu item
	mymenu := menu.NewMenuFromItems(
		menu.SubMenu("", menu.NewMenuFromItems(
			menu.Text("Preferences", keys.CmdOrCtrl(","), b.openPreferences),
			menu.Separator(),
			menu.Text("About", nil, b.about),
			menu.Text("Quit", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
				runtime.Quit(ctx)
			}),
		)),
		menu.EditMenu(),
	)
	runtime.MenuSetApplicationMenu(ctx, mymenu)
}

// OnDomReady is called after the front-end dom has been loaded
func (b *App) OnDomReady(ctx context.Context) {
	// Add your action here
	// check config path and alert a tip
}

// OnShutdown is called at application termination
func (b *App) OnShutdown(ctx context.Context) {
	{
		// statistics
		b.statistics.OpenTime += (time.Now().Unix() - b.statistics.LastOpened)
	}

	err := save(_configFp, b.config, _ext_json)
	if err != nil {
		b.errorf("Failed to save config: %v", err)
	}
	err = save(_statisticsFp, b.statistics, _ext_binary)
	if err != nil {
		b.errorf("Failed to save statistics: %v", err)
	}

	b.infof("OnShutdown called finish")
}

type apolloClusterSetting struct {
	Title         string   `json:"title"`
	Secret        string   `json:"secret"`
	Clusters      []string `json:"clusters"`
	Envs          []string `json:"envs"`
	PortalAddress string   `json:"portalAddr"`
	Account       string   `json:"account"`
	LocalDir      string   `json:"fs"`
}

type config struct {
	_        []func()
	Settings []apolloClusterSetting `json:"clusters"`
}

func (c config) Bytes() []byte {
	bytes, err := json.Marshal(c)
	if err != nil {
		fmt.Println("config marshal error:", err)
	}

	return bytes
}

// statistics struct for app
// DONE(@yeqown): fill statistics indicators here
type statistics struct {
	LastOpened  int64 `json:"lastOpenTs"`
	FirstOpened int64 `json:"firstOpenTs"`
	OpenCount   int64 `json:"openCount"`
	OpenTime    int64 `json:"openTime"` // seconds

	UploadCount       int64 `json:"uploadCount"`       // how many times the user used upload mode
	UploadFileCount   int64 `json:"uploadFileCount"`   // how many files uploaded
	UploadFileSize    int64 `json:"uploadFileSize"`    // how many bytes uploaded
	UploadFailedCount int64 `json:"uploadFailedCount"` // how many times the user failed to upload a file

	DownloadCount       int64 `json:"downloadCount"`       // how many times the user used download mode
	DownloadFileCount   int64 `json:"downloadFileCount"`   // how many files downloaded
	DownloadFileSize    int64 `json:"downloadFileSize"`    // total downloaded file size
	DownloadFailedCount int64 `json:"downloadFailedCount"` // how many times the user failed to download a file
}

func (s *statistics) Bytes() []byte {
	bytes, err := json.Marshal(s)
	if err != nil {
		fmt.Println("statistics marshal error:", err)
	}

	return bytes
}
