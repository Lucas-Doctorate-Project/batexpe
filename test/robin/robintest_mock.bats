teardown() {
    rm -f ./robin ./robin.cover ./ss ./ps
    export PATH="${OLDPATH}"
}

setup() {
    export OLDPATH="${PATH}"
    export PATH=".:${PATH}"
}

@test "mock-robin-hello" {
    ln -f -s $(realpath ./commands/hello) ./robin
    ln -f -s $(realpath ./commands/hello) ./robin.cover

    run robintest batsim_nosched_ok.yaml --test-timeout 10
    [ "$status" -ne 0 ]
}

@test "mock-robin-hello-check-badbatcmd" {
    ln -f -s $(realpath ./commands/hello) ./robin
    ln -f -s $(realpath ./commands/hello) ./robin.cover

    run robintest batsim_nosched_badcmd.yaml --test-timeout 10 \
                  --result-check-script=./commands/success
    [ "$status" -ne 0 ]
}

@test "mock-robin-hello-check-nodescfile" {
    ln -f -s $(realpath ./commands/hello) ./robin
    ln -f -s $(realpath ./commands/hello) ./robin.cover

    run robintest /this/file/does/not/exist --test-timeout 10 \
                  --result-check-script=./commands/success
    [ "$status" -ne 0 ]
}

@test "mock-robin-hello-check-descfile-badoutdir" {
    ln -f -s $(realpath ./commands/hello) ./robin
    ln -f -s $(realpath ./commands/hello) ./robin.cover

    run robintest ./invalid-desc-files/unreachable-outdir.yaml \
                  --test-timeout 10 \
                  --result-check-script=./commands/success
    [ "$status" -ne 0 ]
}

@test "mock-ps-failure" {
    ln -f -s $(realpath ./commands/failure) ./ps

    run robintest batsched_ok.yaml --test-timeout 10
    [ "$status" -ne 0 ]
}
