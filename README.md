# RabbitMQ
This is an RabbitMQ indicator mobile testing project


## 启动RabbitMQ
```shell
cd deployment/docker-compose
docker-compose up -d
```

## 启动消息消费

```shell
go run cmd/entry.go
```

## 启动消息生产

```shell
cd tests
go test -run TestPublishMessage
```