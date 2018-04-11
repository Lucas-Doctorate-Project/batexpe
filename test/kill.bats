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
    export kill_wait_time=0.1
    killall batsim robin robin.cover batsched 2>/dev/null || true
}

# teardown is called after each test
teardown() {
    not_running batsim
    not_running robin
    not_running batsched
}

###################
# Nosched, SIGINT #
###################
@test "kill-sigint-nosched-0.05" {
    "robin" batsim_nosched_ok_long.yaml 3>/dev/null &
    robin_pid=$!
    sleep 0.05
    kill -INT ${robin_pid} 3>/dev/null &
    sleep ${kill_wait_time}
}

@test "kill-sigint-nosched-0.06" {
    "robin" batsim_nosched_ok_long.yaml 3>/dev/null &
    robin_pid=$!
    sleep 0.06
    kill -INT ${robin_pid} 3>/dev/null &
    sleep ${kill_wait_time}
}

@test "kill-sigint-nosched-0.07" {
    "robin" batsim_nosched_ok_long.yaml 3>/dev/null &
    robin_pid=$!
    sleep 0.07
    kill -INT ${robin_pid} 3>/dev/null &
    sleep ${kill_wait_time}
}

@test "kill-sigint-nosched-0.08" {
    "robin" batsim_nosched_ok_long.yaml 3>/dev/null &
    robin_pid=$!
    sleep 0.08
    kill -INT ${robin_pid} 3>/dev/null &
    sleep ${kill_wait_time}
}

@test "kill-sigint-nosched-0.09" {
    "robin" batsim_nosched_ok_long.yaml 3>/dev/null &
    robin_pid=$!
    sleep 0.09
    kill -INT ${robin_pid} 3>/dev/null &
    sleep ${kill_wait_time}
}

@test "kill-sigint-nosched-0.10" {
    "robin" batsim_nosched_ok_long.yaml 3>/dev/null &
    robin_pid=$!
    sleep 0.10
    kill -TERM ${robin_pid} 3>/dev/null &
    sleep ${kill_wait_time}
}

####################
# Nosched, SIGTERM #
####################
@test "kill-sigterm-nosched-0.05" {
    "robin" batsim_nosched_ok_long.yaml 3>/dev/null &
    robin_pid=$!
    sleep 0.05
    kill -TERM ${robin_pid} 3>/dev/null &
    sleep ${kill_wait_time}
}

@test "kill-sigterm-nosched-0.06" {
    "robin" batsim_nosched_ok_long.yaml 3>/dev/null &
    robin_pid=$!
    sleep 0.06
    kill -TERM ${robin_pid} 3>/dev/null &
    sleep ${kill_wait_time}
}

@test "kill-sigterm-nosched-0.07" {
    "robin" batsim_nosched_ok_long.yaml 3>/dev/null &
    robin_pid=$!
    sleep 0.07
    kill -TERM ${robin_pid} 3>/dev/null &
    sleep ${kill_wait_time}
}

@test "kill-sigterm-nosched-0.08" {
    "robin" batsim_nosched_ok_long.yaml 3>/dev/null &
    robin_pid=$!
    sleep 0.08
    kill -TERM ${robin_pid} 3>/dev/null &
    sleep ${kill_wait_time}
}

@test "kill-sigterm-nosched-0.09" {
    "robin" batsim_nosched_ok_long.yaml 3>/dev/null &
    robin_pid=$!
    sleep 0.09
    kill -TERM ${robin_pid} 3>/dev/null &
    sleep ${kill_wait_time}
}

@test "kill-sigterm-nosched-0.10" {
    "robin" batsim_nosched_ok_long.yaml 3>/dev/null &
    robin_pid=$!
    sleep 0.10
    kill -TERM ${robin_pid} 3>/dev/null &
    sleep ${kill_wait_time}
}

