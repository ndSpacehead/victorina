// go:build darwin

package browser

func command() (string, []string) {
	return "open", nil
}
