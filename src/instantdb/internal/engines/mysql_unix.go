// +build !windows

package engines

import (
	"path/filepath"
	"runtime"
)

func getLibraryPathEnv(binaryDir string) string {
	libPath := filepath.Join(binaryDir, "lib")
	if runtime.GOOS == "darwin" {
		return "DYLD_LIBRARY_PATH=" + libPath
	}
	return "LD_LIBRARY_PATH=" + libPath
}
