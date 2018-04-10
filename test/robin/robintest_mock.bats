teardown() {
    rm -f ./robin ./robin.cover
}

@test "robintest-mock-robin-hello" {
    ln -f -s $(realpath ./commands/hello) ./robin
    ln -f -s $(realpath ./commands/hello) ./robin.cover

    run PATH=".:${PATH}" robintest batsim_nosched_ok.yaml --test-timeout 10
    [ "$status" -ne 0 ]
}
