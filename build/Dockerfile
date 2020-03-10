FROM quay.io/ouzi/go-builder:1.14.0 as builder

WORKDIR /credstash-operator

# download modules
COPY go.mod .
COPY go.sum .
RUN go mod download

# add makefile
COPY Makefile .
RUN make setup

COPY . /credstash-operator

RUN make test
RUN TARGETS=linux/amd64 make build

FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/credstash-operator \
    USER_UID=1001 \
    USER_NAME=credstash-operator

# install operator binary
COPY --from=builder /credstash-operator/_dist/linux-amd64/manager/credstash-operator ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
