package tss

import (
	"bytes"
	"encoding/binary"
	"fmt"
	log "github.com/jeanphorn/log4go"
	"github.com/nikita-tomilov/gotsdb/proto"
	"github.com/nikita-tomilov/gotsdb/utils"
	bolt "go.etcd.io/bbolt"
	"math"
	"os"
	"time"
)

type BboltEntry struct {
	Timestamp uint64
	Value     float64
	ExpireAt  uint64
}

type BboltTSS struct {
	Path               string `summer.property:"tss.filePath|/tmp/gotsdb/tss"`
	dbFilePath         string
	db                 *bolt.DB
	periodBetweenWipes time.Duration
	isRunning          bool
}

func (b *BboltTSS) InitStorage() {
	_ = os.MkdirAll(b.Path, os.ModePerm)
	b.dbFilePath = b.Path + "/db.bin"
	if b.periodBetweenWipes == 0*time.Second {
		b.periodBetweenWipes = time.Second * 5
	}

	db, err := bolt.Open(b.dbFilePath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Error(err)
	}
	b.db = db
	go func(b *BboltTSS) {
		time.Sleep(b.periodBetweenWipes)
		for b.isRunning {
			b.expirationCycle()
			time.Sleep(b.periodBetweenWipes)
		}
	}(b)
}

func (b *BboltTSS) CloseStorage() {
	b.db.Close()
	b.isRunning = false
}

func (b *BboltTSS) Save(dataSource string, data map[string]*proto.TSPoints, expirationMillis uint64) {
	b.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(dataSource))
		now := utils.GetNowMillis()
		expireAt := now + expirationMillis
		if expirationMillis == 0 {
			expireAt = 0
		}
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		b.saveToDataSourceBucket(bucket, data, expireAt)
		return nil
	})
}

func (b *BboltTSS) SaveBatch(dataSource string, data []*proto.TSPoint, expirationMillis uint64) {
	b.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(dataSource))
		now := utils.GetNowMillis()
		expireAt := now + expirationMillis
		if expirationMillis == 0 {
			expireAt = 0
		}
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		b.saveBatchToDataSourceBucket(bucket, data, expireAt)
		return nil
	})
}

func (b *BboltTSS) encodeTimestamp(ts uint64) []byte {
	//return utils.Uint64ToByte(ts) does not work because it will not give you "sortable" byte array
	seconds := ts / 1000
	millis := ts - (seconds * 1000)
	x := time.Unix(int64(seconds), 0).Format("2006-01-02T15:04:05Z")
	str := fmt.Sprintf("%s.%04d", x, millis)
	return []byte(str)
}

func (b *BboltTSS) encodeEntry(entry BboltEntry) []byte {
	payloadLen := 8 + 8 + 8
	arr := make([]byte, payloadLen)
	binary.LittleEndian.PutUint64(arr, entry.Timestamp)
	binary.LittleEndian.PutUint64(arr[8:], math.Float64bits(entry.Value))
	binary.LittleEndian.PutUint64(arr[16:], entry.ExpireAt)
	return arr
}

func (b *BboltTSS) decodeEntry(arr []byte) BboltEntry {
	timestamp := binary.LittleEndian.Uint64(arr)
	value := math.Float64frombits(binary.LittleEndian.Uint64(arr[8:]))
	expiresAt := binary.LittleEndian.Uint64(arr[16:])
	return BboltEntry{
		Timestamp: timestamp,
		Value:     value,
		ExpireAt:  expiresAt,
	}
}

func (b *BboltTSS) saveToDataSourceBucket(bucket *bolt.Bucket, data map[string]*proto.TSPoints, expireAt uint64) {
	for tag, points := range data {
		bucketForTag, err := bucket.CreateBucketIfNotExists([]byte(tag))
		if err != nil {
			log.Error("error in saveToDataSourceBucket: {}", err)
			return
		}
		for ts, val := range points.Points {
			e := BboltEntry{
				Timestamp: ts,
				Value:     val,
				ExpireAt:  expireAt,
			}
			key := b.encodeTimestamp(ts)
			value := b.encodeEntry(e)
			bucketForTag.Put(key, value)
		}
	}
}

