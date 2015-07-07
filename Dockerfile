# Pull base image.
FROM eris/base

# grab deps (gmp)
RUN apt-get update && \
  apt-get install -y --no-install-recommends \
    libgmp3-dev jq && \
  rm -rf /var/lib/apt/lists/*

# set the repo and install tendermint
ENV repo /go/src/github.com/eris-ltd/mint-client
ADD . $repo
WORKDIR $repo
RUN go install ./...

# grab eris-keys
RUN go get github.com/eris-ltd/eris-keys

# grab tendermint
RUN go get github.com/tendermint/tendermint/cmd/tendermint
ENV TMROOT /home/eris/.eris/blockchains/tendermint
ADD config.toml $TMROOT/config.toml
RUN chown -R $USER:$USER $TMROOT

ADD ./test.sh /test.sh
RUN chown $USER:$USER /test.sh

# set user
USER $USER
WORKDIR /home/eris

CMD ["/test.sh"]
