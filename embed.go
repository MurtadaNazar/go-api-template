package main

import (
	"embed"
)

//go:embed scaffold
//go:embed scaffold/base/.env.example
//go:embed scaffold/base/.gitignore
var ScaffoldFS embed.FS
