teardown() {
    rm -f ./robin ./robin.cover
    export PATH="${OLDPATH}"
}

setup() {
    export OLDPATH="${PATH}"
    export PATH=".:${PATH}"
}

@test "robintest-mock-robin-hello" {
    ln -f -s $(realpath ./commands/hello) ./robin
    ln -f -s $(realpath ./commands/hello) ./robin.cover

    run robintest batsim_nosched_ok.yaml --test-timeout 10
    [ "$status" -ne 0 ]
}
