#!/bin/bash

# 一般原本 Docker 的数据库容器无法在 containerd 上执行，进行以下修正
mkdir /var/run/mysqld
# RUN useradd -m mysql
chown mysql:mysql /var/run/mysqld
chmod 777 /var/run/mysqld

# user.sql 为一开始执行 mysqld 服务时，所需要执行的 SQL 脚本，会创建一个用户 xiaomi，并且设置密码
mysqld --init-file=/home/mariadb/user.sql
