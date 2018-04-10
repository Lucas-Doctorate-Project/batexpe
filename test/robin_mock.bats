teardown() {
    rm -f ./robin ./robin.cover ./wc
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
