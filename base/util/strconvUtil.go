package util

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"

	log "github.com/cihub/seelog"
)

var randSource = rand.New(rand.NewSource(time.Now().UnixNano()))

func StringToUint64(i string) uint64 {
	if i == "" {
		return 0
	}
	d, e := strconv.ParseUint(i, 10, 64)
	if e != nil {
		log.Info("string convent err ", i, e)
		return 0
	}
	return d
}

func StringToInt64(i string) int64 {
	if i == "" {
		return 0
	}
	d, e := strconv.ParseInt(i, 10, 64)
	if e != nil {
		log.Info("string convent err ", i, e)
		return 0
	}
	return d
}

func StringToUint32(i string) uint32 {
	if i == "" {
		return 0
	}
	d, e := strconv.ParseUint(i, 10, 64)
	if e != nil {
		log.Info("string convent err ", i, e)
		return 0
	}
	return uint32(d)
}

func StringToInt(i string) int {
	if i == "" {
		return 0
	}
	d, e := strconv.ParseInt(i, 10, 64)
	if e != nil {
		log.Info("string convent err ", i, e)
		return 0
	}
	return int(d)
}
func StringToUint16(i string) uint16 {
	d, e := strconv.ParseUint(i, 10, 64)
	if e != nil {
		log.Info("string convent err ", i, e)
		return 0
	}
	return uint16(d)
}
func StringToUint8(i string) uint8 {
	d, e := strconv.ParseUint(i, 10, 64)
	if e != nil {
		log.Info("string convent err ", i, e)
		return 0
	}
	return uint8(d)
}

func StringToFloat64(i string) float64 {
	if i == "" {
		return 0
	}
	d, e := strconv.ParseFloat(i, 64)
	if e != nil {
		log.Info("string convent err ", i, e)
		return 0
	}
	return d
}

func StringToFloat32(i string) float32 {
	d, e := strconv.ParseFloat(i, 64)
	if e != nil {
		log.Info("string convent err ", i, e)
		return 0
	}
	return float32(d)
}

func StringToSliceUint32(str string) []uint32 {
	var r []uint32
	tmp := strings.Split(str, ":")
	for _, t := range tmp {
		r = append(r, StringToUint32(t))
	}
	return r
}
func StringToSliceUint64(str string) []uint64 {
	var r []uint64
	tmp := strings.Split(str, ":")
	for _, t := range tmp {
		r = append(r, StringToUint64(t))
	}
	return r
}
func StringToSliceInt(str string) []int {
	var r []int
	tmp := strings.Split(str, ":")
	for _, t := range tmp {
		r = append(r, StringToInt(t))
	}
	return r
}
func Uint64ToString(i uint64) string {
	return strconv.FormatUint(i, 10)
}

func Uint32ToString(i uint32) string {
	return strconv.FormatUint(uint64(i), 10)
}

func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func IntToString(i int) string {
	return strconv.FormatInt(int64(i), 10)
}

//进行四舍五入，保留n位小数
func RoundFloat(f float64, n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc((f+0.5/pow10_n)*pow10_n) / pow10_n
}

//不进行四舍五入，保留n位小数
func NoRoundFloat(f float64, n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc(f*pow10_n) / pow10_n
}

//[]uint32转化为[]byte
func Uint32sToBytes(args []uint32) []byte {
	bytes := make([]byte, 0)
	for _, arg := range args {
		b := make([]byte, 4)
		binary.BigEndian.PutUint32(b, arg)
		bytes = append(bytes, b...)
	}
	return bytes
}

//[]byte转化为[]uint32
func BytesToUint32s(bytes []byte) []uint32 {
	res := make([]uint32, 0)
	for i := 0; i < len(bytes)/4; i++ {
		b := bytes[4*i : 4*i+4]
		res = append(res, binary.BigEndian.Uint32(b))
	}
	return res
}
func StringToMapU32I64(str, sep1, sep2 string) map[uint32]int64 {
	res := make(map[uint32]int64)
	ss := strings.Split(str, sep1)
	for _, s := range ss {
		ss2 := strings.Split(s, sep2)
		if len(ss2) != 2 {
			continue
		}
		res[StringToUint32(ss2[0])] = StringToInt64(ss2[1])
	}

	return res
}
func StringToMapInt(str, first, second string) map[interface{}]int {
	data := map[interface{}]int{}
	tmp1 := strings.Split(str, first)
	for _, v := range tmp1 {
		tmp2 := strings.Split(v, second)
		if len(tmp2) != 2 {
			continue
		}
		data[tmp2[0]] = StringToInt(tmp2[1])
	}
	return data
}

func StringToMapUint32(str, sep1, sep2 string) map[uint32]uint32 {
	res := make(map[uint32]uint32)
	ss := strings.Split(str, sep1)
	for _, s := range ss {
		ss2 := strings.Split(s, sep2)
		if len(ss2) != 2 {
			continue
		}

		res[StringToUint32(ss2[0])] = StringToUint32(ss2[1])
	}

	return res
}
func StringToMapKStrVUint32(str, sep1, sep2 string) map[string]uint32 {
	res := make(map[string]uint32)
	ss := strings.Split(str, sep1)
	for _, s := range ss {
		ss2 := strings.Split(s, sep2)
		if len(ss2) != 2 {
			continue
		}

		res[ss2[0]] = StringToUint32(ss2[1])
	}

	return res
}

func StringToMapKStrVUint64(str, sep1, sep2 string) map[string]uint64 {
	res := make(map[string]uint64)
	ss := strings.Split(str, sep1)
	for _, s := range ss {
		ss2 := strings.Split(s, sep2)
		if len(ss2) != 2 {
			continue
		}

		res[ss2[0]] = StringToUint64(ss2[1])
	}

	return res
}

func StringToMapKUint32VStr(str, sep1, sep2 string) map[uint32]string {
	res := make(map[uint32]string)
	ss := strings.Split(str, sep1)
	for _, s := range ss {
		ss2 := strings.Split(s, sep2)
		if len(ss2) != 2 {
			continue
		}

		res[StringToUint32(ss2[0])] = ss2[1]
	}

	return res
}

func StringToMapKUint32VF32(str, sep1, sep2 string) map[uint32]float32 {
	res := make(map[uint32]float32)
	ss := strings.Split(str, sep1)
	for _, s := range ss {
		ss2 := strings.Split(s, sep2)
		if len(ss2) != 2 {
			continue
		}

		res[StringToUint32(ss2[0])] = StringToFloat32(ss2[1])
	}

	return res
}

func GetFreePort() (port int, err error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()

	addr := listener.Addr().String()
	fmt.Println(addr)
	_, portString, err := net.SplitHostPort(addr)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(portString)
}
