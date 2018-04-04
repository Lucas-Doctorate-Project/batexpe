@test "cli-robin-noargs" {
    run robin
    [ "$status" -ne 0 ]
}

@test "cli-robin-help" {
    run robin --help
    [ "$status" -eq 0 ]
}

@test "cli-robin-h" {
    run robin -h
    [ "$status" -eq 0 ]
}

@test "cli-robin-version" {
    run robin --version
    [ "$status" -eq 0 ]
    [ $(echo "${lines[0]}" | grep -o -E '^[0-9]+\.[0-9]+\.[0-9]+$') == "${lines[0]}" ]
}

@test "cli-robin-nonexistent-desc-file" {
    run robin /this/file/should/not/exist.yaml
    [ "$status" -ne 0 ]
    [[ "${lines[0]}" =~ 'Cannot open description file' ]]
}

# Verbosity tests
@test "cli-robin-ok-verbose" {
    run robin batsim_nosched_ok.yaml --verbose
    [ "$status" -eq 0 ]
}

@test "cli-robin-ok-quiet" {
    run robin batsim_nosched_ok.yaml --quiet
    [ "$status" -eq 0 ]
}

@test "cli-robin-ok-debug" {
    run robin batsim_nosched_ok.yaml --debug
    [ "$status" -eq 0 ]
}

@test "cli-robin-badbash-preview" {
    run robin batsim_nosched_badbash.yaml --preview-on-error
    [ "$status" -ne 0 ]
}

# generate subcommand test
@test "cli-robin-generate-ok-nosched" {
    run robin generate /tmp/robin_generated.yaml \
                       --output-dir='/tmp/robin/batsim_nosched_ok' \
                       --batcmd='batsim -p ${BATSIM_DIR}/platforms/small_platform.xml -w ${BATSIM_DIR}/workload_profiles/test_workload_profile.json -e /tmp/robin/batsim_nosched_ok/out --batexec' \
                       --schedcmd='' \
                       --simulation-timeout=30 \
                       --ready-timeout=5 \
                       --success-timeout=5 \
                       --failure-timeout=0
    [ "$status" -eq 0 ]

    run yamldiff --file1 /tmp/robin_generated.yaml --file2 batsim_nosched_ok.yaml
    echo ${lines}
    [[ "${lines}" = '' ]]
}

@test "cli-robin-generate-ok-batsched" {
    run robin generate /tmp/robin_generated.yaml \
                       --output-dir='/tmp/robin/batsched_ok' \
                       --batcmd='batsim -p ${BATSIM_DIR}/platforms/small_platform.xml -w ${BATSIM_DIR}/workload_profiles/test_workload_profile.json -e /tmp/robin/batsched_ok/out' \
                       --schedcmd='batsched' \
                       --simulation-timeout=30 \
                       --ready-timeout=5 \
                       --success-timeout=5 \
                       --failure-timeout=0
    [ "$status" -eq 0 ]

    run yamldiff --file1 /tmp/robin_generated.yaml --file2 batsched_ok.yaml
    echo ${lines}
    [[ "${lines}" = '' ]]
}
