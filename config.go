package config

import (
	"bufio"

	"io"
	"os"
	"strconv"
	"strings"
)

var config map[string]string

func Load(file string) {
	config = make(map[string]string)
	f, e := os.OpenFile(file, os.O_RDONLY, 0666)
	if e != nil {
		//TODO log
		return
	}
	defer f.Close()
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
		tmp := strings.Split(strings.TrimSpace(string(line)), "=")
		if len(tmp) != 2 {
			//TODO log
			return
		}
		config[tmp[0]] = tmp[1]
	}
}

func GetMustUint32(name string) uint32 {
	if d, ok := config[name]; ok {
		t, e := strconv.ParseUint(d, 10, 32)
		if e != nil {
			//TODO  log
			return 0
		}
		return uint32(t)
	}
	return 0
}

func GetMustString(n string) string {
	return config[n]
}
