package config

import (
	"bufio"
	"errors"
	"io"
	"os"
	"service/utils"
	"strconv"
	"strings"
)

var config map[string]string
var path string

func Load(file string) error {
	config = make(map[string]string)
	path = file
	f, e := os.OpenFile(file, os.O_RDONLY, 0666)
	if e != nil {
		return errors.New("Load Config File err " + e.Error())
	}
	defer f.Close()
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return errors.New("Read Config File err " + e.Error())
		}
		tmp := strings.Split(strings.TrimSpace(string(line)), "=")
		if len(tmp) != 2 {
			return errors.New("Config File err " + e.Error())
		}
		config[tmp[0]] = tmp[1]
	}
	return nil
}

func GetConfigFile() string {
	return path
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
func GetMustInt(name string) int {
	if d, ok := config[name]; ok {
		t, e := strconv.ParseInt(d, 10, 32)
		if e != nil {
			//TODO  log
			return 0
		}
		return int(t)
	}
	return 0
}
func GetMustString(n string) string {
	return config[n]
}

func GetMustArray(name string) []uint32 {
	if d, ok := config[name]; ok {
		return utils.ChangeStringToArrayUint32(d, "|")
	}
	return nil
}

func GetMustMap(name string) map[uint32]uint32 {
	if d, ok := config[name]; ok {
		return utils.ChangeStringToMapUint32(d, "|", ":")
	}
	return nil
}
