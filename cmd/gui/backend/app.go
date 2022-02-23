package backend

import (
	"context"
	binarypkg "encoding/binary"
	"encoding/json"
	"io"
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

func read(fp string, binary bool) (bytes []byte, err error) {
	fd, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	if binary {
		binarypkg.Read(fd, binarypkg.LittleEndian, bytes)
	} else {
		bytes, err = io.ReadAll(fd)
	}

	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// save file overwrite.
func save(fp string, content []byte, binary bool) error {
	fd, err := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer fd.Close()

	if binary {
		binarypkg.Write(fd, binarypkg.LittleEndian, content)
	} else {
		_, err = fd.Write(content)
	}

	if err != nil {
		return err
	}

	return nil
}

// load config and statistics for app.
func (b *App) load() error {
	configFp := filepath.Join(appConfigRoot(), "asyrc")
	statisticsFp := filepath.Join(appConfigRoot(), "statistics.bat")

	raw, err := read(configFp, false)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		// config file not exist, create a new one.
		save(configFp, []byte("{}"), false)
		raw = []byte("{}")
	}
	_ = json.Unmarshal(raw, b.config)

	raw2, err2 := read(statisticsFp, true)
	if err2 != nil {
		if !os.IsNotExist(err2) {
		}
	}
	_ = json.Unmarshal(raw2, b.statistics)

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

	err := save(filepath.Join(appConfigRoot(), "asyrc"), b.config.Bytes(), false)
	if err != nil {
		b.errorf("Failed to save config: %v", err)
	}
	err = save(filepath.Join(appConfigRoot(), "statistics.bat"), b.statistics.Bytes(), true)
	if err != nil {
		b.errorf("Failed to save statistics: %v", err)
	}

	b.infof("OnShutdown called finish")
}

type apolloClusterSetting struct {
	Title         string   `json:"title"`
	Secret        string   `json:"secret"`
	Clusters      []string `json:"clusters"`
	Env           string   `json:"env"`
	PortalAddress string   `json:"portalAddr"`
	Account       string   `json:"account"`
	LocalDir      string   `json:"fs"`
}

type config struct {
	_        []func()
	Settings []apolloClusterSetting `json:"clusters"`
}

func (c config) Bytes() []byte {
	bytes, _ := json.Marshal(c)
	return bytes
}

// statistics struct for app
// DONE(@yeqown): fill statistics indicators here
type statistics struct {
	LastOpened  int64 `json:"lastOpenTs"`
	FirstOpened int64 `json:"firstOpenTs"`
	OpenCount   int64 `json:"openCount"`
	OpenTime    int64 `json:"openTime"` // seconds

	UploadCount       int64 `json:"uploadFileCount"`   // how many times the user used upload mode
	UploadFileCount   int64 `json:"uploadFileCount"`   // how many files uploaded
	UploadFileSize    int64 `json:"uploadFileSize"`    // how many bytes uploaded
	UploadFailedCount int64 `json:"uploadFailedCount"` // how many times the user failed to upload a file

	DownloadCount       int64 `json:"downloadCount"`       // how many times the user used download mode
	DownloadFileCount   int64 `json:"downloadFileCount"`   // how many files downloaded
	DownloadFileSize    int64 `json:"downloadFileSize"`    // total downloaded file size
	DownloadFailedCount int64 `json:"downloadFailedCount"` // how many times the user failed to download a file
}

func (s *statistics) Bytes() []byte {
	bytes, _ := json.Marshal(s)
	return bytes
}
