## GoLang WebSocket Cluster
### 简介
 Golang Websocket 集群方案demo，通过Redis的Pub&Sub功能，实现多个Websocket Server同步消息，解决大并发需求。
### 依赖
- [Golang](https://golang.google.cn)
- [Redis](https://redis.io/)

### 参考文档
- [How to Scale WebSockets](https://hackernoon.com/scaling-websockets-9a31497af051)
### 参考代码
- [gorilla/websocket](https://github.com/gorilla/websocket/tree/master/examples/chat)
### 配置项
  ```yaml
    server:
      port: 8000 #运行端口
    redis:
      host: redis
      port: 6379
      password:
      db: 0
    message:
      channel: message_channel #消息channel名称
  ```
### Docker 单机运行体验
```shell
docker compose up -d
```
使用websocket客户端连接[ws://localhost:8000](ws://localhost:8000)

### 编译运行
```shell
mv config.example.yml config.yml
go build
./go-websocket-cluster
```
### 集群部署
*部署多套后(多服务器或不同端口)，使用SLB或Nginx反向代理实现负载均衡*

参见：[WebSocket 集群方案总结](https://pathbox.github.io/2018/03/06/socket-io-websocket-cluster-SLB-LVS/)

### 示例功能
  _type 定义请见 [entity/entity.go](https://github.com/yeesuu/go-websocket-cluster/blob/master/entity/entity.go)_
#### 获取在线人数: 
  - request
    ```json
    {
      "type": 3
    }
    ```
  - response
    ```json
    {
      "type": 1,
      "data": 1,
      "timestamp": 1618811942
    }
    ```
#### 获取点赞人数
  - request
    ```json
    {
      "type": 2
    }
    ```
  - response
    ```json
    {
      "type": 4,
      "data": 0,
      "timestamp": 1618811942
    }
    ```
#### 发送普通信息
  - request
    ```json
    {
      "type": 5,
      "data": "hello"
    }
    ```
  - response
    ```json
    {
      "type": 5,
      "data": "hello",
      "timestamp": 1618811942
    }
    ```
    
*更多功能请根据需求自行开发*