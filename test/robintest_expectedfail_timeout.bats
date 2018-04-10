setup() {
    killall batsim robin robin.cover batsched 2>/dev/null || true
}

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

# Context expectation
@test "robintest-efail-context-clean-timeout-schedinuse" {
    batsched 1>/dev/null 2>/dev/null 3>/dev/null &
    run robintest batsched_ok.yaml --test-timeout=10 --expect-ctx-clean
    [ "$status" -ne 0 ]
    [[ "${lines[0]}" =~ "Unexpected context cleanliness during robin's execution" ]]
    killall batsched >/dev/null || true
}

@test "robintest-efail-context-clean-timeout-batsiminuse" {
    batsim -p ${BATSIM_DIR}/platforms/small_platform.xml -w ${BATSIM_DIR}/workload_profiles/test_workload_profile.json -e /tmp/robin/batsched_ok/out 1>/dev/null 2>/dev/null 3>/dev/null &
    run robintest batsched_ok.yaml --test-timeout=10 --expect-ctx-clean
    [ "$status" -ne 0 ]
    [[ "${lines[0]}" =~ "Unexpected context cleanliness during robin's execution" ]]
    killall batsim >/dev/null || true
}
