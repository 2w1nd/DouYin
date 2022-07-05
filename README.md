# DouYin
bytedance youth training camp project

## 1. 环境依赖
- GO 1.16
- gin 1.17
- gorm 1.21
- viper 1.9.0
- redis 8.11.4
- jwt-go v3.2.0
- MySQL 5.6
## 2 目录结构
```
├──common               一方包
    ├──convert          转换器
    ├──context          上下文
    ├──entity           实体
        ├──dto
        ├──request
        ├──response
        ├──vo
    model               数据库模型
    
├──douyinService
    ├── cache           缓存目录
    ├── config          配置文件
    ├── controller      控制器目录
    ├── core            启动器目录，用于加载server，读取配置
    ├── global          存放全局对象
    ├── initialize      存放初始化配置
    ├── middleware      存放中间件
    ├── repository      存放数据库操纵语句
    ├── public          资源文件
    ├── resource        存放配置文件
    ├── router          路由目录
    └── service         业务代码目录
    └── utils           工具目录
```

## 3 项目成果

## 4 文档地址

[开发文档 - Feishu Doc](https://z5g9mu4ryx.feishu.cn/docs/doccnp5EE8GgHSfZPQQrp9HMjag)
