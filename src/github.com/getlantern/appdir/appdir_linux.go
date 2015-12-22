// +build !windows,!darwin
package appdir

import (
	"fmt"
	"path/filepath"
	"strings"
)

func general(app string) string {
    // TODO: Go for Android currently doesn't support Home Directory.
    // Remove as soon as this is available in the future
    return fmt.Sprintf(".%s", strings.ToLower(app))
}

func logs(app string) string {
	return filepath.Join(general(app), "logs")
}
