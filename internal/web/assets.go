package web

import "embed"

//go:embed views
var ViewFiles embed.FS

//go:embed static
var StaticFiles embed.FS
