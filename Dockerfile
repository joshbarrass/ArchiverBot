FROM golang:1.14-buster

# install packages
RUN  apt-get update \
  && apt-get install -y tini git python3 python3-pip \
  && rm -rf /var/lib/apt/lists/*

RUN echo $(python3 --version)

WORKDIR /code

COPY cmd/ /code/ArchiverBot/cmd/
COPY internal/ /code/ArchiverBot/internal/
COPY pkg/ /code/ArchiverBot/pkg/
COPY go.mod /code/ArchiverBot/go.mod
COPY go.sum /code/ArchiverBot/go.sum
COPY LICENSE /code/ArchiverBot/LICENSE

RUN addgroup --gid 10001 --system nonroot && adduser -u 10000 --system --gid 10001 --home /home/nonroot nonroot

USER nonroot

# build bot
WORKDIR /code/ArchiverBot/
RUN GO111MODULES=on go build -o /home/nonroot/ArchiverBot ./cmd/archiverbot.go 

ENTRYPOINT ["/sbin/tini", "--", "/home/nonroot/ArchiverBot"]

WORKDIR /home/nonroot
RUN  git clone https://github.com/joshbarrass/UArchiver \
  && pip3 install -r /home/nonroot/UArchiver/requirements.txt \
  && pip3 install /home/nonroot/UArchiver

ENV PATH /home/nonroot/.local/bin:$PATH

# confirm uarchiver is installed on PATH
RUN uarchiver --version

CMD []