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
    killall batsim robin robin.cover 2>/dev/null || true
    chmod +x $(realpath .)
}

# teardown is called after each test
teardown() {
    not_running batsim
    not_running robin
    rm -rf ./unwritable-dir
}

@test "nosched-badbash" {
    run robintest batsim_nosched_badbash.yaml --test-timeout 30 \
                  --expect-robin-failure --expect-batsim-failure \
                  --expect-no-sched ${RT_CLEAN_CTX}
    good_return_or_print
}

@test "nosched-badinput" {
    run robintest batsim_nosched_badinput.yaml --test-timeout 30 \
                  --expect-robin-failure --expect-batsim-failure \
                  --expect-no-sched ${RT_CLEAN_CTX}
    good_return_or_print
}

@test "nosched-timeout" {
    run robintest batsim_nosched_timeout.yaml --test-timeout 5 \
                  --expect-robin-failure --expect-no-sched ${RT_CLEAN_CTX}
    good_return_or_print
}

@test "nosched-ok" {
    run robintest batsim_nosched_ok.yaml --test-timeout 30 \
                  --expect-robin-success --expect-batsim-success \
                  --expect-no-sched ${RT_CLEAN_CTX}
    good_return_or_print
}

@test "nosched-ok-alt" {
    run robintest batsim_nosched_ok_alt.yaml --test-timeout 30 \
                  --expect-robin-success --expect-batsim-success \
                  --expect-no-sched ${RT_CLEAN_CTX}
    good_return_or_print
}

@test "nosched-ok-mismatching-dir" {
    run robintest batsim_nosched_ok_mismatching_dir.yaml --test-timeout 30 \
                  --expect-robin-success --expect-batsim-success \
                  --expect-no-sched ${RT_CLEAN_CTX}
    good_return_or_print
}

@test "nosched-ok-long" {
    run robintest batsim_nosched_ok.yaml --test-timeout 30 \
                  --expect-robin-success --expect-batsim-success \
                  --expect-no-sched ${RT_CLEAN_CTX}
    good_return_or_print
}

@test "robintest.cover-nevercover" {
    ln -f -s $(which robintest.cover) ./myrobintest

    run ./myrobintest -test.coverprofile=nevercover.covout batsim_nosched_ok.yaml --test-timeout 30
}

@test "nosched-robin-cannot-write-files" {
    mkdir -p unwritable-dir/
    touch unwritable-dir/cmd

    run robintest invalid-desc-files/unwritable_nosched.yaml \
                  --test-timeout 30 \
                  --expect-robin-failure ${RT_CLEAN_CTX}
    good_return_or_print
}
