#!/bin/bash

echo "Benchmark script started"

function run_benchmark {
  BENCHMARK_NAME=$1
  BENCHMARK_RESULTS_FILENAME=benchmark_results_$BENCHMARK_NAME.csv

  rm -rf /tmp/*.txt
  LOGSALL=/tmp/logs_all.txt
  GOGC=off go test 'github.com/nikita-tomilov/gotsdb/services/storage/tss' -bench=$BENCHMARK_NAME -test.run== | tee $LOGSALL

  LOGSLSM=/tmp/logs_lsm.txt
  LOGSSQL=/tmp/logs_sqlite.txt
  LOGSCSV=/tmp/logs_csv.txt
  LOGSBCSV=/tmp/logs_bcsv.txt
  LOGSINM=/tmp/logs_inmem.txt
  LOGSQLP=/tmp/logs_ql.txt
  LOGSBBL=/tmp/logs_bbolt.txt

  grep 'LSM-based' < $LOGSALL > $LOGSLSM
  grep 'SqliteTSS' < $LOGSALL > $LOGSSQL
  grep 'CSV-based_TSS' < $LOGSALL > $LOGSCSV
  grep 'Binary_CSV-based' < $LOGSALL > $LOGSBCSV
  grep 'In-Memory' < $LOGSALL > $LOGSINM
  grep 'QlBasedPersistentTSS' < $LOGSALL > $LOGSQLP
  grep 'BboltTSS' < $LOGSALL > $LOGSBBL

  LOGSLSMCLEAR=/tmp/logs_lsm_clear.txt
  LOGSSQLCLEAR=/tmp/logs_sqlite_clear.txt
  LOGSCSVCLEAR=/tmp/logs_csv_clear.txt
  LOGSBCSVCLEAR=/tmp/logs_bcsv_clear.txt
  LOGSINMCLEAR=/tmp/logs_inm_clear.txt
  LOGSQLPCLEAR=/tmp/logs_qlp_clear.txt
  LOGSBBLCLEAR=/tmp/logs_bbolt_clear.txt

  cat $LOGSLSM | sed 's/.*_|//g' | sed 's/|-6//g' | awk '{print($1,$3)}' > $LOGSLSMCLEAR
  cat $LOGSSQL | sed 's/.*_|//g' | sed 's/|-6//g' | awk '{print($3)}' > $LOGSSQLCLEAR
  cat $LOGSCSV | sed 's/.*_|//g' | sed 's/|-6//g' | awk '{print($3)}' > $LOGSCSVCLEAR
  cat $LOGSBCSV | sed 's/.*_|//g' | sed 's/|-6//g' | awk '{print($3)}' > $LOGSBCSVCLEAR
  cat $LOGSINM | sed 's/.*_|//g' | sed 's/|-6//g' | awk '{print($3)}' > $LOGSINMCLEAR
  cat $LOGSQLP | sed 's/.*_|//g' | sed 's/|-6//g' | awk '{print($3)}' > $LOGSQLPCLEAR
  cat $LOGSBBL | sed 's/.*_|//g' | sed 's/|-6//g' | awk '{print($3)}' > $LOGSBBLCLEAR

  echo "s LSM SQLITE INMEM CSV BCSV QL BBOLT" > $BENCHMARK_RESULTS_FILENAME
  pr -W 128 -t -m  $LOGSLSMCLEAR $LOGSSQLCLEAR $LOGSINMCLEAR $LOGSCSVCLEAR $LOGSBCSVCLEAR $LOGSQLPCLEAR $LOGSBBLCLEAR >> $BENCHMARK_RESULTS_FILENAME

}

function parse_git_dirty() {
  git diff --quiet --ignore-submodules HEAD 2>/dev/null; [ $? -eq 1 ] && echo "*"
}

function parse_git_branch() {
  git branch --no-color 2> /dev/null | sed -e '/^[^*]/d' -e "s/* \(.*\)/\1$(parse_git_dirty)/" | sed 's/[*]//g'
}

function parse_git_hash() {
  git rev-parse --short HEAD 2> /dev/null | sed "s/\(.*\)/_at_\1/"
}

function move_all_to_target_folder() {
  GIT_BRANCH=$(parse_git_branch)$(parse_git_hash)
  TS=$(date --iso-8601=minutes | sed 's/[+].*//g' | sed 's/:/-/g')
  BENCH_TARGET_DIR="./benchmark-results/"$GIT_BRANCH"_at_"$TS

  echo $BENCH_TARGET_DIR
  mkdir -p $BENCH_TARGET_DIR
  mv benchmark_results*.csv $BENCH_TARGET_DIR/
}

run_benchmark BenchmarkDataReading
run_benchmark BenchmarkLatestDataReading
run_benchmark BenchmarkLinearDataWriting
run_benchmark BenchmarkRandomDataWriting
run_benchmark BenchmarkBatchDataWriting

times

move_all_to_target_folder
