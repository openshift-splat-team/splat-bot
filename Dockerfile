FROM registry.access.redhat.com/ubi8/go-toolset:1.20.12-2

USER root

WORKDIR /usr/src/app

COPY . .
RUN go mod tidy && go mod vendor

RUN ./hack/build.sh

CMD ["./slack-bot" , "--slack-token-path", "$SLACK_BOT_TOKEN", "--slack-signing-secret-path", "$SLACK_BOT_SECRET"]