####################
# Batsched, SIGINT #
####################
@test "kill-sigint-batsched-0.05" {
    "robin" batsched_schedcrash_end_loop.yaml 3>/dev/null &
    robin_pid=$!
    sleep 0.05
    kill -INT ${robin_pid} 3>/dev/null &
    sleep ${kill_wait_time}
}

@test "kill-sigint-batsched-0.06" {
    "robin" batsched_schedcrash_end_loop.yaml 3>/dev/null &
    robin_pid=$!
    sleep 0.06
    kill -INT ${robin_pid} 3>/dev/null &
    sleep ${kill_wait_time}
}

@test "kill-sigint-batsched-0.07" {
    "robin" batsched_schedcrash_end_loop.yaml 3>/dev/null &
    robin_pid=$!
    sleep 0.07
    kill -INT ${robin_pid} 3>/dev/null &
    sleep ${kill_wait_time}
}

@test "kill-sigint-batsched-0.08" {
    "robin" batsched_schedcrash_end_loop.yaml 3>/dev/null &
    robin_pid=$!
    sleep 0.08
    kill -INT ${robin_pid} 3>/dev/null &
    sleep ${kill_wait_time}
}

@test "kill-sigint-batsched-0.09" {
    "robin" batsched_schedcrash_end_loop.yaml 3>/dev/null &
    robin_pid=$!
    sleep 0.09
    kill -INT ${robin_pid} 3>/dev/null &
    sleep ${kill_wait_time}
}

@test "kill-sigint-batsched-0.10" {
    "robin" batsched_schedcrash_end_loop.yaml 3>/dev/null &
    robin_pid=$!
    sleep 0.10
    kill -TERM ${robin_pid} 3>/dev/null &
    sleep ${kill_wait_time}
}

@test "kill-sigint-batsched-5.00" {
    "robin" batsched_schedcrash_end_loop.yaml 3>/dev/null &
    robin_pid=$!
    sleep 5.00
    kill -INT ${robin_pid}
    sleep 2.00
}

#####################
# Batsched, SIGTERM #
#####################
@test "kill-sigterm-batsched-0.05" {
    "robin" batsched_schedcrash_end_loop.yaml 3>/dev/null &
    robin_pid=$!
    sleep 0.05
    kill -TERM ${robin_pid} 3>/dev/null &
    sleep ${kill_wait_time}
}

@test "kill-sigterm-batsched-0.06" {
    "robin" batsched_schedcrash_end_loop.yaml 3>/dev/null &
    robin_pid=$!
    sleep 0.06
    kill -TERM ${robin_pid} 3>/dev/null &
    sleep ${kill_wait_time}
}

@test "kill-sigterm-batsched-0.07" {
    "robin" batsched_schedcrash_end_loop.yaml 3>/dev/null &
    robin_pid=$!
    sleep 0.07
    kill -TERM ${robin_pid} 3>/dev/null &
    sleep ${kill_wait_time}
}

@test "kill-sigterm-batsched-0.08" {
    "robin" batsched_schedcrash_end_loop.yaml 3>/dev/null &
    robin_pid=$!
    sleep 0.08
    kill -TERM ${robin_pid} 3>/dev/null &
    sleep ${kill_wait_time}
}

@test "kill-sigterm-batsched-0.09" {
    "robin" batsched_schedcrash_end_loop.yaml 3>/dev/null &
    robin_pid=$!
    sleep 0.09
    kill -TERM ${robin_pid} 3>/dev/null &
    sleep ${kill_wait_time}
}

@test "kill-sigterm-batsched-0.10" {
    "robin" batsched_schedcrash_end_loop.yaml 3>/dev/null &
    robin_pid=$!
    sleep 0.10
    kill -TERM ${robin_pid} 3>/dev/null &
    sleep ${kill_wait_time}
}

@test "kill-sigterm-batsched-5.00" {
    "robin" batsched_schedcrash_end_loop.yaml 3>/dev/null &
    robin_pid=$!
    sleep 5.00
    kill -TERM ${robin_pid}
}
