// go:build !windows && !darwin

package browser

import (
	"os/exec"
	"strings"
)

func command() (string, []string) {
	var isWSL bool
	if releaseData, err := exec.Command("uname", "-r").Output(); err == nil {
		isWSL = strings.Contains(strings.ToLower(string(releaseData)), "microsoft")
	}
	if isWSL {
		return "cmd.exe", []string{"/c", "start"}
	}
	return "xdg-open", nil
}
