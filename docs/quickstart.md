# 快速入门

## 编译安装

gaea基于go开发，基于glide进行版本管理，并依赖goyacc、gofmt等工具。

* 首先安装依赖包
glide install

* 编译二进制包
make

## 执行

编译之后在bin目录会有gaea、gaea-cc两个可执行文件。etc目录下为配置文件，如果想快速体验gaea功能，可以采用file配置方式，然后在etc/file/namespace下添加对应租户的json文件，该目录下目前有两个示例，可以直接修改使用。
./bin/gaea --help显示如下，其中-config是指定配置文件位置，默认为./etc/gaea.ini，具体配置见[配置说明](configuration.md)。

```bash
Usage of ./bin/gaea:
  -config string
    gaea config file (default "etc/gaea.ini")
```
