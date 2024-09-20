# parser


## goyacc
基于 tidb 的 goyacc 版本：https://github.com/pingcap/parser/releases/tag/v2.1.5

### 使用

```bash
# 编译 goyacc
make goyacc

# 生成 parser.go
make parser

```

- 修改 parser.y 后请添加对应的测试用例，并保证测试通过
```bash
make test
```
