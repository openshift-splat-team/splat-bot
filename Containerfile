FROM registry.ci.openshift.org/ocp/4.16:base-rhel9

USER root

WORKDIR /usr/src/app

COPY . .
RUN go mod tidy && go mod vendor

RUN ./hack/build.sh

CMD ["./slack-bot" , "--slack-token-path", "/creds/token", "--slack-signing-secret-path", "/creds/secret"]
