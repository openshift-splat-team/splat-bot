FROM golang:1.22

WORKDIR /usr/src/app
RUN cat /etc/os-release
RUN apt update
RUN apt install -y
COPY . .
RUN ./hack/build.sh

CMD ["./slack-bot"]
 