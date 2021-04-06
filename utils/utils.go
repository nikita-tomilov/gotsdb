package utils

import (
	"io/ioutil"
	"log"
	"os"
	"time"
	"unsafe"
)

func Check(e error) {
	if e != nil {
		log.Fatal(e)
		panic(e)
	}
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func DeleteFile(filename string) {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return
	}
	os.Remove(filename)
}

func GetFileNames(path string) []string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	ans := make([]string, len(files))
	for i, file := range files {
		ans[i] = file.Name() //wow no map/reduce/filter as well
	}
	return ans
}

func ComputeHashCode(arr []byte) uint32 {
	var p, hash uint32
	p = 16777619
	hash = 2166136261

	for _, b := range arr {
		b2 := uint32(b)
		hash = (hash ^ b2) * p
	}

	hash += hash << 13
	hash ^= hash >> 7
	hash += hash << 3
	hash ^= hash >> 17
	hash += hash << 5
	return hash
}

func Min(a uint64, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

func MinInt(a int, b int) int {
	if a < b {
		return a
	}
	return b
}


func Max(a uint64, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func GetNowMillis() uint64 {
	return uint64(time.Now().UnixNano() / 1000000)
}

type Envelope struct {
	Value float64
}

const EnvelopeSize = int(unsafe.Sizeof(Envelope{}))

func Float64ToByte(f float64) []byte {
	env := Envelope{Value:f}
	arr := *(*[EnvelopeSize]byte)(unsafe.Pointer(&env))
	return arr[:]
}

func ByteToFloat64(b []byte) float64 {
	rawPointer := unsafe.Pointer(&b[0])
	castedPointer := (*Envelope)(rawPointer)
	return (*castedPointer).Value
}

func UnixTsToString(ts uint64) string {
	tm := time.Unix(int64(ts/1000), 0)
	return tm.String()
}