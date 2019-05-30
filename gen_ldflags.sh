#!/bin/bash

VERBOSE=${VERBOSE:-"0"}
V=""
if [[ "${VERBOSE}" == "1" ]];then
    V="-x"
    set -x
fi

ROOT="$(pwd)"

OUT=${1:?"output path"}
VERSION_PACKAGE=${2:?"version go package"} # istio.io/istio/pkg/version
BUILDPATH=${3:?"path to build"}

set -e

GOOS=${GOOS:-linux}
GOARCH=${GOARCH:-amd64}
GOBINARY=${GOBINARY:-go}
GOPKG="$GOPATH/pkg"
BUILDINFO=${BUILDINFO:-""}
STATIC=${STATIC:-1}
LDFLAGS="-extldflags -static"
GOBUILDFLAGS=${GOBUILDFLAGS:-""}
GCFLAGS=${GCFLAGS:-}
export CGO_ENABLED=0

if [[ "${STATIC}" !=  "1" ]];then
    LDFLAGS=""
fi

# gather buildinfo if not already provided
# For a release build BUILDINFO should be produced
# at the beginning of the build and used throughout
if [[ -z ${BUILDINFO} ]];then
    BUILDINFO=$(mktemp)
    ${ROOT}/gen_version.sh > ${BUILDINFO}
fi

# BUILD LD_VERSIONFLAGS
LD_VERSIONFLAGS=""
while read line; do
    read SYMBOL VALUE < <(echo $line)
    LD_VERSIONFLAGS=${LD_VERSIONFLAGS}" -X ${VERSION_PACKAGE}.${SYMBOL}=${VALUE}"
done < "${BUILDINFO}"

echo -pkgdir=${GOPKG}/${GOOS}_${GOARCH} -ldflags "\"${LDFLAGS} ${LD_VERSIONFLAGS}\"" ${BUILDPATH}
