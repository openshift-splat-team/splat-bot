FROM registry.access.redhat.com/ubi8/go-toolset:1.20.12-2

USER root

WORKDIR /usr/src/app

COPY . .
RUN go mod tidy && go mod vendor

RUN ./hack/build.sh

CMD ["./slack-bot" , "--slack-token-path", "/creds/token", "--slack-signing-secret-path", "/creds/secret"]