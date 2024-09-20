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
    -o -path './tests/e2e/cmd/*' \
    \)))

NUM=0
declare FAILED_FILE

for f in ${FILELIST[@]}; do
    c=$(tail -c 1 "$f" | wc -l)
    if [ "$c" -eq 0 ]; then
        FAILED_FILE+=($f)
        NUM=$((NUM + 1))
    fi
done

if [ $NUM -ne 0 ]; then
    echo "error: following files do not end with newline, please run hack/update-EOF.sh to fix them"
    printf '%s\n' "${FAILED_FILE[@]}"
    exit 1
else
    echo "all files pass checking."
fi
