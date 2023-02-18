#!/bin/bash

# Copyright 2023 Scott M. Long
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# 	http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

ARG DEV=0

FROM golang:1.19 as builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go install ./cmd/mu-app

FROM scratch as image-0
COPY --from=builder /go/bin/mu-app /mu-app
ENTRYPOINT ["/mu-app"]

FROM alpine:3.17.2 as image-1
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/bin/mu-app /mu-app
ENTRYPOINT ["/mu-app"]

FROM image-$DEV as image
