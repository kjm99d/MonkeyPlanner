//go:build prod

package web

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var embeddedDist embed.FS

// Dist는 prod 빌드 시 embed.FS 의 dist/ 하위 내용을 반환합니다.
func Dist() (fs.FS, error) {
	return fs.Sub(embeddedDist, "dist")
}
