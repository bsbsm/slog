package main

import (
	"time"

	"github.com/bsbsm/slog"
)

func main() {
	// init
	slog.Init(3, 16)

	time.Sleep(10 * time.Second)
	//slog.ArchiveLogFilesLoop(config.paths, config.archives)
	//fmt.Println(strconv.Itoa(slog.LogLevel))
}
