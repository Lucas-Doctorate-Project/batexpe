not_running() {
    set +e
    nb_running=$(ps -e -o command| cut -d' ' -f1| grep -E "\b$1$"| wc -l)
    set -e

    if [ "${nb_running}" -ne 0 ]; then
        (>&2 echo "A '$1' process is still running")
        return 1
    fi
}

good_return_or_print() {
    if [ "${status}" -ne 0 ]; then
        (>&2 echo "${output}")
        return 1
    fi
}

# setup is called before each test
setup() {
    export RT_CLEAN_CTX="--expect-ctx-clean --expect-ctx-clean-at-begin --expect-ctx-clean-at-end"
    killall batsim robin robin.cover batsched 2>/dev/null || true
}

# teardown is called after each test
teardown() {
    not_running batsim
    not_running robin
    not_running batsched
}

@test "batsched-badport" {
    run robintest batsched_badport.yaml --test-timeout 20 \
                  --expect-robin-failure --expect-batsim-killed \
                  --expect-sched-killed ${RT_CLEAN_CTX}
    good_return_or_print
}

@test "batsched-schedcrash-begin-loop" {
    run robintest batsched_schedcrash_begin_loop.yaml --test-timeout 20 \
                  --expect-robin-failure --expect-batsim-killed \
                  --expect-sched-killed ${RT_CLEAN_CTX}
    good_return_or_print
}

@test "batsched-schedcrash-mid-loop" {
    run robintest batsched_schedcrash_mid_loop.yaml --test-timeout 10 \
                  --expect-robin-killed --expect-batsim-killed \
                  --expect-sched-killed ${RT_CLEAN_CTX}
    good_return_or_print
}

@test "batsched-schedcrash-end-loop" {
    run robintest batsched_schedcrash_end_loop.yaml --test-timeout 10 \
                  --expect-robin-killed --expect-batsim-killed \
                  --expect-sched-killed ${RT_CLEAN_CTX}
    good_return_or_print
}
