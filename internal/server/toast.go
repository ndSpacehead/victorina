package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	levelInfo    = "info"
	levelSuccess = "success"
	levelWarning = "warning"
	levelDanger  = "danger"
)

type toast struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

func newToast(level, message string) toast {
	return toast{
		Level:   level,
		Message: base64.StdEncoding.EncodeToString([]byte(message)),
	}
}

func infoToast(message string) toast {
	return newToast(levelInfo, message)
}

func successToast(message string) toast {
	return newToast(levelSuccess, message)
}

func warningToast(message string) toast {
	return newToast(levelWarning, message)
}

func dangerToast(message string) toast {
	return newToast(levelDanger, message)
}

func (t toast) serialize() string {
	m := map[string]toast{
		"makeToast": t,
	}
	if data, err := json.Marshal(m); err == nil {
		return string(data)
	}
	return ""
}

func setHXTriggerHeader(w http.ResponseWriter, constructor func(string) toast, format string, args ...any) {
	w.Header().Set(
		"HX-Trigger",
		constructor(fmt.Sprintf(format, args...)).serialize(),
	)
}
