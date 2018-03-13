not_running() {
    set +e
    nb_running=$(ps -e -o command| cut -d' ' -f1| grep -E "\b$1$"| wc -l)
    set -e

    [ "${nb_running}" -eq 0 ]
}

# setup is called before each test
setup() {
    export RT_CLEAN_CTX="--expect-ctx-clean --expect-ctx-clean-at-begin --expect-ctx-clean-at-end"
    killall batsim robin 2>/dev/null || true
}

# teardown is called after each test
teardown() {
    not_running batsim
    not_running robin
}

@test "nosched-badbash" {
    run robintest batsim_nosched_badbash.yaml --test-timeout 30 \
                  --expect-robin-failure --expect-batsim-failure \
                  --expect-no-sched ${RT_CLEAN_CTX}
    [ "$status" -eq 0 ]
}

@test "nosched-badinput" {
    run robintest batsim_nosched_badinput.yaml --test-timeout 30 \
                  --expect-robin-failure --expect-batsim-failure \
                  --expect-no-sched ${RT_CLEAN_CTX}
    [ "$status" -eq 0 ]
}

@test "nosched-timeout" {
    run robintest batsim_nosched_timeout.yaml --test-timeout 5 \
                  --expect-robin-failure --expect-no-sched ${RT_CLEAN_CTX}
    [ "$status" -eq 0 ]
}

@test "nosched-ok" {
    run robintest batsim_nosched_ok.yaml --test-timeout 30 \
                  --expect-robin-success --expect-batsim-success \
                  --expect-no-sched ${RT_CLEAN_CTX}
    [ "$status" -eq 0 ]
}


