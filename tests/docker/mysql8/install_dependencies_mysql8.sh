#!/bin/bash

set -euo pipefail

# 更新源
sed -i 's/mirrorlist/#mirrorlist/g' /etc/yum.repos.d/CentOS-*
sed -i 's|#baseurl=http://mirror.centos.org|baseurl=http://vault.centos.org|g' /etc/yum.repos.d/CentOS-*

# install mysql
yum install -y wget perl net-tools etcd curl libaio libaio-devel numactl make git gcc openssl  percona-icu-data-files
yum install -y glibc
yum update -y
# modify timezone
ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

mkdir -p /dist
cd /dist

wget https://downloads.percona.com/downloads/Percona-Server-8.0/Percona-Server-8.0.30-22/binary/redhat/7/x86_64/percona-server-server-8.0.30-22.1.el7.x86_64.rpm
wget https://downloads.percona.com/downloads/Percona-Server-8.0/Percona-Server-8.0.30-22/binary/redhat/7/x86_64/percona-server-shared-compat-8.0.30-22.1.el7.x86_64.rpm
wget https://downloads.percona.com/downloads/Percona-Server-8.0/Percona-Server-8.0.30-22/binary/redhat/7/x86_64/percona-server-shared-8.0.30-22.1.el7.x86_64.rpm
wget https://downloads.percona.com/downloads/Percona-Server-8.0/Percona-Server-8.0.30-22/binary/redhat/7/x86_64/percona-server-client-8.0.30-22.1.el7.x86_64.rpm
wget https://downloads.percona.com/downloads/Percona-Server-8.0/Percona-Server-8.0.30-22/binary/redhat/7/x86_64/percona-icu-data-files-8.0.30-22.1.el7.x86_64.rpm

rpm -ivh percona-server-server-8.0.30-22.1.el7.x86_64.rpm \
  percona-server-shared-compat-8.0.30-22.1.el7.x86_64.rpm \
  percona-server-shared-8.0.30-22.1.el7.x86_64.rpm \
  percona-server-client-8.0.30-22.1.el7.x86_64.rpm \
  percona-icu-data-files-8.0.30-22.1.el7.x86_64.rpm

rm -rf percona-server-server-8.0.30-22.1.el7.x86_64.rpm \
   percona-server-shared-compat-8.0.30-22.1.el7.x86_64.rpm \
   percona-server-shared-8.0.30-22.1.el7.x86_64.rpm \
   percona-server-client-8.0.30-22.1.el7.x86_64.rpm \
   percona-icu-data-files-8.0.30-22.1.el7.x86_64.rpm

groupadd -r --gid 2000 work
useradd -r -g work --uid 1000 work
mkdir -p /data/mysql /data/tmp /data/etc
chown -R work:work /data

# install golang
wget https://go.dev/dl/go1.16.15.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.15.linux-amd64.tar.gz
