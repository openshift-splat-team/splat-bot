FROM registry.ci.openshift.org/openshift/release:golang-1.18 AS builder
WORKDIR /go/src/github.com/openshift-eng/splat-sandbox
COPY . .
RUN go build cmd/slack-bot/slack-bot.go

FROM registry.ci.openshift.org/openshift/origin-v4.0:base
COPY --from=builder /go/src/github.com/openshift-eng/splat-sandbox/slack-bot /slack-bot
CMD ["/slack-bot" , "--slack-token-path", "/creds/bot-token", "--slack-signing-secret-path", "/creds/signing-secret"]