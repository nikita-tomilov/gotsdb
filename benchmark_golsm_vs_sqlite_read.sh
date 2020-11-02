#!/bin/bash
rm -rf /tmp/*.txt
LOGSALL=/tmp/logs_all.txt
GOGC=off go test 'github.com/nikita-tomilov/gotsdb/services/storage/tss' -bench=BenchmarkDataReading_LSMvsSQLite -test.run== | tee $LOGSALL

LOGSLSM=/tmp/logs_lsm.txt
LOGSSQL=/tmp/logs_sqlite.txt

grep 'LSM-based' < $LOGSALL > $LOGSLSM
grep 'SqliteTSS' < $LOGSALL > $LOGSSQL

LOGSLSMCLEAR=/tmp/logs_lsm_clear.txt
LOGSSQLCLEAR=/tmp/logs_sqlite_clear.txt
cat $LOGSLSM | sed 's/.*_|//g' | sed 's/|-6//g' | awk '{print($1,$3)}' > $LOGSLSMCLEAR
cat $LOGSSQL | sed 's/.*_|//g' | sed 's/|-6//g' | awk '{print($3)}' > $LOGSSQLCLEAR

echo "s LSM SQLITE"
pr -t -m $LOGSLSMCLEAR $LOGSSQLCLEAR