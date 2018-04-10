@test "robintest-mock-robin-hello" {
    ln -f -s $(realpath ./commands/hello) ./robin

    run PATH=".:${PATH}" robintest batsim_nosched_ok.yaml --test-timeout 10
    [ "$status" -ne 0 ]
}
