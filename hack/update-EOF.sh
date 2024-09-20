#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

ROOTDIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)"

cd $ROOTDIR

FILELIST=($(find . -type f -not \( -path './bin/*' \
    -o -path './etc/*' \
    -o -path './.git/*' \
    -o -path '*.png' \
    -o -path './.idea/*' \
    -o -path './.DS_Store' \
    -o -path './*/.DS_Store' \
    -o -path './docs/*' \
    -o -path './logs/*' \
    -o -path './parser/goyacc' \
    \)))

for f in ${FILELIST[@]}; do
    c=$(tail -c 1 "$f" | wc -l)
    if [ "$c" -eq 0 ]; then
        echo "find file $f do not end with newline, fixing it"
        printf "\n" >> $f
    fi
done
