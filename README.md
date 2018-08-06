[![pipeline status](https://gitlab.inria.fr/batsim/batexpe/badges/master/pipeline.svg)](https://gitlab.inria.fr/batsim/batexpe/commits/master)
[![coverage report](https://gitlab.inria.fr/batsim/batexpe/badges/master/coverage.svg)](https://gitlab.inria.fr/batsim/batexpe/commits/master)

This repository contains a set of Go tools around
[Batsim](https://gitlab.inria.fr/batsim/batsim) to simplify experiments.

## Install
### Via the go tool
```bash
go get gitlab.inria.fr/batsim/batexpe/cmd/robin
go get gitlab.inria.fr/batsim/batexpe/cmd/robintest
```

### Via nix
```bash
git clone https://gitlab.inria.fr/vreis/datamove-nix.git ./datamovepkgs
nix-env --file ./datamovepkgs --install --attr batexpe
```

## Proposed tools
- [robin](doc/robin.md) manages the execution of **one** simulation.  
  It is meant to be as robust as possible, as it is the core building block
  to create experiment workflows with Batsim.
- *robintest* is a *robin* wrapper mainly used to test robin.
  *robintest* notably allows to specify what (robin/batsim/scheduler)
  result is expected.
- the multiple commands are just wrappers around the *batexpe* library
  (written in Go).  
  This allows users to build their own tools (in Go) with decent code reuse.
