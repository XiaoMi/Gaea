#!/bin/bash
ginkgo --v --progress --trace --flake-attempts=1 --skip '^.*only mysql5:.*$' --skip '^.*shard join support test in.*$' --skip 'test dml set variables' --skip 'simple sql test' --skip='Unshard DML Support Test' ./tests/e2e/