func (b *BboltTSS) saveBatchToDataSourceBucket(bucket *bolt.Bucket, data []*proto.TSPoint, expireAt uint64) {
	for _, entry := range data {
		tag := entry.Tag
		ts := entry.Timestamp
		val := entry.Value

		bucketForTag, err := bucket.CreateBucketIfNotExists([]byte(tag))
		if err != nil {
			log.Error("error in saveToDataSourceBucket: {}", err)
			return
		}

		e := BboltEntry{
			Timestamp: ts,
			Value:     val,
			ExpireAt:  expireAt,
		}
		key := b.encodeTimestamp(ts)
		value := b.encodeEntry(e)
		bucketForTag.Put(key, value)
	}
}

func (b *BboltTSS) Retrieve(dataSource string, tags []string, fromTimestamp uint64, toTimestamp uint64) map[string]*proto.TSPoints {
	ans := make(map[string]*proto.TSPoints)
	now := utils.GetNowMillis()
	b.db.View(func(tx *bolt.Tx) error {
		bucketForDataSource := tx.Bucket([]byte(dataSource))
		if bucketForDataSource == nil {
			return nil
		}
		for _, tag := range tags {
			ansForTag := make(map[uint64]float64)
			bucketForTag := bucketForDataSource.Bucket([]byte(tag))
			if bucketForTag == nil {
				ans[tag] = &proto.TSPoints{Points: ansForTag}
				continue
			}

			c := bucketForTag.Cursor()
			min := b.encodeTimestamp(fromTimestamp)
			max := b.encodeTimestamp(toTimestamp)

			for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
				entry := b.decodeEntry(v)
				if (entry.ExpireAt == 0) || (entry.ExpireAt > now) {
					ansForTag[entry.Timestamp] = entry.Value
				}
			}
			ans[tag] = &proto.TSPoints{Points: ansForTag}
		}

		return nil
	})
	return ans
}

func (b *BboltTSS) Availability(dataSource string, fromTimestamp uint64, toTimestamp uint64) []*proto.TSAvailabilityChunk {
	now := utils.GetNowMillis()
	ansMin := uint64(math.MaxUint64)
	ansMax := uint64(0)

	b.db.View(func(tx *bolt.Tx) error {
		bucketForDataSource := tx.Bucket([]byte(dataSource))
		if bucketForDataSource == nil {
			return nil
		}
		bucketForDataSource.ForEach(func(k, v []byte) error {
			if v != nil {
				return nil
			}

			bucketForTag := bucketForDataSource.Bucket(k)
			c := bucketForTag.Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {
				entry := b.decodeEntry(v)
				if (entry.ExpireAt == 0) || (entry.ExpireAt > now) {
					ansMin = utils.Min(entry.Timestamp, ansMin)
					break
				}
			}

			for k, v := c.Last(); k != nil; k, v = c.Prev() {
				entry := b.decodeEntry(v)
				if (entry.ExpireAt == 0) || (entry.ExpireAt > now) {
					ansMax = utils.Max(entry.Timestamp, ansMax)
					break
				}
			}

			return nil
		})
		return nil
	})

	if ansMax < ansMin {
		ans := make([]*proto.TSAvailabilityChunk, 0)
		return ans
	}

	ansMin = utils.Max(fromTimestamp, ansMin)
	ansMax = utils.Min(toTimestamp, ansMax)

	ans := make([]*proto.TSAvailabilityChunk, 1)
	ans[0] = &proto.TSAvailabilityChunk{FromTimestamp: ansMin, ToTimestamp: ansMax}
	return ans
}

func (b *BboltTSS) String() string {
	return fmt.Sprintf("BboltTSS on file %s", b.dbFilePath)
}

func (b *BboltTSS) expirationCycle() {
	b.db.Update(func(tx *bolt.Tx) error {
		tx.ForEach(func(name []byte, bucket *bolt.Bucket) error {
			b.expirationCycleForDataSource(string(name), bucket)
			return nil
		})
		return nil
	})
}

func (b *BboltTSS) expirationCycleForDataSource(ds string, bucket *bolt.Bucket) {
	log.Debug(fmt.Sprintf("Performing expiration for ds %s", ds))
	now := utils.GetNowMillis()
	bucket.ForEach(func(k, v []byte) error {
		if v != nil {
			return nil
		}
		var keysForDeletion [][]byte

		bucketForTag := bucket.Bucket(k)
		bucketForTag.ForEach(func(k, v []byte) error {
			entry := b.decodeEntry(v)
			if (entry.ExpireAt != 0) || (entry.ExpireAt <= now) {
				keysForDeletion = append(keysForDeletion, k)
			}
			return nil
		})

		for _, key := range keysForDeletion {
			bucketForTag.Delete(key)
		}
		return nil
	})
}
