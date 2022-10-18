# Copyright 2017 Heptio Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Build the manager binary
FROM amaas-eos-mw1.cec.lab.emc.com:5047/mw/golang:1.18.1 as builder

WORKDIR /workspace
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o kafka-producer

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM amaas-eos-mw1.cec.lab.emc.com:5047/mw/gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/kafka-producer .
COPY --from=builder /workspace/server.crt /
COPY --from=builder /workspace/client.crt /
COPY --from=builder /workspace/client.key /
USER nobody:nobody

CMD ["/kafka-producer"]
