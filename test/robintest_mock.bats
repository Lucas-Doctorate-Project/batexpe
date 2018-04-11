teardown() {
    rm -f ./robin ./robin.cover ./ps
    killall fake-batsim 2>/dev/null
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

@test "robintest-mock-robin-unlaunchable" {
    echo '#!/this/sha/bang/should/be/invalid' > ./robin
    chmod +x ./robin
    echo '#!/this/sha/bang/should/be/invalid' > ./robin.cover
    chmod +x ./robin.cover

    run robintest batsim_nosched_ok.yaml --test-timeout 10 \
                  --expect-robin-failure
    [[ "${lines[0]}" =~ 'Could not start robin' ]]
}

@test "robintest-mock-robin-hello-check-badbatcmd" {
    ln -f -s $(realpath ./commands/hello) ./robin
    ln -f -s $(realpath ./commands/hello) ./robin.cover

    run robintest batsim_nosched_badcmd.yaml --test-timeout 10 \
                  --result-check-script=./commands/success
    [ "$status" -ne 0 ]
}

@test "robintest-mock-robin-hello-check-nodescfile" {
    ln -f -s $(realpath ./commands/hello) ./robin
    ln -f -s $(realpath ./commands/hello) ./robin.cover

    run robintest /this/file/does/not/exist --test-timeout 10 \
                  --result-check-script=./commands/success
    [ "$status" -ne 0 ]
}

@test "robintest-mock-robin-hello-check-descfile-badoutdir" {
    ln -f -s $(realpath ./commands/hello) ./robin
    ln -f -s $(realpath ./commands/hello) ./robin.cover

    run robintest ./invalid-desc-files/unreachable_outdir.yaml \
                  --test-timeout 10 \
                  --result-check-script=./commands/success
    [ "$status" -ne 0 ]
}

@test "robintest-mock-ps-failure" {
    ln -f -s $(realpath ./commands/failure) ./ps

    run robintest batsched_ok.yaml --test-timeout 10
    [ "$status" -ne 0 ]
}

@test "robintest-mock-robin-success-expectbatsimfailure" {
    ln -f -s $(realpath ./commands/success) ./robin
    ln -f -s $(realpath ./commands/success) ./robin.cover

    run robintest batsim_nosched_ok.yaml --test-timeout 10 \
                  --expect-batsim-failure
    [ "$status" -eq 0 ]
}

@test "robintest-mock-robin.cover-fake-badreturncode" {
    ln -f -s $(realpath ./commands/fakerobin.cover-success-nonint-return) ./robin.cover

    run robintest batsim_nosched_ok.yaml --test-timeout 10 \
                  --expect-robin-failure
    [ "$status" -ne 0 ]
}

@test "robintest-mock-batsim-batsched-background" {
    cp -f $(realpath $(which batsched)) ./fake-batsim
    ./fake-batsim 3>/dev/null &

    run robintest batsched_ok.yaml --test-timeout 10 \
                  --expect-robin-failure
    [ "$status" -ne 0 ]
}
