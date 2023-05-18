package dist

import "embed"

//go:generate make -C .. dist
//go:embed static
var StaticFS embed.FS
