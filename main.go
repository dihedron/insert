package main

import (
	"os"

	log "github.com/dihedron/go-log"
	"github.com/dihedron/zed/cmd"
)

func init() {
	log.SetLevel(log.DBG)
	log.SetStream(os.Stdout, true)
	log.SetTimeFormat("15:04:05.000")
	log.SetPrintCallerInfo(true)
	log.SetPrintSourceInfo(log.SourceInfoShort)
}

func main() {
	cmd.Execute()
}
