teardown() {
    rm -f ./robin ./robin.cover ./wc ./head ./tail ./ss
    export PATH="${OLDPATH}"
}

setup() {
    export OLDPATH="${PATH}"
    export PATH=".:${PATH}"
}

@test "robin-mock-wc-failure-preview" {
    ln -f -s $(realpath ./commands/failure) ./wc

    run robin batsim_nosched_badinput.yaml --preview-on-error
    [ "$status" -ne 0 ]
}

@test "robin-mock-wc-hello-preview" {
    ln -f -s $(realpath ./commands/hello) ./wc

    run robin batsim_nosched_badinput.yaml --preview-on-error
    [ "$status" -ne 0 ]
}

@test "robin-mock-head-failure-preview" {
    ln -f -s $(realpath ./commands/failure) ./head

    run robin batsched_schedcrash_end_segfault_long.yaml --preview-on-error
    [ "$status" -ne 0 ]
}

@test "robin-mock-tail-failure-preview" {
    ln -f -s $(realpath ./commands/failure) ./tail

    run robin batsched_schedcrash_end_segfault_long.yaml --preview-on-error
    [ "$status" -ne 0 ]
}

@test "robintest-mock-ss-failure" {
    ln -f -s $(realpath ./commands/failure) ./ss

    run robin batsched_ok.yaml
    [ "$status" -ne 0 ]
}
