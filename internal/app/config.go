package app

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
)

const defaultDBFilename = ":default:"

const defaultPort uint = 9000

type config struct {
	port       uint
	dbFilename string
}

func readConfig() (config, error) {
	var out config
	flag.UintVar(&out.port, "p", defaultPort, "HTTP-server port")
	flag.StringVar(&out.dbFilename, "b", defaultDBFilename, "database filename")
	flag.Parse()
	if out.port >= 1<<16 {
		return out, errors.New("port must be between 0 (default) and 65535")
	}
	if out.port == 0 {
		out.port = defaultPort
	}
	if out.dbFilename == defaultDBFilename {
		ex, err := os.Executable()
		if err != nil {
			return out, err
		}
		out.dbFilename = filepath.Join(filepath.Dir(ex), "vic.db")
	}
	return out, nil
}
