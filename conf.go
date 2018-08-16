package slog

import (
	"errors"
	"os"
	"strings"

	"github.com/valyala/fastjson"
)

type configuration struct {
	paths    map[string]string
	archives map[string]string
	rules    map[string]rule
}

type rule struct {
	destination string
	level       int
	// string format may be here
}

func readConfiguration(path string) (*configuration, error) {
	//  Open file
	f, e := os.Open(path)
	defer f.Close()

	if e != nil {
		return nil, e
	}

	// Get file info
	fi, e := f.Stat()

	if e != nil {
		return nil, e
	}

	// Read file to array
	buf := make([]byte, fi.Size())

	_, e = f.Read(buf)

	if e != nil {
		return nil, e
	}

	// Parse file
	var p fastjson.Parser

	v, e := p.ParseBytes(buf)

	if e != nil {
		return nil, e
	}

	// Mapping to struct
	return getNewConfFile(v)
}

func getNewConfFile(v *fastjson.Value) (*configuration, error) {
	// Create instance
	conf := &configuration{
		paths:    make(map[string]string),
		archives: make(map[string]string),
		rules:    make(map[string]rule),
	}

	// Process log file paths
	a := v.GetArray("paths")
	var name, path string

	for i := 0; i < len(a); i++ {
		name = string(a[i].GetStringBytes("name"))
		if len(name) == 0 {
			continue
		}

		path = string(a[i].GetStringBytes("path"))

		conf.paths[name] = path
		conf.archives[path] = string(a[i].GetStringBytes("archive"))
	}

	if len(conf.paths) == 0 {
		return nil, errors.New("Config parse error. Paths not found")
	}

	// Process rules
	a = v.GetArray("rules")

	for i := 0; i < len(a); i++ {
		name = string(a[i].GetStringBytes("logger"))
		if len(name) == 0 {
			continue
		}

		conf.rules[name] = rule{
			level:       convertToIntLogLevel(string(a[i].GetStringBytes("level"))),
			destination: string(a[i].GetStringBytes("file")),
		}
	}

	if len(conf.rules) == 0 {
		return nil, errors.New("Config parse error. Rules not found")
	}

	return conf, nil
}

func convertToIntLogLevel(level string) int {
	if len(level) == 0 {
		return 0
	}

	lvl := strings.ToLower(level)

	switch {
	case strings.HasPrefix(lvl, "info"):
		return 1
	case strings.HasPrefix(lvl, "warn"):
		return 2
	case strings.HasPrefix(lvl, "error"):
		return 3
	default:
		return 0
	}
}
