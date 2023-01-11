package main

import (
	"embed"

	"messageboard/cmd"
)

//go:embed public
//go:embed views
var embeddedFiles embed.FS

func main() {
	cmd.ExecServer(embeddedFiles)
}
