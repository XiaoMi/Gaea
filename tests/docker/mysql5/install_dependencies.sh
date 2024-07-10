#!/bin/bash

set -euo pipefail

# 更新源
sed -i 's/mirrorlist/#mirrorlist/g' /etc/yum.repos.d/CentOS-*
sed -i 's|#baseurl=http://mirror.centos.org|baseurl=http://vault.centos.org|g' /etc/yum.repos.d/CentOS-*

# install mysql
yum install -y wget perl net-tools etcd curl libaio libaio-devel numactl make git gcc
yum install -y glibc
yum update -y
# modify timezone
ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

cd /dist
wget https://downloads.percona.com/downloads/Percona-Server-5.7/Percona-Server-5.7.25-28/binary/redhat/7/x86_64/Percona-Server-server-57-5.7.25-28.1.el7.x86_64.rpm
wget https://downloads.percona.com/downloads/Percona-Server-5.7/Percona-Server-5.7.25-28/binary/redhat/7/x86_64/Percona-Server-shared-compat-57-5.7.25-28.1.el7.x86_64.rpm
wget https://downloads.percona.com/downloads/Percona-Server-5.7/Percona-Server-5.7.25-28/binary/redhat/7/x86_64/Percona-Server-shared-57-5.7.25-28.1.el7.x86_64.rpm
wget https://downloads.percona.com/downloads/Percona-Server-5.7/Percona-Server-5.7.25-28/binary/redhat/7/x86_64/Percona-Server-client-57-5.7.25-28.1.el7.x86_64.rpm

rpm -ivh Percona-Server-server-57-5.7.25-28.1.el7.x86_64.rpm \
  Percona-Server-shared-compat-57-5.7.25-28.1.el7.x86_64.rpm \
  Percona-Server-shared-57-5.7.25-28.1.el7.x86_64.rpm \
  Percona-Server-client-57-5.7.25-28.1.el7.x86_64.rpm

rm -rf Percona-Server-server-57-5.7.25-28.1.el7.x86_64.rpm \
   Percona-Server-shared-compat-57-5.7.25-28.1.el7.x86_64.rpm \
   Percona-Server-shared-57-5.7.25-28.1.el7.x86_64.rpm \
   Percona-Server-client-57-5.7.25-28.1.el7.x86_64.rpm

groupadd -r --gid 2000 work
useradd -r -g work --uid 1000 work
mkdir -p /data/mysql /data/tmp /data/etc
chown -R work:work /data

# install golang
wget https://go.dev/dl/go1.16.15.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.15.linux-amd64.tar.gz

# 提前安装一些依赖库，跑 CI 的时间就可以不用拉了，减少测试用例时间
/usr/local/go/bin/go install github.com/onsi/ginkgo/v2/ginkgo@v2.3.1
/usr/local/go/bin/go mod download github.com/shirou/gopsutil@v2.20.9+incompatible
/usr/local/go/bin/go mod download go.uber.org/atomic@v1.4.0
/usr/local/go/bin/go mod download github.com/golang/mock@v1.4.4
/usr/local/go/bin/go mod download github.com/gin-contrib/gzip@v0.0.1
/usr/local/go/bin/go mod download github.com/gin-gonic/gin@v1.7.2
/usr/local/go/bin/go mod download github.com/pingcap/errors@v0.11.1
/usr/local/go/bin/go mod download github.com/shopspring/decimal@v1.3.1
/usr/local/go/bin/go mod download github.com/hashicorp/go-version@v1.6.0
/usr/local/go/bin/go mod download github.com/go-ini/ini@v1.42.0
/usr/local/go/bin/go mod download github.com/pingcap/tipb@v0.0.0-20190226124958-833c2ffd2fe7
/usr/local/go/bin/go mod download github.com/coreos/etcd@v3.3.13+incompatible
/usr/local/go/bin/go mod download github.com/emirpasic/gods@v1.12.0
/usr/local/go/bin/go mod download github.com/prometheus/client_golang@v0.9.2
/usr/local/go/bin/go mod download github.com/cznic/mathutil@v0.0.0-20181122101859-297441e03548
/usr/local/go/bin/go mod download github.com/prometheus/common@v0.0.0-20181126121408-4724e9255275
/usr/local/go/bin/go mod download github.com/prometheus/client_model@v0.0.0-20180712105110-5c3871d89910
/usr/local/go/bin/go mod download github.com/beorn7/perks@v0.0.0-20180321164747-3a771d992973
/usr/local/go/bin/go mod download github.com/golang/protobuf@v1.5.2
/usr/local/go/bin/go mod download github.com/prometheus/procfs@v0.0.0-20181204211112-1dc9a6cbc91a
/usr/local/go/bin/go mod download github.com/remyoudompheng/bigfft@v0.0.0-20190321074620-2f0d2b0e0001
/usr/local/go/bin/go mod download github.com/mattn/go-isatty@v0.0.13
/usr/local/go/bin/go mod download github.com/gin-contrib/sse@v0.1.0
/usr/local/go/bin/go mod download gopkg.in/yaml.v2@v2.4.0
/usr/local/go/bin/go mod download github.com/ugorji/go@v1.1.7
/usr/local/go/bin/go mod download github.com/go-playground/validator/v10@v10.8.0
/usr/local/go/bin/go mod download golang.org/x/sys@v0.0.0-20220722155257-8c9f86f7a55f
/usr/local/go/bin/go mod download github.com/ugorji/go/codec@v1.1.7
/usr/local/go/bin/go mod download google.golang.org/protobuf@v1.28.0
/usr/local/go/bin/go mod download github.com/matttproud/golang_protobuf_extensions@v1.0.1
/usr/local/go/bin/go mod download golang.org/x/text@v0.3.7
/usr/local/go/bin/go mod download golang.org/x/crypto@v0.0.0-20210921155107-089bfa567519
/usr/local/go/bin/go mod download github.com/leodido/go-urn@v1.2.1
/usr/local/go/bin/go mod download github.com/go-playground/universal-translator@v0.17.0
/usr/local/go/bin/go mod download github.com/go-playground/locales@v0.13.0
/usr/local/go/bin/go mod download github.com/modern-go/reflect2@v1.0.1
/usr/local/go/bin/go mod download github.com/json-iterator/go@v1.1.11
/usr/local/go/bin/go mod download github.com/coreos/go-semver@v0.2.0
/usr/local/go/bin/go mod download google.golang.org/grpc@v1.21.0
/usr/local/go/bin/go mod download github.com/gogo/protobuf@v1.2.1
/usr/local/go/bin/go mod download golang.org/x/net@v0.0.0-20220722155237-a158d28d115b
/usr/local/go/bin/go mod download github.com/modern-go/concurrent@v0.0.0-20180306012644-bacd9c7ef1dd
/usr/local/go/bin/go mod download google.golang.org/genproto@v0.0.0-20180817151627-c66870c02cf8
