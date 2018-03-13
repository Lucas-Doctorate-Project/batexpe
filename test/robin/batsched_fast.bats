not_running() {
    set +e
    nb_running=$(ps -e -o command| cut -d' ' -f1| grep -E "\b$1$"| wc -l)
    set -e

    [ "${nb_running}" -eq 0 ]
}

# setup is called before each test
setup() {
    export RT_CLEAN_CTX="--expect-ctx-clean --expect-ctx-clean-at-begin --expect-ctx-clean-at-end"
    killall batsim robin batsched 2>/dev/null || true
}

# teardown is called after each test
teardown() {
    not_running batsim
    not_running robin
    not_running batsched
}

@test "batsched-ok" {
    run robintest batsched_ok.yaml --test-timeout 30 \
                  --expect-robin-success --expect-batsim-success \
                  --expect-sched-success ${RT_CLEAN_CTX}
    [ "$status" -eq 0 ]
}

@test "batsched-badwritedir" {
    run robintest batsched_badwritedir.yaml --test-timeout 30 \
                  --expect-robin-failure --expect-batsim-failure \
                  --expect-sched-killed ${RT_CLEAN_CTX}
    [ "$status" -eq 0 ]
}

@test "batsched-timeout" {
    run robintest batsched_timeout.yaml --test-timeout 5 \
                  --expect-robin-failure ${RT_CLEAN_CTX}
    [ "$status" -eq 0 ]
}

@test "batsched-schedcrash-badargs" {
    run robintest batsched_schedcrash_badargs.yaml --test-timeout 30 \
                  --expect-robin-failure --expect-batsim-killed \
                  --expect-sched-failure ${RT_CLEAN_CTX}
    [ "$status" -eq 0 ]
}

@test "batsched-schedcrash-begin-segfault" {
    run robintest batsched_schedcrash_begin_segfault.yaml --test-timeout 30 \
                  --expect-robin-failure --expect-batsim-killed \
                  --expect-sched-failure ${RT_CLEAN_CTX}
    [ "$status" -eq 0 ]
}

@test "batsched-schedcrash-mid-segfault" {
    run robintest batsched_schedcrash_mid_segfault.yaml --test-timeout 30 \
                  --expect-robin-failure --expect-batsim-killed \
                  --expect-sched-failure ${RT_CLEAN_CTX}
    [ "$status" -eq 0 ]
}

@test "batsched-schedcrash-end-segfault" {
    run robintest batsched-schedcrash-end-segfault.yaml --test-timeout 30 \
                  --expect-robin-failure --expect-batsim-killed \
                  --expect-sched-failure ${RT_CLEAN_CTX}
    [ "$status" -eq 0 ]
}

