#!/usr/bin/env bash
set -eux

# This command replaces "robintest" by "robintest --cover=AB.covout"
# in each of the .bats file of the current directory, where:
# - A is the name of the .bats file
# - B is an increasing number
#
# This way, the coverage of each bats test is written to a unique file
find . -name '*.bats' | \
    sed 'sW\(\./\(.*\)\.bats\)Wawk -i inplace '\''{for(x=1;x<=NF;x++)if($x~/robintest/){sub(/robintest/,"robintest --cover=\2"++i".covout")}}1'\'' \1W' | \
    bash -x

# Run the tests, so coverage files can be obtained
bats . || true

# Merge all coverage files into one
gocovmerge *.covout > merged.covout

# Get a readable coverage report
gocov convert merged.covout | gocov report > coverage-report.txt

# Revert .bats files to their previous state (if in a git repo)
git checkout -- *.bats || true
