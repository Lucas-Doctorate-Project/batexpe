#!/usr/bin/env python3
import re

#############
# Functions #
#############

def generate_bats_file(input_filename, output_filename):
    with open(input_filename, "r") as in_file:
        content = [x.rstrip() for x in in_file.readlines()]

        with open(output_filename, "w") as out_file:
            repl_count = 0
            for line in content:
                if "run robintest" in line:
                    line = re.sub("""run robintest\\b""",
                        "run robintest --cover={}{}.covout".format(
                            input_filename, repl_count), line)
                    repl_count += 1
                elif "run robin" in line:
                    line = re.sub("""run robin\\b""",
                        "run robin.cover -test.coverprofile={}{}.covout".format(
                            input_filename, repl_count), line)
                    line = re.sub("""--help\\b""", "__bypass--help", line)
                    line = re.sub("""-h\\b""", "__bypass-h", line)
                    line = re.sub("""--version\\b""", "__bypass--version", line)
                    repl_count += 1
                elif '[ "$status" -ne 0 ]' in line:
                    line = re.sub("""\[ "\$status" -ne 0 \]""",
                        '[ "$status" -eq 0 ]', line)
                out_file.write("{}\n".format(line))

##########
# Script #
##########

# Input files definition
ROBINTEST_FILES = ["batsched_fast.bats",
                   "batsched_timeout.bats",
                   "nosched.bats"]
ROBIN_FILES = ["robin_cli.bats"]

for robintest_file in ROBINTEST_FILES + ROBIN_FILES:
    generate_bats_file(robintest_file, 
        re.sub(".bats", "-cover.bats", robintest_file))

