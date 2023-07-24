#!/bin/bash

set -euo pipefail

# prepare env
go install github.com/onsi/ginkgo/v2/ginkgo@v2.3.1

# prepare etcd env
# shellcheck disable=SC2069
etcd 2>&1 1>>etcd.log &
echo "pwd"
pwd
echo "ls"
ls

# start mysql master
cp ./tests/docker/my3319_example.cnf /data/etc/my3319.cnf
mysqld --defaults-file=/data/etc/my3319.cnf --user=work --initialize-insecure
mysqld --defaults-file=/data/etc/my3319.cnf --user=work &
sleep 3

# start mysql slave
cp ./tests/docker/my3329_example.cnf /data/etc/my3329.cnf
mysqld --defaults-file=/data/etc/my3329.cnf --user=work --initialize-insecure
mysqld --defaults-file=/data/etc/my3329.cnf --user=work &
sleep 3

# 授权
mysql -hlocalhost -P3319 -uroot -S/data/tmp/mysql3319.sock <<EOF
reset master;
GRANT REPLICATION SLAVE, REPLICATION CLIENT on *.* to 'mysqlsync'@'%' IDENTIFIED BY 'mysqlsync';
-- User for superroot system
GRANT ALL ON *.* TO 'superroot'@'%' IDENTIFIED BY 'superroot' WITH GRANT OPTION;
EOF

# 构建主从关系
mysql -hlocalhost -P3329 -uroot -S/data/tmp/mysql3329.sock <<EOF
GRANT REPLICATION SLAVE, REPLICATION CLIENT on *.* to 'mysqlsync'@'%' IDENTIFIED BY 'mysqlsync';
-- User for superroot system
GRANT ALL ON *.* TO 'superroot'@'%' IDENTIFIED BY 'superroot' WITH GRANT OPTION;
CHANGE MASTER TO MASTER_HOST='127.0.0.1', MASTER_PORT=3319, MASTER_USER='mysqlsync', MASTER_PASSWORD='mysqlsync', MASTER_AUTO_POSITION=1;
START SLAVE;
DO SLEEP(1);
EOF


# start gaea
ps aux | grep mysql
# start gaea_cc
echo "run here"

# TODO: run test,but err in this file
#ginkgo --v --progress --trace --flake-attempts=1 ./tests/e2e/
