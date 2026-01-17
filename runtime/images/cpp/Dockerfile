FROM alpine:3.19

# Install minimal build toolchain + ccache
# build-base includes: gcc, g++, musl-dev, make, binutils, busybox
RUN apk add --no-cache \
    build-base \
    ccache \
    coreutils \
    bash \
    && rm -rf /var/cache/apk/*

# --- CCache configuration ---
# Put cache in tmpfs-mounted path (fast, auto-cleared)
ENV CCACHE_DIR=/tmp/ccache
ENV CC="ccache gcc"
ENV CXX="ccache g++"
ENV CCACHE_COMPRESS=1
ENV CCACHE_MAXSIZE=50M

# Create sandbox (rootfs will be read-only at runtime)
RUN mkdir -p /sandbox

WORKDIR /sandbox

# Default shell (used by docker exec)
CMD ["sh"]
