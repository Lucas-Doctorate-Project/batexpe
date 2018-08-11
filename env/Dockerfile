# Base image with Nix installed: https://github.com/LnL7/nix-docker
FROM lnl7/nix:1.11.16

# Survival kit
RUN nix-env -i git gnutar gzip

# Retrieve datamove nix packages
RUN git clone https://gitlab.inria.fr/vreis/datamove-nix.git /datamovepkgs

# batsim + batsched scheduler
RUN nix-env -f /datamovepkgs -i batsim batsched-v1.1.0

# Go
RUN nix-env -i go

# Retrieve batsim directory (workloads/platforms...)
RUN git clone --branch v2.0.0 --depth 1 --single-branch \
              https://gitlab.inria.fr/batsim/batsim.git /batsim

# Robin dependencies
RUN nix-env -i procps iproute2

# Test and CI dependencies
RUN nix-env -i gnugrep psmisc bats findutils gnused gawk python3
RUN go get github.com/wadey/gocovmerge
RUN go get github.com/axw/gocov/gocov
RUN go get github.com/sahilm/yamldiff
