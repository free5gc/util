# SPDX-FileCopyrightText: 2022-present Intel Corporation
# Copyright 2019-present Open Networking Foundation
#
# SPDX-License-Identifier: Apache-2.0
#

FROM golang:1.18.0-stretch AS test

LABEL maintainer="ONF <omec-dev@opennetworking.org>"

RUN apt-get update && apt-get -y install vim 

RUN cd $GOPATH/src && mkdir -p dbtestapp
COPY . $GOPATH/src/dbtestapp
RUN cd $GOPATH/src/dbtestapp && go install 

FROM alpine:3.16 AS dbtestapp
RUN apk update && apk add -U gcompat vim strace net-tools curl netcat-openbsd bind-tools bash

RUN mkdir -p /dbtestapp/bin
COPY --from=test /go/bin/* /dbtestapp/bin/
WORKDIR /dbtestapp
