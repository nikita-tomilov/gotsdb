package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"os"
	"time"
)

func ToString(e interface{}) string {
	b, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func ToNetworkByteArray(e interface{}) []byte {
	var network bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&network) // Will write to network.
	//var dec = gob.NewDecoder(&network) // Will read from network.
	err := enc.Encode(e)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	b := network.Bytes()
	return append([]byte{byte(len(b))}, b...)
}

func FromNetworkByteArray(arr []byte, x interface{}) (interface{}, error) {
	var network *bytes.Buffer

	l := arr[0]
	network = bytes.NewBuffer(arr[1 : l+1])

	dec := gob.NewDecoder(network) // Will read from network.
	err := dec.Decode(x)
	if err != nil {
		log.Fatal("decode error:", err)
		return nil, err
	}
	return x, nil
}

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

func Float64ToByte(f float64) []byte {
	var buf [8]byte
	n := math.Float64bits(f)
	buf[7] = byte(n >> 56)
	buf[6] = byte(n >> 48)
	buf[5] = byte(n >> 40)
	buf[4] = byte(n >> 32)
	buf[3] = byte(n >> 24)
	buf[2] = byte(n >> 16)
	buf[1] = byte(n >> 8)
	buf[0] = byte(n)
	return buf[:]
}

func ByteToFloat64(b []byte) float64 {
	bits := binary.LittleEndian.Uint64(b)
	return math.Float64frombits(bits)
}

func UnixTsToString(ts uint64) string {
	tm := time.Unix(int64(ts/1000), 0)
	return tm.String()
}