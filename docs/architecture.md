# 架构设计

## 模块划分

gaea包含四个模块，分别是gaea-proxy、gaea-cc、gaea-agent、gaea-web。gaea-proxy为在线代理，负责承接sql流量，gaea-cc是中控模块，负责gaea-proxy的配置管理及一些后台任务，gaea-agent部署在mysql所在的机器上，负责实例创建、管理、回收等工作，gaea-web是gaea的一个管理界面，使gaea整体使用更加方便。

## 架构图

![gaea架构图](assets/architecture.png)
