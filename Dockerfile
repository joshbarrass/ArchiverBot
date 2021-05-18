FROM golang:1.14-buster

# install packages
RUN  apt-get update \
  && apt-get install -y tini git python3 python3-pip ffmpeg \
  && rm -rf /var/lib/apt/lists/*

RUN echo $(python3 --version)

WORKDIR /code

RUN addgroup --gid 10001 --system nonroot && adduser -u 10000 --system --gid 10001 --home /home/nonroot nonroot

USER nonroot

COPY requirements.txt /code/requirements.txt

WORKDIR /home/nonroot
RUN  pip3 install -r /code/requirements.txt

ENV PATH /home/nonroot/.local/bin:$PATH

# confirm uarchiver is installed on PATH
RUN echo "Installed UArchiver version $(uarchiver --version)"

COPY cmd/ /code/ArchiverBot/cmd/
COPY internal/ /code/ArchiverBot/internal/
COPY pkg/ /code/ArchiverBot/pkg/
COPY go.mod /code/ArchiverBot/go.mod
COPY go.sum /code/ArchiverBot/go.sum
COPY LICENSE /code/ArchiverBot/LICENSE

# build bot
WORKDIR /code/ArchiverBot/
RUN GO111MODULES=on go build -o /home/nonroot/ArchiverBot ./cmd/archiverbot.go

WORKDIR /home/nonroot

ENTRYPOINT ["/usr/bin/tini", "--", "/home/nonroot/ArchiverBot"]
CMD []

ENV AB_OUT_DIR=/downloads
VOLUME /downloads