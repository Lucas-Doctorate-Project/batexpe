# robin

robin is a program in charge of managing one Batsim simulation.

## Why?
A batsim simulation involves several processes that communicate with each other.

Executing these processes manually can be tricky, bothersome and risky,
as many problems may rise:
- One process got bad arguments and stopped. The other one is waiting forever.
- One process crashed. The other one is waiting forever.
- The simulation is launched but another one is already running.
- ...

The main goal of this script is to facilitate the execution of one single
simulation instance so it can be seen as a black box with inputs and outputs.

## Features
- Avoid to hinder other simulations. Does not execute the simulation if:
  - The communication socket is in use.
  - Another Batsim instance is running on the desired socket.
- Robust termination:
  - The simulation is stopped after a used-specified ``simulation-timeout``.  
    This allows to detect problems such as infinite loops,
    infinite scheduling cycles...
  - If one process fails (non-zero return / could not be executed at all),
    the second is stopped after a user-specified ``failure-timeout``.  
    This allows to stop the simulation when errors occur in any process.
  - When the first process finishes, the second is stopped after a
    user-specified ``success-timeout``.  
    This allows to detect bad scheduler (or rarely Batsim) termination.
  - Robin's exit code is 0 if and only if the simulation has been executed
    and has completed successfully.
- Cleanup:
  - *stopping* a process means *killing it and all its subprocesses*.
- Support no-scheduler mode (Batsim's ``--batexec`` option).
- Create executable command files (that can be hacked for painless debugging).
- Log the outputs of the involved processes.

## How does it work?
The main idea behind Robin is shown on the workflow below.
![robin main idea](automata/smooth2.svg "Robin main idea")
