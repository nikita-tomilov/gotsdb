#!/bin/bash
rm -rf /tmp/*.txt
LOGSALL=/tmp/logs_all.txt
GOGC=off go test 'github.com/nikita-tomilov/gotsdb/services/storage/tss' -bench=BenchmarkDataReading_LSMvsSQLite -test.run== | tee $LOGSALL

LOGSLSM=/tmp/logs_lsm.txt
LOGSSQL=/tmp/logs_sqlite.txt

grep 'LSM-based' < $LOGSALL > $LOGSLSM
grep 'SqliteTSS' < $LOGSALL > $LOGSSQL

echo "LSM"
cat $LOGSLSM

echo "SQLITE"
cat $LOGSSQL