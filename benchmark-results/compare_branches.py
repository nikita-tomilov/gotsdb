#!/usr/bin/env python3
import argparse
import os
import sys

import matplotlib.pyplot as plt

from common import read_all_lines, print_available_branches, print_available_implementations


def compare_branches(old_branch_name, new_branch_name, bench_name, target_name):
    print("Going to compare " + new_branch_name + " against older version " + old_branch_name)
    file_old = os.curdir + "/" + old_branch_name + "/benchmark_results_Benchmark" + bench_name + ".csv"
    file_new = os.curdir + "/" + new_branch_name + "/benchmark_results_Benchmark" + bench_name + ".csv"
    print("Old: " + file_old)
    print("New: " + file_new)
    file_old = read_all_lines(file_old)
    file_new = read_all_lines(file_new)
    header_old = file_old[0]
    header_new = file_new[0]
    idx_old = header_old.split(" ").index(target_name)
    idx_new = header_new.split(" ").index(target_name)
    rows_old = list(map(lambda x: x.split(), file_old[1:]))
    rows_new = list(map(lambda x: x.split(), file_new[1:]))
    x_axis = list(map(lambda x: int(x[0]), rows_old))
    y_axis_old = list(map(lambda x: int(x[idx_old]), rows_old))
    y_axis_new = list(map(lambda x: int(x[idx_new]), rows_new))
    plt.title('Performance difference in benchmark ' + bench_name + ' for impl ' + target_name, wrap=True)
    plt.plot(x_axis, y_axis_old, 'r-', label="in " + old_branch_name)
    plt.plot(x_axis, y_axis_new, 'b-', label="in " + new_branch_name)
    plt.legend()
    plt.show()


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Compares two CSV files.')
    parser.add_argument('old', metavar='old', type=str, help='Older branch')
    parser.add_argument('new', metavar='new', type=str, help='Newer branch')
    parser.add_argument('bench', metavar='bench', type=str, help='Benchmark name')
    parser.add_argument('target', metavar='target', type=str, help='Benchmark target impl')
    if len(sys.argv) == 1:
        parser.print_help()
        print_available_branches()
        print_available_implementations()
        print("\nExample: compare_branches.py master<...> branch_name<...> DataReading LSM")
        sys.exit(1)
    args = parser.parse_args()
    compare_branches(args.old, args.new, args.bench, args.target)
