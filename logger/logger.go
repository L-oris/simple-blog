package logger

import (
	"os"

	logging "github.com/op/go-logging"
)

// Log allows access to formatted logging functionalities
var Log = logging.MustGetLogger("yabb")

func init() {
	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortpkg}.%{shortfunc} ▶ %{level:.10s}%{color:reset} %{message}`,
	)
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)
}
