#!/bin/bash
ginkgo --v --progress --trace --flake-attempts=1 --skip "^.*only mysql8:.*$" ./tests/e2e/
