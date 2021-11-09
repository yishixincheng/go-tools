# go-tools

# 下载
go get -u github.com/yishixincheng/go-tools

# 使用说明
- 配置application.toml 如数据库配置

#根据数据库表自动生成model文件 示例如下

- 单表生成
go run main.go sql2code --d=数据库名 -t=表名

- 多表生成
go run main.go sql2code --d=数据库名 -t=”表名1,表名2“

- 数据库下面全部表生成

go run main.go sql2code --d=数据库名 -t=*