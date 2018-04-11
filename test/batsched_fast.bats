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
    rm -rf ./unwritable-dir
}

@test "batsched-ok" {
    run robintest batsched_ok.yaml --test-timeout 30 \
                  --expect-robin-success --expect-batsim-success \
                  --expect-sched-success ${RT_CLEAN_CTX}
    good_return_or_print
}

@test "batsched-badwritedir" {
    run robintest batsched_badwritedir.yaml --test-timeout 30 \
                  --expect-robin-failure --expect-batsim-failure \
                  --expect-sched-killed ${RT_CLEAN_CTX}
    good_return_or_print
}

@test "batsched-timeout" {
    run robintest batsched_timeout.yaml --test-timeout 5 \
                  --expect-robin-failure ${RT_CLEAN_CTX}
    good_return_or_print
}

@test "batsched-schedcrash-badargs" {
    run robintest batsched_schedcrash_badargs.yaml --test-timeout 30 \
                  --expect-robin-failure --expect-batsim-killed \
                  --expect-sched-failure ${RT_CLEAN_CTX}
    good_return_or_print
}

@test "batsched-schedcrash-begin-segfault" {
    run robintest batsched_schedcrash_begin_segfault.yaml --test-timeout 30 \
                  --expect-robin-failure --expect-batsim-killed \
                  --expect-sched-failure ${RT_CLEAN_CTX}
    good_return_or_print
}

@test "batsched-schedcrash-mid-segfault" {
    run robintest batsched_schedcrash_mid_segfault.yaml --test-timeout 30 \
                  --expect-robin-failure --expect-batsim-killed \
                  --expect-sched-failure ${RT_CLEAN_CTX}
    good_return_or_print
}

@test "batsched-schedcrash-end-segfault" {
    run robintest batsched_schedcrash_end_segfault.yaml --test-timeout 30 \
                  --expect-robin-failure --expect-batsim-killed \
                  --expect-sched-failure ${RT_CLEAN_CTX}
    good_return_or_print
}

@test "batsched-schedcrash-end-segfault-long-preview" {
    run robin batsched_schedcrash_end_segfault_long.yaml --preview-on-error
    [ "$status" -ne 0 ]
}

@test "batsched-schedcrash-end-segfault-long-preview-quiet" {
    run robin batsched_schedcrash_end_segfault_long_quiet.yaml --preview-on-error
    [ "$status" -ne 0 ]
}

@test "batsched-robin-cannot-write-files" {
    mkdir -p unwritable-dir/cmd
    mkdir -p unwritable-dir/log

    chmod -w unwritable-dir/cmd
    chmod -w unwritable-dir/log

    run robintest invalid-desc-files/unwritable_batsched.yaml \
                  --test-timeout 30 \
                  --expect-robin-failure ${RT_CLEAN_CTX}
    good_return_or_print
}
