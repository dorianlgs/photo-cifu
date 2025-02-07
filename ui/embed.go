package ui

import (
	"embed"
	"io/fs"
)

//go:generate yarn
//go:generate yarn run build
//go:embed all:build
var distDir embed.FS

var DistDirFS, _ = fs.Sub(distDir, "build")
