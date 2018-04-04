# Robin expectation
@test "robintest-efail-robin-success" {
    run robintest /this/file/should/not/exist.yaml --test-timeout=10 --expect-robin-success
    [ "$status" -ne 0 ]
    [[ "${lines[0]}" =~ 'Unexpected robin success state' ]]
}

@test "robintest-efail-robin-failure" {
    run robintest batsim_nosched_ok.yaml --test-timeout=10 --expect-robin-failure
    [ "$status" -ne 0 ]
    [[ "${lines[0]}" =~ 'Unexpected robin success state' ]]
}

@test "robintest-efail-robin-killed" {
    run robintest batsim_nosched_ok.yaml --test-timeout=10 --expect-robin-killed
    [ "$status" -ne 0 ]
    [[ "${lines[0]}" =~ 'Unexpected robin success state' ]]
    [[ "${lines[1]}" =~ 'Unexpected robin kill state' ]]
}

# Batsim expectation
@test "robintest-efail-batsim-success" {
    run robintest batsim_nosched_badbash.yaml --test-timeout=10 --expect-batsim-success
    [ "$status" -ne 0 ]
    [[ "${lines[0]}" =~ 'Unexpected batsim success state' ]]
}

@test "robintest-efail-batsim-failure" {
    run robintest batsim_nosched_ok.yaml --test-timeout=10 --expect-batsim-failure
    [ "$status" -ne 0 ]
    [[ "${lines[0]}" =~ 'Unexpected batsim success state' ]]
}

@test "robintest-efail-batsim-killed" {
    run robintest batsim_nosched_ok.yaml --test-timeout=10 --expect-batsim-killed
    [ "$status" -ne 0 ]
    [[ "${lines[0]}" =~ 'Unexpected batsim success state' ]]
    [[ "${lines[1]}" =~ 'Unexpected batsim kill state' ]]
}

# Sched expectation
@test "robintest-efail-sched-success" {
    run robintest batsim_badsched_nosuchcmd.yaml --test-timeout=10 --expect-sched-success
    [ "$status" -ne 0 ]
    [[ "${lines[0]}" =~ 'Unexpected sched success state' ]]
}

@test "robintest-efail-sched-failure" {
    run robintest batsched_ok.yaml --test-timeout=10 --expect-sched-failure
    [ "$status" -ne 0 ]
    [[ "${lines[0]}" =~ 'Unexpected sched success state' ]]
}

@test "robintest-efail-sched-killed" {
    run robintest batsched_ok.yaml --test-timeout=10 --expect-sched-killed
    [ "$status" -ne 0 ]
    [[ "${lines[0]}" =~ 'Unexpected sched success state' ]]
    [[ "${lines[1]}" =~ 'Unexpected sched kill state' ]]
}

@test "robintest-efail-sched-nosched" {
    run robintest batsched_ok.yaml --test-timeout=10 --expect-no-sched
    [ "$status" -ne 0 ]
    [[ "${lines[0]}" =~ 'Unexpected sched success state' ]]
    [[ "${lines[1]}" =~ 'Unexpected sched presence state' ]]
}
