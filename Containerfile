FROM golang:1.21

WORKDIR /usr/src/app

COPY . .
RUN go mod tidy && go mod vendor

RUN ./hack/build.sh

CMD ["./slack-bot" , "--slack-token-path", "/creds/token", "--slack-signing-secret-path", "/creds/secret"]
