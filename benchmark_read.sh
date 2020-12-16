#!/bin/bash
rm -rf /tmp/*.txt
LOGSALL=/tmp/logs_all.txt
GOGC=off go test 'github.com/nikita-tomilov/gotsdb/services/storage/tss' -bench=BenchmarkDataReading -test.run== | tee $LOGSALL

LOGSLSM=/tmp/logs_lsm.txt
LOGSSQL=/tmp/logs_sqlite.txt
LOGSCSV=/tmp/logs_csv.txt
LOGSINM=/tmp/logs_inmem.txt

grep 'LSM-based' < $LOGSALL > $LOGSLSM
grep 'SqliteTSS' < $LOGSALL > $LOGSSQL
grep 'CSV-based' < $LOGSALL > $LOGSCSV
grep 'In-Memory' < $LOGSALL > $LOGSINM

LOGSLSMCLEAR=/tmp/logs_lsm_clear.txt
LOGSSQLCLEAR=/tmp/logs_sqlite_clear.txt
LOGSCSVCLEAR=/tmp/logs_csv_clear.txt
LOGSINMCLEAR=/tmp/logs_inm_clear.txt

cat $LOGSLSM | sed 's/.*_|//g' | sed 's/|-6//g' | awk '{print($1,$3)}' > $LOGSLSMCLEAR
cat $LOGSSQL | sed 's/.*_|//g' | sed 's/|-6//g' | awk '{print($3)}' > $LOGSSQLCLEAR
cat $LOGSCSV | sed 's/.*_|//g' | sed 's/|-6//g' | awk '{print($3)}' > $LOGSCSVCLEAR
cat $LOGSINM | sed 's/.*_|//g' | sed 's/|-6//g' | awk '{print($3)}' > $LOGSINMCLEAR

echo "s LSM INMEM CSV SQLITE"
pr -W 128 -t -m  $LOGSLSMCLEAR $LOGSINMCLEAR $LOGSCSVCLEAR $LOGSSQLCLEAR