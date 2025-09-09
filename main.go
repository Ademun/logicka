package main

import (
	"embed"
	"logicka/lib"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "logicka",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		Bind: []interface{}{
			&lib.Logicka{},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
