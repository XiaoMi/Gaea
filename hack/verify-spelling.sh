#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

ROOT=$(unset CDPATH && cd $(dirname "${BASH_SOURCE[0]}")/.. && pwd)
cd $ROOT
OUTPUT_BIN=$ROOT/bin/output/
mkdir -p $OUTPUT_BIN

# ensure_misspell
echo "Installing misspell..."
GOBIN=$OUTPUT_BIN go install github.com/client9/misspell/cmd/misspell@v0.3.4


ignore_words=(
    "importas"
    "etc"
)

ret=0
git ls-files | xargs ${OUTPUT_BIN}/misspell -i ${ignore_words[@]} -error -o stderr || ret=$?
if [ $ret -eq 0 ]; then
    echo "Spellings all good!"
else
    echo "Found some typos, please fix them!"
    exit 1
fi
