package browser

import (
	"os/exec"
)

// OpenURL opens given URL in default browser.
func OpenURL(url string) error {
	cmd, args := command()
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
