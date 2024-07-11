package autocomplete

import (
	"embed"
)

//go:embed  "zsh_autocomplete"
var EmbeddedFiles embed.FS
