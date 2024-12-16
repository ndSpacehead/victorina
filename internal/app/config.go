package app

import "flag"

const defaultPathToDB = "./vic.db"

type config struct {
	dbFilename string
}

func readConfig() config {
	var out config
	flag.StringVar(&out.dbFilename, "b", defaultPathToDB, "database filename")
	flag.Parse()
	return out
}
