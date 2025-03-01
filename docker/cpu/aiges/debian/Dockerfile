FROM golang:1.18-bullseye as go-builder

LABEL maintainer="ybyang7@iflytek.com"

ENV LANG=C.UTF-8 LC_ALL=C.UTF-8

# bullseye

ENV TZ=Asia/Shanghai \
DEBIAN_FRONTEND=noninteractive

MAINTAINER ybyang7@iflytek.com
COPY docker/repos/bullseye/sources.list /etc/apt/sources.list

ENV GOPROXY=https://goproxy.cn,direct
RUN apt-get update && apt-get install -y libnuma-dev build-essential git sudo
COPY . /home/aiges
WORKDIR /home/aiges
RUN make build-pack


FROM golang:1.18-bullseye as final
LABEL maintainer="ybyang7@iflytek.com"

ENV LANG=C.UTF-8 LC_ALL=C.UTF-8
MAINTAINER ybyang7@iflytek.com
ENV TZ Asia/Shanghai

COPY docker/repos/bullseye/sources.list /etc/apt/sources.list

RUN DEBIAN_FRONTEND=noninteractive apt update &&apt install -y libnuma-dev git vim
WORKDIR /home/aiges

# setup native python

# runtime dependencies
RUN set -eux; \
	apt-get update; \
	apt-get install -y --no-install-recommends \
		libbluetooth-dev \
		tk-dev \
		uuid-dev \
    	libffi-dev \
    	libssl-dev \
        libbz2-dev \
        liblzma-dev \
    ; \
	rm -rf /var/lib/apt/lists/*

ENV GPG_KEY E3FF2839C048B25C084DEBE9B26995E310250568
ENV PYTHON_VERSION 3.9.13

RUN set -eux; \
	\
	wget -O python.tar.xz "https://registry.npmmirror.com/-/binary/python/${PYTHON_VERSION%%[a-z]*}/Python-$PYTHON_VERSION.tar.xz"; \
	wget -O python.tar.xz.asc "https://registry.npmmirror.com/-/binary/python/${PYTHON_VERSION%%[a-z]*}/Python-$PYTHON_VERSION.tar.xz.asc"; \
	GNUPGHOME="$(mktemp -d)"; export GNUPGHOME; \
	gpg --batch --keyserver hkps://keys.openpgp.org --recv-keys "$GPG_KEY"; \
	gpg --batch --verify python.tar.xz.asc python.tar.xz; \
	command -v gpgconf > /dev/null && gpgconf --kill all || :; \
	rm -rf "$GNUPGHOME" python.tar.xz.asc; \
	mkdir -p /usr/src/python; \
	tar --extract --directory /usr/src/python --strip-components=1 --file python.tar.xz; \
	rm python.tar.xz; \
	\
	cd /usr/src/python; \
	gnuArch="$(dpkg-architecture --query DEB_BUILD_GNU_TYPE)"; \
	./configure \
		--build="$gnuArch" \
		--enable-loadable-sqlite-extensions \
		--enable-optimizations \
		--enable-option-checking=fatal \
		--enable-shared \
		--with-system-expat \
		--without-ensurepip \
	; \
	nproc="$(nproc)"; \
	make -j "$nproc" \
	; \
	make install; \
	\
# enable GDB to load debugging data: https://github.com/docker-library/python/pull/701
	bin="$(readlink -ve /usr/local/bin/python3)"; \
	dir="$(dirname "$bin")"; \
	mkdir -p "/usr/share/gdb/auto-load/$dir"; \
	cp -vL Tools/gdb/libpython.py "/usr/share/gdb/auto-load/$bin-gdb.py"; \
	\
	cd /; \
	rm -rf /usr/src/python; \
	\
	find /usr/local -depth \
		\( \
			\( -type d -a \( -name test -o -name tests -o -name idle_test \) \) \
			-o \( -type f -a \( -name '*.pyc' -o -name '*.pyo' -o -name 'libpython*.a' \) \) \
		\) -exec rm -rf '{}' + \
	; \
	\
	ldconfig; \
	\
	python3 --version

# make some useful symlinks that are expected to exist ("/usr/local/bin/python" and friends)
RUN set -eux; \
	for src in idle3 pydoc3 python3 python3-config; do \
		dst="$(echo "$src" | tr -d 3)"; \
		[ -s "/usr/local/bin/$src" ]; \
		[ ! -e "/usr/local/bin/$dst" ]; \
		ln -svT "$src" "/usr/local/bin/$dst"; \
	done

# if this is called "PIP_VERSION", pip explodes with "ValueError: invalid truth value '<VERSION>'"
ENV PYTHON_PIP_VERSION 22.0.4
# https://github.com/docker-library/python/issues/365
ENV PYTHON_SETUPTOOLS_VERSION 58.1.0
# https://github.com/pypa/get-pip
ENV PYTHON_GET_PIP_URL https://bootstrap.pypa.io/get-pip.py
ENV PYTHON_GET_PIP_SHA256 ba3ab8267d91fd41c58dbce08f76db99f747f716d85ce1865813842bb035524d

# will cause error here
COPY ./docker/get-pip.py /

RUN set -eux; \
	\
#	wget -O get-pip.py "$PYTHON_GET_PIP_URL"; \
#	echo "$PYTHON_GET_PIP_SHA256 *get-pip.py" | sha256sum -c -; \
	\
	export PYTHONDONTWRITEBYTECODE=1; \
	\
	python /get-pip.py \
		--disable-pip-version-check \
		--no-cache-dir \
		--no-compile \
		"pip==$PYTHON_PIP_VERSION" \
		"setuptools==$PYTHON_SETUPTOOLS_VERSION" \
	; \
	rm -f get-pip.py; \
	\
	pip --version


ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/home/aiges:/home/aiges/library


# need xtest
# need aiges.toml
# xtest.toml
COPY --from=go-builder /home/aiges/dist/aiges_linux_amd64 .
COPY --from=go-builder /home/aiges/ai_cpython_wrapper  /home/aiges/ai_cpython_wrapper

RUN cd /home/aiges/ai_cpython_wrapper && make && cp -r /home/aiges/ai_cpython_wrapper/wrapper_lib/* /home/aiges/library && rm -rf /home/aiges/ai_cpython_wrapper

CMD ["sh", "-c", "./AIservice -m=0 -c=aiges.toml -s=svcName -u=http://companion.xfyun.iflytek:6868 -p=AIaaS -g=dx"]