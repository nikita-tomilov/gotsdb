import os
from pathlib import Path


def read_all_lines(filename):
    with open(filename) as f:
        content = f.readlines()
    content = [x.strip() for x in content]
    return content


def print_available_branches():
    print("\nAvailable branches:")
    branches = [x[0] for x in os.walk(os.curdir)]
    branches = list(filter(lambda x: "_at_" in x, branches))
    for branch in branches:
        print(" - " + branch)


def print_available_implementations():
    print("\nAvailable implementations:")
    csvs = list(Path(".").rglob("*.csv"))
    csvs = list(map(lambda x: str(x), csvs))
    if len(csvs) == 0:
        print("none.")
        return
    csv = csvs[0]
    impls = read_all_lines(csv)[0].split(" ")[1:]
    for impl in impls:
        print(" - " + impl)
