# sqladuit

这是sqlaudit的主要组件，主要用来监听本机或[sqlpacket组件](http://gitlab.jixindatech.com/sql/sqlpacket) 发送过来的以太网帧，然后重新还原成原来的sql流量，依据规则进行匹配。同时提供了web接口进行规则管理告警配置等,prometheus提供了客户端数量和packet的监控。欢迎issue和star!

## Docker
```
docker build -t sqlaudit .
docker run -d \
     --name sqlaudit \
     -p 9696:9696 \
     -p 9797:9797 \
     -p 9898:9898 \
     -v /docker/sqlaudit/etc:/opt/sqlaudit/etc \
     sqlaudit
```

## Usage
### 修改配置文件
- web_addr 是web接口，web_user 和web_password 是登陆web管理后台的用户名和密码。
- 日志默认是输出到标准输出/错误的， 如果配置log_path所有日志记录到该目录下
- es_config 是配置sql事件的存储。
- alert_email 是配置邮件告警的信息，interval是指多长时间再次告警，防止邮件风暴。
- Elasticsearch 中的 settings配置在 etc/mapping.json中，比较简单。
- Prometheus的端口(9898)是在程序中写的，监控内容是发送数据的客户端和sql请求的数量，统计间隔是prometheus的scrap_interval。

### 启动方式
- 本地编译需要安装go 和npm 环境(vue)，然后make 生成可执行文件和web 文件，需要在编译目录执行可执行程序， 注意配置路径都是相对目录，否则可能找不到文件退出或异常出错。
- docker 形式, docker build -t sqlaudit . 即可生成sqlaudit镜像， 注意docker方式只能启动服务接受sqlpacket的帧数据，不能监听网卡流量。
- Sqlaudit -config etc/config.yml 执行配置文件，并开启接受sqlpacket发送的数据帧。
- Sqlaudit -config etc/config.yml  -c -i eth0(本地网卡接口)， 是开启监听本地网卡接口，接受sql流量。

### web管理
- 访问http://ip:9797 的web管理接口，配置文件中对应的 web_addr参数，使用配置文件中的web_user 和 web_password 登陆。
- Dashboard 显示相关时间段内统计的可视化内容，默认时间是一周，然后可以根据数据库进行查询。
- Sql配置菜单是配置相关的sql规则，其中操作类型的UNKNOWN是对应sql解析失败的情况，另外匹配条件是与的关系。
- 日志查询可以根据查询条件进行查询，注意Sql关键字查询是elasticsearch中的查询。
## Basic structure
![Image text](https://raw.githubusercontent.com/jixindatech/sqlaudit/master/doc/images/sqlaudit.jpg)
## Contributing
PRs accepted.

## License

Unlicense
