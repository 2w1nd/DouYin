# DouYin
bytedance youth training camp project

## 1. 环境配置
- GO 1.16
- gin 1.17
- gorm 1.21
## 2 目录结构
```
├──DouYin
    ├── config          配置文件
    ├── controller      控制器目录
    ├── core            启动器目录，用于加载server，读取配置
    ├── global          存放全局对象
    ├── initialize      存放初始化配置
    ├── middleware      存放中间件
    ├── model           数据库访问目录
    ├── public          资源文件
    ├── router          路由目录
    └── service         业务代码目录
```