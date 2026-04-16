//go:build !prod

package web

import (
	"errors"
	"io/fs"
)

// Dist는 dev 빌드에서 항상 에러를 반환합니다.
// dev 모드는 Vite dev server(5173)가 별도 서빙하므로 embed 가 필요 없음.
func Dist() (fs.FS, error) {
	return nil, errors.New("web.Dist: dev build does not embed frontend (use Vite dev server on :5173)")
}
