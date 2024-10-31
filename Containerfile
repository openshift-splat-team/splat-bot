FROM golang:1.22


ENV HOME="/usr/src/app"
WORKDIR ${HOME}


RUN cat /etc/os-release
RUN apt update
RUN apt-get install -y make build-essential libssl-dev zlib1g-dev \
    libbz2-dev libreadline-dev libsqlite3-dev wget curl llvm libncurses5-dev \
    libncursesw5-dev xz-utils tk-dev libffi-dev liblzma-dev python3-openssl libpcre3 libpcre3-dev git


ENV PYENV_ROOT="${HOME}/.pyenv"
ENV PATH="${PYENV_ROOT}/shims:${PYENV_ROOT}/bin:${PATH}"

RUN git clone --depth=1 https://github.com/pyenv/pyenv.git .pyenv


ENV PYTHON_VERSION=3.12

RUN pyenv install ${PYTHON_VERSION}
RUN pyenv global ${PYTHON_VERSION}
RUN pip install ansible pyvmomi pyvim requests omsdk && \
    ansible-galaxy collection install ansible.posix && \
    ansible-galaxy collection install community.vmware


COPY . .

RUN ./hack/build.sh

CMD ["./slack-bot" , "--slack-token-path", "/creds/token", "--slack-signing-secret-path", "/creds/secret"]
