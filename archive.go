package slog

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
)

const tickPeriod time.Duration = 24 * time.Hour

const hourToTick = 23
const minuteToTick = 00
const secondsToTick = 03

type ticker struct {
	*time.Timer
}

// TODO сделать, чтобы проверка происходила при записи в лог, а не в отдельной горутине (что тоже неплохая идея)

func Work(logsPath, archivePath map[string]string) {
	yday := time.Now().YearDay()

	for _, v := range logsPath {
		err := filepath.Walk(path.Dir(v),
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				fmt.Printf("Folder: %s\n\tPath: %s, ModTime: %s\n", path, archivePath[v]+path, info.ModTime().String())

				if info.ModTime().YearDay() < yday {
					//return os.Rename(path, archivePath+path)
				}
				return nil
			})
		if err != nil {
			log.Println(err)
		}
	}
}

func ArchiveLogFilesLoop(logsPath, archivePath map[string]string) {
	Work(logsPath, archivePath)
	// t := ticker{time.NewTimer(getNextTickDuration())}
	// for {
	// 	<-t.C

	// 	Work(logsPath, archivePath)
	// 	t.updateTicker()
	// }
}

func getNextTickDuration() time.Duration {
	now := time.Now()
	nextTick := time.Date(now.Year(), now.Month(), now.Day(), hourToTick, minuteToTick, secondsToTick, 0, time.Local)
	if nextTick.Before(now) {
		nextTick = nextTick.Add(tickPeriod)
	}
	return nextTick.Sub(time.Now())
}

func (t *ticker) updateTicker() {
	t.Reset(getNextTickDuration())
}
