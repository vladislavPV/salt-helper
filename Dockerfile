FROM centos

ENV GOLANG_VERSION 1.8.3
ENV goRelSha256 1862f4c3d3907e59b04a757cfda0ea7aa9ef39274af99a784f5be843c80c6772

##################Salt install###################################
RUN yum install -y https://repo.saltstack.com/yum/redhat/salt-repo-latest-2.el7.noarch.rpm
RUN yum install -y epel-release
RUN yum update -y
RUN yum install -y \
  file \
  htop \
  iproute \
  less \
  lsof \
  nano \
  net-tools \
  python-augeas \
  python2-pip \
  salt-master \
  salt-minion \
  screen \
  vim

RUN echo "auto_accept: True" >/etc/salt/master.d/autoaccept.conf \
	&& echo "master: localhost" >/etc/salt/minion.d/master.conf \
	&& echo "test-minion.system" >/etc/salt/minion_id
##################Golang install###################################

RUN yum update -y                                                  \
 && yum install -y wget                                            \
                  tar                                              \
                  g++                                              \
                  gcc                                              \
                  libc6-dev                                        \
                  make                                             \
                  git                                              \
                  python-setuptools

RUN set -eux; \
	url="https://golang.org/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz"; \
	wget -O go.tgz "$url"; \
	echo "${goRelSha256} *go.tgz" | sha256sum -c -; \
	tar -C /usr/local -xzf go.tgz; \
	rm go.tgz; \
	\
	if [ "linux-amd64" = 'src' ]; then \
		echo >&2; \
		echo >&2 'error: UNIMPLEMENTED'; \
		echo >&2 'TODO install golang-any from jessie-backports for GOROOT_BOOTSTRAP (and uninstall after build)'; \
		echo >&2; \
		exit 1; \
	fi; \
	\
	export PATH="/usr/local/go/bin:$PATH"; \
	go version

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

RUN   echo '#!/usr/bin/env bash' > /bin/entrypoint.sh \
	&& echo 'set -o errexit' >> /bin/entrypoint.sh \
	&& echo 'set -o pipefail' >> /bin/entrypoint.sh \
	&& echo 'set -o nounset' >> /bin/entrypoint.sh \
	&& echo 'salt-master -d' >> /bin/entrypoint.sh \
	&& echo 'salt-minion -d' >> /bin/entrypoint.sh \
	&& echo 'sleep 5' >> /bin/entrypoint.sh \
	&& echo 'exec /bin/bash' >> /bin/entrypoint.sh \
	&& cat /bin/entrypoint.sh \
	&& chmod 777 /bin/entrypoint.sh \
	&& wget https://raw.githubusercontent.com/docker-library/golang/2a15dfff04accfd31c2a45b3bb7423aa86aa2d60/1.8/jessie/go-wrapper -P /usr/local/bin/

WORKDIR $GOPATH/src/app
