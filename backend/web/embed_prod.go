//go:build prod

package web

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var embeddedDist embed.FS

// Dist returns the embedded dist/ tree from the prod build.
func Dist() (fs.FS, error) {
	return fs.Sub(embeddedDist, "dist")
}
