# grpc 微服务

设计一个简单的数据存储的服务，通过get获取数据，通过put设置数据，类似于Redis的key,value存储。

服务有两个API:

- `get`: 根据key获取内容
- `put`: 根据key设置内容

那么初步定义一个Service实现这两个功能:

```golang
type Service interface {
	Get(ctx context.Context, key string) (val string, err error)
	Put(ctx context.Context, key, val string) (err error)
}
```

以下内容全部通过go-kit的一些组件来实现。

![](http://source.qiniu.cnd.nsini.com/images/2020/03/df/94/b1/20200329-42107a33e5251f65bbf9b27e941dffa0.jpeg?imageView2/2/w/1280/interlace/0/q/70)

## 测试

```bash
$ make run
GO111MODULE=on /usr/local/go/bin/go run ./cmd/main.go -http-addr :8080 -grpc-addr :8081
level=error ts=2020-03-28T10:45:05.923565Z caller=service.go:106 transport=HTTP addr=:8080
```

执行客户端测试命令:

```bash
$ go run ./client/grpc/client.go
$ go run ./client/http/client.go
level=info ts=2020-03-28T10:45:44.793353Z caller=logging.go:41 method=Put key=hello val=world took=2.142µs err=null
level=info ts=2020-03-28T10:45:44.794983Z caller=logging.go:28 method=Get key=hello val=world took=1.248µs err=null
level=info ts=2020-03-28T10:47:02.666247Z caller=logging.go:28 method=Get key=hello val=world took=1.396µs err=null
```

## 尾巴

本章所用的测试代码已经更新到了Github上，如果您觉得有参考价值的，可以将代码clone 下来，最好能给个star。

Github: `https://github.com/icowan/grpc-world`

**谢谢了**

**如果我写的内容对您有用，谢谢大家了**

![](http://source.qiniu.cnd.nsini.com//static/pay/wechat-pay.JPG?imageView2/2/w/500/q/75|imageslim)

