package app

import (
	"flag"
	"os"
	"path/filepath"
)

const defaultDBFilename = ":default:"

type config struct {
	dbFilename string
}

func readConfig() (config, error) {
	var out config
	flag.StringVar(&out.dbFilename, "b", defaultDBFilename, "database filename")
	flag.Parse()
	if out.dbFilename == defaultDBFilename {
		ex, err := os.Executable()
		if err != nil {
			return out, err
		}
		out.dbFilename = filepath.Join(filepath.Dir(ex), "vic.db")
	}
	return out, nil
}
