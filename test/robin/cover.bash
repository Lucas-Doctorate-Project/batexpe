#!/usr/bin/env bash
set -eux

# Create a copy of each bats file
rm -f *-cover.bats
find . -name '*.bats' | \
    sed 'sW\(\(.*\)/\(.*\)\.bats\)Wcp -f \1 \2/\3-cover.batsW' | bash -x

# This command replaces "robintest" by "robintest --cover=AB.covout"
# in each of the .bats file of the current directory, where:
# - A is the name of the .bats file
# - B is an increasing number
#
# This way, the coverage of each bats test is written to a unique file
find . -name '*-cover.bats' | \
    sed 'sW\(\./\(.*\)\.bats\)Wawk -i inplace '\''{for(x=1;x<=NF;x++)if($x~/robintest/){sub(/robintest/,"robintest --cover=\2"++i".covout")}}1'\'' \1W' | \
    bash -x

# Clean previous coverage results if needed
rm -f *.covout coverage-report.txt

# Run the tests, so coverage files can be obtained
bats *-cover.bats || true

# Merge all coverage files into one
gocovmerge *.covout > merged.covout

# Get a readable coverage report
gocov convert merged.covout | gocov report > coverage-report.txt
