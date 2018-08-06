not_running() {
    set +e
    nb_running=$(ps -e -o command| cut -d' ' -f1| grep -E "\b$1$"| wc -l)
    set -e

    if [ "${nb_running}" -ne 0 ]; then
        (>&2 echo "A '$1' process is still running")
        return 1
    fi
}

# setup is called before each test
setup() {
    export RT_CLEAN_CTX="--expect-ctx-clean --expect-ctx-clean-at-begin --expect-ctx-clean-at-end"
    killall batsim robin robin.cover batsched 2>/dev/null || true
}

good_return_or_print() {
    if [ "${status}" -ne 0 ]; then
        (>&2 echo "${output}")
        return 1
    fi
}

# teardown is called after each test
teardown() {
    not_running batsim
    not_running robin
    not_running batsched
    rm -rf ./dir-without-write-permissions
}

@test "badinputfiles-unreachable-outdir" {
    run robintest invalid-desc-files/unreachable_outdir.yaml \
                  --test-timeout 30 \
                  --expect-robin-failure ${RT_CLEAN_CTX}
    good_return_or_print
}

@test "badinputfiles-unwritable-outdir" {
    mkdir dir-without-write-permissions
    chmod -w dir-without-write-permissions

    run robintest invalid-desc-files/unwritable_outdir.yaml \
                  --test-timeout 30 \
                  --expect-robin-failure ${RT_CLEAN_CTX}
    good_return_or_print
}

@test "badinputfiles-unsplittable-batcmd" {
    run robintest invalid-desc-files/unsplittable_batcmd.yaml \
                  --test-timeout 30 \
                  --expect-robin-failure ${RT_CLEAN_CTX}
    good_return_or_print
}

@test "badinputfiles-nonintegral-port" {
    run robintest invalid-desc-files/nonintegral_port.yaml \
                  --test-timeout 30 \
                  --expect-robin-failure ${RT_CLEAN_CTX}
    good_return_or_print
}

@test "badinputfiles-huge-port" {
    run robintest invalid-desc-files/huge_port.yaml \
                  --test-timeout 30 \
                  --expect-robin-failure ${RT_CLEAN_CTX}
    good_return_or_print
}

@test "badinputfiles-nonfloat-timeout" {
    run robintest invalid-desc-files/nonfloat_timeout.yaml \
                  --test-timeout 30 \
                  --expect-robin-failure ${RT_CLEAN_CTX}
    good_return_or_print
}

@test "badinputfiles-missing-timeout" {
    run robintest invalid-desc-files/missing_timeout.yaml \
                  --test-timeout 30 \
                  --expect-robin-failure ${RT_CLEAN_CTX}
    good_return_or_print
}

@test "badinputfiles-nonstring-batcmd" {
    run robintest invalid-desc-files/nonstring_batcmd.yaml \
                  --test-timeout 30 \
                  --expect-robin-failure ${RT_CLEAN_CTX}
    good_return_or_print
}

@test "badinputfiles-missing-batcmd" {
    run robintest invalid-desc-files/missing_batcmd.yaml \
                  --test-timeout 30 \
                  --expect-robin-failure ${RT_CLEAN_CTX}
    good_return_or_print
}

@test "badinputfiles-list" {
    run robintest invalid-desc-files/list.yaml \
                  --test-timeout 30 \
                  --expect-robin-failure ${RT_CLEAN_CTX}
    good_return_or_print
}

@test "badinputfiles-notyaml" {
    run robintest invalid-desc-files/notyaml.yaml \
                  --test-timeout 30 \
                  --expect-robin-failure ${RT_CLEAN_CTX}
    good_return_or_print
}

@test "badinputfiles-nosched-with-sched" {
    run robintest invalid-desc-files/nosched_with_sched.yaml \
                  --test-timeout 30 \
                  --expect-robin-failure ${RT_CLEAN_CTX}
    good_return_or_print
}

@test "badinputfiles-batsched-without-sched" {
    run robintest invalid-desc-files/batsched_without_sched.yaml \
                  --test-timeout 30 \
                  --expect-robin-failure ${RT_CLEAN_CTX}
    good_return_or_print
}
