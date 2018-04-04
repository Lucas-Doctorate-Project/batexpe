@test "robintest-efail-batsim-success-timeout-crash" {
    run robintest batsim_badsched_wrongcmd.yaml --test-timeout=10 --expect-batsim-success
    [ "$status" -ne 0 ]
    [[ "${lines[0]}" =~ 'Unexpected batsim success state' ]]
    [[ "${lines[1]}" =~ 'Unexpected batsim kill state' ]]
}

@test "robintest-efail-batsim-failure-timeout-crash" {
    run robintest batsim_badsched_wrongcmd.yaml --test-timeout=10 --expect-batsim-success
    [ "$status" -ne 0 ]
    [[ "${lines[0]}" =~ 'Unexpected batsim success state' ]]
    [[ "${lines[1]}" =~ 'Unexpected batsim kill state' ]]
}

@test "robintest-efail-sched-success-timeout-loop" {
    run robintest batsched_schedcrash_begin_loop.yaml --test-timeout=10 --expect-sched-success
    [ "$status" -ne 0 ]
    [[ "${lines[0]}" =~ 'Test timeout reached' ]]
    [[ "${lines[1]}" =~ 'Unexpected sched success state' ]]
    [[ "${lines[2]}" =~ 'Unexpected sched kill state' ]]
}

@test "robintest-efail-sched-failure-timeout-loop" {
    run robintest batsched_schedcrash_begin_loop.yaml --test-timeout=10 --expect-sched-failure
    [ "$status" -ne 0 ]
    [[ "${lines[0]}" =~ 'Test timeout reached' ]]
    [[ "${lines[1]}" =~ 'Unexpected sched kill state' ]]
}
