package main

import (
	"b2n3/api"
	"b2n3/config"
	"context"
	"fmt"
)

// App struct
type App struct {
	ctx context.Context
}

//type Result struct {
//	Name   string `json:"name"`
//	Status bool   `json:"status"`
//}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// So we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	video := api.NewVideoInfo(name)
	if video == nil {
		return "出现错误"
	}
	return fmt.Sprintf("这是%v", video.Titles)
}

func (a *App) GetVideoInfo(name string) *api.VideoInfo {
	video := api.NewVideoInfo(name)
	if video == nil {
		return nil
	}
	return video
}

func (a *App) GetBangumiInfo(name string) *api.Bangumi {
	bangumi := api.NewBangumiInfo(name)
	if bangumi == nil {
		return nil
	}
	return bangumi
}

func (a *App) SubmitVideoInfo() {
	api.SumbitVideo()
}
func (a *App) SubmitBangumiInfo() {
	// api.SumbitBangumi()
}

func (a *App) ScanConfiguraitonFiles() (configs []string) {
	loadConfig := config.LoadOptions()
	configs = append(loadConfig, "创建新配置")
	return
}

func (a *App) CreateConfiguration(conf config.Config, name string) string {
	err := config.CreateConfiguration(&conf, name)
	if err != nil {
		return "配置文件创建失败"
	}
	return "配置文件创建成功"
}

func (a *App) GetCurrentConfiguration(name string) *config.Config {
	conf, err := config.LoadConfigFile(name)
	if err != nil {
		return nil
	}
	return conf
}

func (a *App) ChangeConfiguration(conf config.Config) int {
	err := config.ModifyConfiguration(&conf)
	if err != nil {
		return 1
	}
	return 0
}

func (a *App) DeleteConfigurationFile(name string) int {
	err := config.DeleteConfigurationFile(name)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	return 0
}
