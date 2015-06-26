# Pull base image.
FROM eris/base

# Set the env variables to non-interactive
ENV DEBIAN_FRONTEND noninteractive
ENV DEBIAN_PRIORITY critical
ENV DEBCONF_NOWARNINGS yes
ENV TERM linux
RUN echo 'debconf debconf/frontend select Noninteractive' | debconf-set-selections

# grab deps (gmp)
RUN apt-get update && \
  apt-get install -y --no-install-recommends \
    libgmp3-dev && \
  rm -rf /var/lib/apt/lists/*

# set the repo and install tendermint
ENV repo /go/src/github.com/eris-ltd/mint-client
ADD . $repo
WORKDIR $repo
RUN go install ./...

# set user
USER $USER
ENV TMROOT /home/eris/.eris/
WORKDIR /home/eris
CMD ["mintx"]
