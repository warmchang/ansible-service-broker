FROM centos:7

RUN yum upgrade -y
RUN yum install -y git btrfs-progs-devel device-mapper-devel make gcc-c++

ENV GOPATH=/tmp/go GOBIN=/tmp/go/bin GOROOT=/usr/local/go
RUN mkdir -p /tmp/go/src && mkdir -p /tmp/go/bin && \
  mkdir -p /tmp/go/src/github.com/fusor && \
  mkdir -p /usr/local/ansible-service-broker/bin && \
  mkdir -p /etc/ansible-service-broker
ENV PATH=/usr/local/go/bin:/usr/local/ansible-service-broker/bin:$PATH
RUN curl -L "https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz" \
  > /usr/local/go.tar.gz && cd /usr/local && tar xf go.tar.gz
RUN curl -L "https://github.com/Masterminds/glide/releases/download/v0.12.3/glide-v0.12.3-linux-amd64.tar.gz" \
  > /tmp/glide.tar.gz && cd /tmp && tar xf glide.tar.gz --strip-components=1 && \
  mv glide /usr/bin

RUN git clone https://github.com/fusor/ansible-service-broker \
  /tmp/go/src/github.com/fusor/ansible-service-broker && \
  cd /tmp/go/src/github.com/fusor/ansible-service-broker && git checkout -b demo-broker origin/demo-broker
RUN cd /tmp/go/src/github.com/fusor/ansible-service-broker && glide install
RUN cd /tmp/go/src/github.com/fusor/ansible-service-broker && make build

RUN cp $GOBIN/broker /usr/bin/asbd

COPY docker/entrypoint.sh /usr/bin
COPY docker/ansible-service-broker /usr/local/ansible-service-broker/bin
COPY etc/ex.demo.config.yaml /etc/ansible-service-broker/config.yaml

ENTRYPOINT ["/usr/bin/entrypoint.sh"]
