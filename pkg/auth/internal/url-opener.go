package internal

import (
	"os/exec"
	"runtime"
)

type URLOpener struct{}

func (u *URLOpener) OpenURL(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	exec.Command(cmd, args...).Start()
}

// NewURLOpener returns a new URLOpener
func NewURLOpener() *URLOpener {
	return &URLOpener{}
}
