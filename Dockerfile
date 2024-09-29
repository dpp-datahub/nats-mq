# Stage 1: Build the binary
FROM golang:1.12.4 AS builder

LABEL maintainer="Stephen Asbury <sasbury@nats.io>"
LABEL "ProductName"="NATS-MQ Bridge" \
      "ProductVersion"="0.5"

# Install the MQ client from the Redistributable package
RUN mkdir -p /opt/mqm && cd /opt/mqm \
 && curl -LO "https://public.dhe.ibm.com/ibmdl/export/pub/software/websphere/messaging/mqdev/redist/9.2.0.0-IBM-MQC-Redist-LinuxX64.tar.gz" \
 && tar -zxf ./*.tar.gz \
 && rm -f ./*.tar.gz

ENV CGO_CFLAGS="-I/opt/mqm/inc/"
ENV CGO_LDFLAGS_ALLOW="-Wl,-rpath.*"

# Build the go ibmmq library
RUN go get github.com/ibm-messaging/mq-golang/ibmmq
RUN chmod -R a+rx $GOPATH/src
RUN cd $GOPATH/src \
  && go install github.com/ibm-messaging/mq-golang/ibmmq

# Copy and build the nats-mq code
RUN mkdir -p /nats-mq \
  && chmod -R 777 /nats-mq
COPY . /nats-mq
RUN rm -rf /nats-mq/build /nats-mq/.vscode
RUN chmod -R a+rx /nats-mq

RUN cd /nats-mq && go mod download && go install ./...

# Stage 2: Create the final image
FROM debian:bullseye-slim

# Install necessary dependencies
RUN apt-get update && apt-get install -y libstdc++6 && rm -rf /var/lib/apt/lists/*

# Create a non-root user and set a working directory
RUN useradd -m -s /bin/bash natsmq

# Add directories that are expected by MQ client
RUN mkdir -p /IBM/MQ/data/errors \
  && mkdir -p /.mqm \
  && chmod -R 777 /IBM \
  && chmod -R 777 /.mqm

WORKDIR /home/natsmq

# Copy the nats-mq binary and MQ libraries from the builder stage
COPY --from=builder /go/bin/nats-mq /usr/local/bin/nats-mq
COPY --from=builder /opt/mqm/lib64/* /opt/mqm/lib64
COPY --from=builder /opt/mqm/samp/ccsid_part2.tbl /opt/mqm/samp/ccsid_part2.tbl

# COPY --from=builder /opt/mqm/lib64/libmqm_r.so /usr/local/lib64/libmqm_r.so
# COPY --from=builder /opt/mqm/lib64/libmqmcs_r.so /usr/local/lib64/libmqmcs_r.so
# COPY --from=builder /opt/mqm/lib64/libmqz_r.so /usr/local/lib64/libmqz_r.so
# COPY --from=builder /opt/mqm/lib64/libmqe_r.so /usr/local/lib64/libmqe_r.so

# Set the library path
ENV LD_LIBRARY_PATH=/opt/mqm/lib64

# Change ownership of the binary and libraries to the non-root user
RUN chown -R natsmq:natsmq /usr/local/bin/nats-mq /opt/mqm/

# Switch to the non-root user
USER natsmq

# Run the bridge
ENTRYPOINT ["/usr/local/bin/nats-mq", "-c", "/mqbridge.conf"]
