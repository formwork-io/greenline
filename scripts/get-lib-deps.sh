#!/usr/bin/env bash
# Builds the latest version of libsodium, zeromq, and czmq
set -e  # exit on error
dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"/../
source "$dir"/env.sh
assert-env-or-die "BUILD"
assert-env-or-die "GL_SODIUM_VER"
assert-env-or-die "GL_ZEROMQ_VER"
assert-env-or-die "GL_CZMQ_VER"

gl_lib_deps="$BUILD/deps/libs"
mkdir -p "$gl_lib_deps"
cd "$gl_lib_deps"

voverride PKG_CONFIG_PATH "$gl_lib_deps"/lib/pkgconfig

libsodium_tgz="libsodium-${GL_SODIUM_VER}.tar.gz"
github_url="https://github.com/jedisct1/libsodium/releases"
libsodium_url="$github_url/download/$GL_SODIUM_VER/$libsodium_tgz"

zeromq_tgz="zeromq-${GL_ZEROMQ_VER}.tar.gz"
github_url="https://github.com/zeromq/zeromq4-1/releases"
zeromq_url="$github_url/download/v$GL_ZEROMQ_VER/${zeromq_tgz}"

czmq_tgz="czmq-${GL_CZMQ_VER}.tar.gz"
github_url="https://github.com/zeromq/czmq/releases"
czmq_url="$github_url/download/v$GL_CZMQ_VER/${czmq_tgz}"

# where we keep archives
download_root="$gl_lib_deps"/download
mkdir -p "$download_root"
# where we build from
work_root="$gl_lib_deps"/work
mkdir -p "$work_root"

# clear old work root
rm -fr "$work_root"

# download everything
cd "$download_root"

function download() {
    if [ ! -r "$3" ]; then
        vrun-or-die wget "$2" -O "$3"
    fi
}

download "libsodium" "$libsodium_url" "$libsodium_tgz"
download "zeromq" "$zeromq_url" "$zeromq_tgz"
download "czmq" "$czmq_url" "$czmq_tgz"

function extract() {
    echo -en "Extracting $1 $2... "
    tarfile="$download_root/$3"
    tar --transform="s|$1-$2|$2|" -xzf "$tarfile"
    echo "OK"
}

function cmmi() {
    echo -en "Configuring $1 $2... "
    cd "$3/$2"
    ./configure --silent --prefix="$gl_lib_deps"
    echo "OK"
    echo -e "Starting $1 build.\n--"
    make --silent all install
    echo -e "--\n$1 build complete."
}

libsodium_work="$work_root-libsodium"
mkdir -p "$libsodium_work"
cd "$libsodium_work"
extract "libsodium" "$GL_SODIUM_VER" "$libsodium_tgz"
cmmi "libsodium" "$GL_SODIUM_VER" "$libsodium_work"

zeromq_work="$work_root-zeromq"
mkdir -p "$zeromq_work"
cd "$zeromq_work"
extract "zeromq" "$GL_ZEROMQ_VER" "$zeromq_tgz"
cmmi "zeromq" "$GL_ZEROMQ_VER" "$zeromq_work"

czmq_work="$work_root-czmq"
mkdir -p "$czmq_work"
cd "$czmq_work"
extract "czmq" "$GL_CZMQ_VER" "$czmq_tgz"
cmmi "czmq" "$GL_CZMQ_VER" "$czmq_work"

touch "$gl_lib_deps"/.done

