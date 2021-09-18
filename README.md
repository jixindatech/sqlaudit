# sqladuit

这是sqlaudit的主要组件，主要用来监听本机或[sqlpacket组件](http://gitlab.jixindatech.com/sql/sqlpacket) 发送过来的以太网帧，然后重新还原成原来的sql流量，依据规则进行匹配。同时提供了web接口进行规则管理告警配置等,prometheus提供了客户端数量和packet的监控。

## Compile

```
make
```

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
1. 启动监听本地端口的服务时需要指定 -c 和 -i inf 的参数，-c 是监听端口的bool值， -i 指定本地的端口接口
2. 启动接受以太网帧的服务时，需要用到[sqlpacket组件](http://gitlab.jixindatech.com/sql/sqlpacket)
3. 说明:主要要有三个服务，分别是解析包、web api、及prometheus 服务， docker安装时需要注意开发其端口

## Basic structure
![Image text](http://gitlab.jixindatech.com/sql/sqlaudit/-/raw/master/doc/images/sqlaudit.jpg?inline=true)
## Contributing

PRs accepted.

## License

Unlicense
