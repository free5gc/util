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

FROM golang:1.18.0-stretch AS dbtestapp
RUN apt-get update && apt-get -y install vim 

RUN mkdir -p /dbtestapp/bin
COPY --from=test $GOPATH/bin/* /dbtestapp/bin/
WORKDIR /dbtestapp
