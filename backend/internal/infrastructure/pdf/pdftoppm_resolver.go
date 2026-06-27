package pdf

import (
	"fmt"
	"os/exec"
	"runtime"
)

// resolvePdftoppm returns the path to the pdftoppm executable.
//
// On Unix systems (Linux, macOS) it expects pdftoppm to be available in $PATH
// (installed via `apt install poppler-utils` or `brew install poppler`).
//
// On Windows it first checks $PATH (in case the user installed Poppler and
// added it manually), then falls back to the conventional installation path
// used by the official Poppler Windows builds:
//
//	C:\Program Files\poppler\Library\bin\pdftoppm.exe
//
// To install Poppler on Windows download the release from
// https://github.com/oschwartz10612/poppler-windows/releases and unzip to
// C:\Program Files\poppler  (or set POPPLER_BIN_DIR env var to override).
func resolvePdftoppm() (string, error) {
	// Always try PATH first — works on all platforms.
	if path, err := exec.LookPath("pdftoppm"); err == nil {
		return path, nil
	}

	if runtime.GOOS != "windows" {
		// On Unix pdftoppm must be in PATH; nothing else to try.
		return "", fmt.Errorf(
			"pdftoppm not found in $PATH — install Poppler:\n" +
				"  Debian/Ubuntu: sudo apt install poppler-utils\n" +
				"  macOS:         brew install poppler",
		)
	}

	// Windows: check well-known installation locations.
	candidates := []string{
		`C:\Program Files\poppler\Library\bin\pdftoppm.exe`,
		`C:\Program Files\poppler\bin\pdftoppm.exe`,
		`C:\poppler\bin\pdftoppm.exe`,
		`C:\tools\poppler\bin\pdftoppm.exe`,
	}

	for _, p := range candidates {
		if path, err := exec.LookPath(p); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf(
		"pdftoppm not found on Windows.\n" +
			"Download Poppler for Windows from:\n" +
			"  https://github.com/oschwartz10612/poppler-windows/releases\n" +
			"Unzip to C:\\Program Files\\poppler  and restart the application.\n" +
			"Alternatively add the bin\\ directory to your PATH.",
	)
}
