#!/bin/bash

function run_benchmark {
  BENCHMARK_NAME=$1
  BENCHMARK_RESULTS_FILENAME=benchmark_results_$BENCHMARK_NAME.csv

  rm -rf /tmp/*.txt
  LOGSALL=/tmp/logs_all.txt
  GOGC=off go test 'github.com/nikita-tomilov/gotsdb/services/storage/tss' -bench=$BENCHMARK_NAME -test.run== | tee $LOGSALL

  LOGSLSM=/tmp/logs_lsm.txt
  LOGSSQL=/tmp/logs_sqlite.txt
  LOGSCSV=/tmp/logs_csv.txt
  LOGSINM=/tmp/logs_inmem.txt
  LOGSQLP=/tmp/logs_ql.txt

  grep 'LSM-based' < $LOGSALL > $LOGSLSM
  grep 'SqliteTSS' < $LOGSALL > $LOGSSQL
  grep 'CSV-based' < $LOGSALL > $LOGSCSV
  grep 'In-Memory' < $LOGSALL > $LOGSINM
  grep 'QlBasedPersistentTSS' < $LOGSALL > $LOGSQLP

  LOGSLSMCLEAR=/tmp/logs_lsm_clear.txt
  LOGSSQLCLEAR=/tmp/logs_sqlite_clear.txt
  LOGSCSVCLEAR=/tmp/logs_csv_clear.txt
  LOGSINMCLEAR=/tmp/logs_inm_clear.txt
  LOGSQLPCLEAR=/tmp/logs_qlp_clear.txt

  cat $LOGSLSM | sed 's/.*_|//g' | sed 's/|-6//g' | awk '{print($1,$3)}' > $LOGSLSMCLEAR
  cat $LOGSSQL | sed 's/.*_|//g' | sed 's/|-6//g' | awk '{print($3)}' > $LOGSSQLCLEAR
  cat $LOGSCSV | sed 's/.*_|//g' | sed 's/|-6//g' | awk '{print($3)}' > $LOGSCSVCLEAR
  cat $LOGSINM | sed 's/.*_|//g' | sed 's/|-6//g' | awk '{print($3)}' > $LOGSINMCLEAR
  cat $LOGSQLP | sed 's/.*_|//g' | sed 's/|-6//g' | awk '{print($3)}' > $LOGSQLPCLEAR

  echo "s LSM INMEM CSV SQLITE QL" > $BENCHMARK_RESULTS_FILENAME
  pr -W 128 -t -m  $LOGSLSMCLEAR $LOGSINMCLEAR $LOGSCSVCLEAR $LOGSSQLCLEAR $LOGSQLPCLEAR >> $BENCHMARK_RESULTS_FILENAME

}

run_benchmark BenchmarkDataReading
run_benchmark BenchmarkLatestDataReading
run_benchmark BenchmarkLinearDataWriting
run_benchmark BenchmarkRandomDataWriting

times