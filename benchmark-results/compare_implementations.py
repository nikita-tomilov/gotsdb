#!/usr/bin/env python3
import argparse
import os
import re
import sys
from pathlib import Path

import matplotlib.pyplot as plt

from common import print_available_branches, print_available_implementations, read_all_lines


def find_available_benchmarks(branch):
    csvs = list(Path("./" + branch + "/").rglob("*.csv"))
    csvs = list(map(lambda x: str(x), csvs))
    if len(csvs) == 0:
        print("none.")
        return
    csvs = list(map(lambda x: re.sub(r"^.+[/]", "", x), csvs))
    csvs = list(map(lambda x: x.replace("benchmark_results_Benchmark", ""), csvs))
    csvs = list(map(lambda x: x.replace(".csv", ""), csvs))
    return csvs


def compare_implementations(branch_name, bench_name, target_names):
    print("Going to compare " + str(target_names) + " at benchmark " + bench_name + " at branch " + branch_name)
    file = os.curdir + "/" + branch_name + "/benchmark_results_Benchmark" + bench_name + ".csv"
    file = read_all_lines(file)
    header = file[0]
    avail_impls = header.split()

    indexes = list(map(lambda x: avail_impls.index(x), target_names))
    rows = list(map(lambda x: x.split(), file[1:]))
    x_axis = list(map(lambda x: int(x[0]), rows))

    plt.figure(bench_name)
    plt.title('Performance difference in benchmark ' + bench_name + ' at branch ' + branch_name, wrap=True)
    for i in range(0, len(target_names)):
        target_name = target_names[i]
        idx = indexes[i]
        y_axis_old = list(map(lambda x: int(x[idx]), rows))
        plt.plot(x_axis, y_axis_old, '-', label=target_name)
    plt.legend()


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Compares two CSV files.')
    parser.add_argument('branch', metavar='branch', type=str, help='Branch')
    parser.add_argument('bench', metavar='bench', type=str, help='Coma-separated benchmark names (or "all" for all available)')
    parser.add_argument('targets', metavar='target', type=str, nargs='+', help='Benchmark target implementation[s]')
    if len(sys.argv) == 1:
        parser.print_help()
        print_available_branches()
        print_available_implementations()
        print("\nExample: compare_implementations.py master<...> DataReading INMEM LSM")
        sys.exit(1)
    args = parser.parse_args()
    if args.bench == "all":
        benches = find_available_benchmarks(args.branch)
    else:
        benches = list(map(lambda x: x.strip(), args.bench.split(",")))
    for bench in benches:
        compare_implementations(args.branch, bench, args.targets)
    plt.show()
