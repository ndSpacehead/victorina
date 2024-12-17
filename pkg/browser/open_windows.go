// go:build windows

package browser

func command() (string, []string) {
	return "cmd", []string{"/c", "start"}
}
