package tss

import (
	"fmt"
	"github.com/nikita-tomilov/gotsdb/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

//TODO: unify tests!!!
//https://github.com/segmentio/testdemo/blob/master/partyparrot/partyparrot_test.go
func TestQlTss_SanityCheck(t *testing.T) {
	//given
	s := buildQlStorage()
	defer s.CloseStorage()
	dataToStore := buildData()
	//when
	s.Save(testDataSource, dataToStore, 1000)
	retrievedData := s.Retrieve(testDataSource, []string{testTag1, testTag2}, may040520, may050520)
	//then
	assert.Equal(t, dataToStore[testTag1], retrievedData[testTag1], "Data for tag1 for 04.05-05.05 should be same")
	assert.Equal(t, make(map[uint64]float64), retrievedData[testTag2].Points, "Data for tag2 for 04.05-05.05 should be empty")
}

func buildQlStorage() *QlBasedPersistentTSS {
	s := QlBasedPersistentTSS{Path:fmt.Sprintf("/tmp/gotsdb_test/test%d", utils.GetNowMillis())}
	s.InitStorage()
	return &s
}