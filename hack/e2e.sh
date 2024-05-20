#!/bin/bash

set -euo pipefail

# prepare env
if [ $(ls $HOME/go/bin/|grep ginkgo|wc -l) -eq 0 ];then
    go install github.com/onsi/ginkgo/v2/ginkgo@v2.3.1
fi

# prepare etcd env
# shellcheck disable=SC2069

function check_pid() {
    value=$1
    pid_num=$(ps aux|grep "$value"|grep -v 'grep'|wc -l)
    echo $pid_num
}

function check_mysql_dir() {
    value=$1
    pid_num=$(ps aux|grep etcd|grep -v 'grep'|wc -l)
    dir_num=$(ls /data/mysql/|grep "$value"|wc -l)
    echo $dir_num
}

if [ $(check_pid "etcd") -eq 0 ];then
    etcd --data-dir bin/etcd 2>&1 1>>bin/etcd.log &
fi

# Start 2 Mysql Cluster
# Cluster-1: 3319(master), 3329(slave), 3339(slave)
# Prepare Cluster-1 for master
if [ $(check_mysql_dir "3319") -eq 0 ];then
    cp ./tests/docker/my3319.cnf /data/etc/my3319.cnf
    mysqld --defaults-file=/data/etc/my3319.cnf --user=work --initialize-insecure
    mysqld --defaults-file=/data/etc/my3319.cnf --user=work &
    sleep 3
    mysql -h127.0.0.1 -P3319 -uroot -S/data/tmp/mysql3319.sock <<EOF
    reset master;
    GRANT REPLICATION SLAVE, REPLICATION CLIENT on *.* to 'mysqlsync'@'%' IDENTIFIED BY 'mysqlsync';
    GRANT ALL ON *.* TO 'superroot'@'%' IDENTIFIED BY 'superroot' WITH GRANT OPTION;
    GRANT SELECT, INSERT, UPDATE, DELETE ON *.* TO 'gaea_backend_user'@'%' IDENTIFIED BY 'gaea_backend_pass';
    GRANT  REPLICATION SLAVE, REPLICATION CLIENT ON *.* TO 'gaea_backend_user'@'%' IDENTIFIED BY 'gaea_backend_pass';
EOF
else
    if [ $(check_pid "my3319") -eq 0 ];then
        mysqld --defaults-file=/data/etc/my3319.cnf --user=work &
    fi
fi

# Prepare Cluster-1 slave 1
if [ $(check_mysql_dir "3329") -eq 0 ];then
    cp ./tests/docker/my3329.cnf   /data/etc/my3329.cnf
    mysqld --defaults-file=/data/etc/my3329.cnf --user=work --initialize-insecure
    mysqld --defaults-file=/data/etc/my3329.cnf --user=work &
    sleep 3
    mysql -h127.0.0.1 -P3329 -uroot -S/data/tmp/mysql3329.sock <<EOF
    CHANGE MASTER TO MASTER_HOST='127.0.0.1', MASTER_PORT=3319, MASTER_USER='mysqlsync', MASTER_PASSWORD='mysqlsync', MASTER_AUTO_POSITION=1;
    START SLAVE;
    DO SLEEP(1);
EOF
else
    if [ $(check_pid "my3329") -eq 0 ];then
        mysqld --defaults-file=/data/etc/my3329.cnf --user=work &
    fi
fi

# Prepare Cluster-1 slave 2
if [ $(check_mysql_dir "3339") -eq 0 ];then
    cp ./tests/docker/my3339.cnf   /data/etc/my3339.cnf
    mysqld --defaults-file=/data/etc/my3339.cnf --user=work --initialize-insecure
    mysqld --defaults-file=/data/etc/my3339.cnf --user=work &
    sleep 3
    mysql -h127.0.0.1 -P3339 -uroot -S/data/tmp/mysql3339.sock <<EOF
    CHANGE MASTER TO MASTER_HOST='127.0.0.1', MASTER_PORT=3319, MASTER_USER='mysqlsync', MASTER_PASSWORD='mysqlsync', MASTER_AUTO_POSITION=1;
    START SLAVE;
    DO SLEEP(1);
EOF
else
    if [ $(check_pid "my3339") -eq 0 ];then
        mysqld --defaults-file=/data/etc/my3339.cnf --user=work &
    fi
fi

# Cluster-2: 3379(master)
if [ $(check_mysql_dir "3379") -eq 0 ];then
    cp ./tests/docker/my3379.cnf /data/etc/my3379.cnf
    mysqld --defaults-file=/data/etc/my3379.cnf --user=work --initialize-insecure
    mysqld --defaults-file=/data/etc/my3379.cnf --user=work &
    sleep 3
    mysql -h127.0.0.1 -P3379 -uroot -S/data/tmp/mysql3379.sock <<EOF
    reset master;
    GRANT REPLICATION SLAVE, REPLICATION CLIENT on *.* to 'mysqlsync'@'%' IDENTIFIED BY 'mysqlsync';
    GRANT ALL ON *.* TO 'superroot'@'%' IDENTIFIED BY 'superroot' WITH GRANT OPTION;
    GRANT SELECT, INSERT, UPDATE, DELETE ON *.* TO 'gaea_backend_user'@'%' IDENTIFIED BY 'gaea_backend_pass';
    GRANT REPLICATION SLAVE, REPLICATION CLIENT ON *.* TO 'gaea_backend_user'@'%' IDENTIFIED BY 'gaea_backend_pass';
EOF
else
    if [ $(check_pid "my3379") -eq 0 ];then
        mysqld --defaults-file=/data/etc/my3379.cnf --user=work &
    fi
fi

# Cluster-3: 3349(master)
if [ $(check_mysql_dir "3349") -eq 0 ];then
cp ./tests/docker/my3349.cnf /data/etc/my3349.cnf
    mysqld --defaults-file=/data/etc/my3349.cnf --user=work --initialize-insecure
    mysqld --defaults-file=/data/etc/my3349.cnf --user=work &
    sleep 3
    mysql -h127.0.0.1 -P3349 -uroot -S/data/tmp/mysql3349.sock <<EOF
    reset master;
    GRANT REPLICATION SLAVE, REPLICATION CLIENT on *.* to 'mysqlsync'@'%' IDENTIFIED BY 'mysqlsync';
    GRANT ALL ON *.* TO 'superroot'@'%' IDENTIFIED BY 'superroot' WITH GRANT OPTION;
    GRANT SELECT, INSERT, UPDATE, DELETE ON *.* TO 'gaea_backend_user'@'%' IDENTIFIED BY 'gaea_backend_pass';
    GRANT REPLICATION SLAVE, REPLICATION CLIENT ON *.* TO 'gaea_backend_user'@'%' IDENTIFIED BY 'gaea_backend_pass';
EOF
else
    if [ $(check_pid "my3349") -eq 0 ];then
        mysqld --defaults-file=/data/etc/my3349.cnf --user=work &
    fi
fi
