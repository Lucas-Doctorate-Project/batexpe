#!/usr/bin/env python3
import re

#############
# Functions #
#############

def generate_bats_file(input_filename, output_filename):
    options_to_bypass = ['--help', '-h', '--version',
                         '--output-dir', '--batcmd', '--schedcmd',
                         '--simulation-timeout', '--ready-timeout',
                         '--success-timeout', '--failure-timeout']

    with open(input_filename, "r") as in_file:
        content = [x.rstrip() for x in in_file.readlines()]

        with open(output_filename, "w") as out_file:
            repl_count = 0
            for line in content:
                if "run robintest" in line:
                    line = re.sub("""run robintest\\b""",
                        "run robintest.cover "
                        "-test.coverprofile={f}.rt.{c}.covout "
                        "__bypass--cover={f}.r.{c}.covout".format(
                            f=input_filename, c=repl_count), line)
                    repl_count += 1
                elif "run robin" in line:
                    line = re.sub("""run robin\\b""",
                        "run robin.cover "
                        "-test.coverprofile={f}.r.{c}.covout".format(
                            f=input_filename, c=repl_count), line)
                    repl_count += 1
                elif '[ "$status" -ne 0 ]' in line:
                    line = re.sub("""\[ "\$status" -ne 0 \]""",
                        '[ "$status" -eq 0 ]', line)

                for option in options_to_bypass:
                    line = re.sub(option + '\\b',
                                  '__bypass' + option, line)

                out_file.write("{}\n".format(line))

##########
# Script #
##########

# Input files definition
ROBINTEST_FILES = ["batsched_fast.bats",
                   "batsched_timeout.bats",
                   "nosched.bats",
                   "robintest_cli.bats"]
ROBIN_FILES = ["robin_cli.bats"]

for robintest_file in ROBINTEST_FILES + ROBIN_FILES:
    generate_bats_file(robintest_file, re.sub(
        ".bats", "-cover.bats", robintest_file))
