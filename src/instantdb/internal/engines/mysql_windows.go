// +build windows

package engines

import (
	"path/filepath"
)

func getLibraryPathEnv(binaryDir string) string {
	libPath := filepath.Join(binaryDir, "lib")
	return "PATH=" + libPath
}
