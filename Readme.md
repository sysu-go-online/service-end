# go-online service end

## 简介

## 维护者

## 依赖的环境变量

- DATABASE_ADDRESS 数据库地址，默认为localhost
- DATABASE_PORT 数据库端口，默认为3306
- DEVELOP 是否为开发环境，默认为false
- DOCKER_ADDRESS 容器服务地址，默认为localhost
- DOCKER_PORT 容器服务端口，默认为8888

## 依赖的外部软件

- mysql

  需要包含有go-online数据库，该库包含有project表，具体说明参见技术文档

## 运行方式

`go run main`

## 测试版本相关说明

- 根目录

  每个用户根目录始终为/home/username/

  每个项目根目录始终为/home/username/src/github.com/projectname

- 用户

  默认全局用户名为golang

- 项目

  默认全局项目名为test
