@test "cli-robintest-noargs" {
    run robintest
    [ "$status" -ne 0 ]
}

@test "cli-robintest-help" {
    run robintest --help
    [ "$status" -eq 0 ]
}

@test "cli-robintest-h" {
    run robintest -h
    [ "$status" -eq 0 ]
}

@test "cli-robintest-version" {
    run robintest --version
    [ "$status" -eq 0 ]
    [ $(echo "${lines[0]}" | grep -o -E 'v[0-9]+\.[0-9]+\.[0-9]+.*') == "${lines[0]}" ]
}

@test "cli-robintest-nonexistent-desc-file" {
    run robintest /this/file/should/not/exist.yaml --test-timeout=10 \
                  --expect-robin-failure
    [ "$status" -eq 0 ]
}

@test "cli-robintest-bad-test-timeout" {
    run robintest batsim_nosched_ok.yaml --test-timeout='not a valid simeout'
    [ "$status" -ne 0 ]
    [[ "${lines[0]}" =~ 'Invalid test timeout' ]]
}

@test "cli-robintest-ok-debug" {
    run robintest batsim_nosched_ok.yaml --test-timeout=10 --debug
    [ "$status" -eq 0 ]
}

@test "cli-robintest-resultcheckscript-success" {
    run robintest batsim_nosched_ok.yaml --test-timeout=10 --result-check-script=./checkscript_success.bash
    [ "$status" -eq 0 ]
}

@test "cli-robintest-resultcheckscript-failure" {
    run robintest batsim_nosched_ok.yaml --test-timeout=10 --result-check-script=./checkscript_failure.bash
    [ "$status" -ne 0 ]
    [[ "${lines[0]}" =~ 'Check subprocess failed' ]]
}

@test "cli-robintest-resultcheckscript-badfile" {
    run robintest batsim_nosched_ok.yaml --test-timeout=10 --result-check-script=/does/not/exist.bash
    [ "$status" -ne 0 ]
    [[ "${lines[0]}" =~ 'Could not start Check subprocess' ]]
}
