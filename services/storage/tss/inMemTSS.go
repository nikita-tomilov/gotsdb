package tss

import (
	"encoding/json"
	"github.com/google/btree"
	pb "github.com/nikita-tomilov/gotsdb/proto"
	"github.com/nikita-tomilov/gotsdb/utils"
	"sync"
	"time"
)

type InMemTSS struct {
	isRunning          bool
	data               map[string]map[string]*btree.BTree
	lock               sync.Mutex
	periodBetweenWipes time.Duration
	maxEntriesPerTag   int
}

func (f *InMemTSS) InitStorage() {
	f.isRunning = true
	if f.periodBetweenWipes == 0*time.Second {
		f.periodBetweenWipes = time.Second * 5
	}
	if f.maxEntriesPerTag == 0 {
		f.maxEntriesPerTag = 21600
	}
	f.data = make(map[string]map[string]*btree.BTree)
	go func(f *InMemTSS) {
		time.Sleep(f.periodBetweenWipes)
		for f.isRunning {
			f.expirationCycle()
			time.Sleep(f.periodBetweenWipes)
		}
	}(f)
}

func (f *InMemTSS) CloseStorage() {
	f.isRunning = false
}

func (f *InMemTSS) Save(dataSource string, data map[string]*pb.TSPoints, expirationMillis uint64) {
	expireAt := utils.GetNowMillis() + expirationMillis
	if expirationMillis == 0 {
		expireAt = 0
	}
	f.lock.Lock()
	for tag, values := range data {
		tree := f.getTree(dataSource, tag)
		for ts, val := range values.Points {
			entry := InMemTssEntry{Timestamp: ts, ExpiresAt: expireAt, Value: val}
			tree.ReplaceOrInsert(&entry)
		}
	}
	f.lock.Unlock()
}

func (f *InMemTSS) SaveBatch(dataSource string, data []*pb.TSPoint, expirationMillis uint64) {
	expireAt := utils.GetNowMillis() + expirationMillis
	if expirationMillis == 0 {
		expireAt = 0
	}
	f.lock.Lock()
	for _, value := range data {
		tag := value.Tag
		ts := value.Timestamp
		val := value.Value
		tree := f.getTree(dataSource, tag)
		entry := InMemTssEntry{Timestamp: ts, ExpiresAt: expireAt, Value: val}
		tree.ReplaceOrInsert(&entry)
	}
	f.lock.Unlock()
}

func (f *InMemTSS) Retrieve(dataSource string, tags []string, fromTimestamp uint64, toTimestamp uint64) map[string]*pb.TSPoints {
	f.lock.Lock()
	ans := make(map[string]*pb.TSPoints)
	for _, tag := range tags {
		tree := f.getTree(dataSource, tag)
		ansForTag := make(map[uint64]float64)
		tree.AscendRange(buildIndexKey(fromTimestamp), buildIndexKey(toTimestamp+1), func(i btree.Item) bool {
			oe := i.(*InMemTssEntry)
			ansForTag[oe.Timestamp] = oe.Value
			return true
		})
		ans[tag] = &pb.TSPoints{Points: ansForTag}
	}
	f.lock.Unlock()
	return ans
}

func (f *InMemTSS) Availability(dataSource string, fromTimestamp uint64, toTimestamp uint64) []*pb.TSAvailabilityChunk {
	f.lock.Lock()
	ans := make([]*pb.TSAvailabilityChunk, 0)
	if f.containsDs(dataSource) {
		minTs := ^uint64(0) - 1
		maxTs := uint64(0)
		tags, _ := f.data[dataSource]
		for tag, _ := range tags {
			tree := f.getTree(dataSource, tag)
			if tree.Len() == 0 {
				continue
			}
			minE := tree.Min().(*InMemTssEntry)
			maxE := tree.Max().(*InMemTssEntry)

			min := minE.Timestamp
			max := maxE.Timestamp

			if min < minTs {
				minTs = min
			}
			if max > maxTs {
				maxTs = max
			}
		}
		if minTs < fromTimestamp {
			minTs = fromTimestamp
		}
		if maxTs > toTimestamp {
			maxTs = toTimestamp
		}
		if minTs < maxTs {
			ans = append(ans, &pb.TSAvailabilityChunk{FromTimestamp: minTs, ToTimestamp: maxTs})
		}
	}
	f.lock.Unlock()
	return ans
}

func (f *InMemTSS) String() string {
	return "In-Memory TSS"
}

func (f *InMemTSS) containsDs(dataSource string) bool {
	_, found := f.data[dataSource]
	return found
}

func (f *InMemTSS) getTree(dataSource string, tag string) *btree.BTree {
	treesForDs, dsExists := f.data[dataSource]
	if !dsExists {
		f.data[dataSource] = make(map[string]*btree.BTree)
		treesForDs = f.data[dataSource]
	}
	tree, tagExists := treesForDs[tag]
	if !tagExists {
		f.data[dataSource][tag] = f.createTree()
		tree = f.data[dataSource][tag]
	}
	return tree
}

func (f *InMemTSS) createTree() *btree.BTree {
	return btree.New(4)
}

func (f *InMemTSS) expirationCycle() {
	f.lock.Lock()
	toBeDeleted := make([]*InMemTssEntry, 0, 20)
	now := utils.GetNowMillis()
	for _, treesForDs := range f.data {
		for _, treeForTag := range treesForDs {
			treeForTag.Ascend(func(i btree.Item) bool {
				oe := i.(*InMemTssEntry)
				if (oe.ExpiresAt != 0) && (oe.ExpiresAt < now) {
					toBeDeleted = append(toBeDeleted, oe)
				}
				return true
			})
			for _, i := range toBeDeleted {
				treeForTag.Delete(i)
			}
		}
	}
	f.lock.Unlock()
}

func buildIndexKey(ts uint64) btree.Item {
	return &InMemTssEntry{Timestamp: ts}
}

type InMemTssEntry struct {
	Timestamp uint64
	ExpiresAt uint64
	Value     float64
}

func (e *InMemTssEntry) ToString() string {
	b, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (e *InMemTssEntry) Less(than btree.Item) bool {
	oe := than.(*InMemTssEntry)
	return e.Timestamp < oe.Timestamp
}
