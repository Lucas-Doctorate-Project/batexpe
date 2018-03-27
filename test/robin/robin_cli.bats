@test "cli-robin-noargs" {
    run robin
    [ "$status" -ne 0 ]
}

@test "cli-robin-help" {
    run robin --help
    [ "$status" -eq 0 ]
}

@test "cli-robin-h" {
    run robin -h
    [ "$status" -eq 0 ]
}

@test "cli-robin-version" {
    run robin --version
    [ "$status" -eq 0 ]
    [ $(echo "${lines[0]}" | grep -o -E '^[0-9]+\.[0-9]+\.[0-9]+$') == "${lines[0]}" ]
}

@test "cli-robin-nonexistent-desc-file" {
    run robin /this/file/should/not/exist.yaml
    [ "$status" -ne 0 ]
    [[ "${lines[0]}" =~ 'Cannot open description file' ]]
}
