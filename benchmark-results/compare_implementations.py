#!/usr/bin/env python3
import argparse
import os
import sys

import matplotlib.pyplot as plt

from common import print_available_branches, print_available_implementations, read_all_lines


def compare_implementations(branch_name, bench_name, target_names):
    print("Going to compare " + str(target_names) + " at benchmark " + bench_name + " at branch " + branch_name)
    file = os.curdir + "/" + branch_name + "/benchmark_results_Benchmark" + bench_name + ".csv"
    file = read_all_lines(file)
    header = file[0]
    avail_impls = header.split()

    indexes = list(map(lambda x: avail_impls.index(x), target_names))
    rows = list(map(lambda x: x.split(), file[1:]))
    x_axis = list(map(lambda x: int(x[0]), rows))

    plt.title('Performance difference in benchmark ' + bench_name + ' at branch ' + branch_name, wrap=True)
    for i in range(0, len(target_names)):
        target_name = target_names[i]
        idx = indexes[i]
        y_axis_old = list(map(lambda x: int(x[idx]), rows))
        plt.plot(x_axis, y_axis_old, '-', label=target_name)

    plt.legend()
    plt.show()


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Compares two CSV files.')
    parser.add_argument('branch', metavar='branch', type=str, help='Branch')
    parser.add_argument('bench', metavar='bench', type=str, help='Benchmark name')
    parser.add_argument('targets', metavar='target', type=str, nargs='+', help='Benchmark target implementation[s]')
    if len(sys.argv) == 1:
        parser.print_help()
        print_available_branches()
        print_available_implementations()
        print("\nExample: compare_implementations.py master<...> DataReading INMEM LSM")
        sys.exit(1)
    args = parser.parse_args()
    compare_implementations(args.branch, args.bench, args.targets)
