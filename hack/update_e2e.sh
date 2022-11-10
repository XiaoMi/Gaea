#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

MASTER_USER=$1
MASTER_PASSWORD=$2
ROOTDIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)"

cd $ROOTDIR
echo $ROOTDIR

cp ./bin/gaea ./tests/e2e/cmd/
sed -i "s#^master_user.*#master_user: $MASTER_USER#" ./tests/e2e/config/config.yaml
sed -i "s#^master_password.*#master_user: $MASTER_PASSWORD#" ./tests/e2e/config/config.yaml
sed -i "s#^slave_user.*#slave_user: $MASTER_USER#" ./tests/e2e/config/config.yaml
sed -i "s#^slave_password.*#slave_user: $MASTER_PASSWORD#" ./tests/e2e/config/config.yaml
