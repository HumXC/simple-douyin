# simple-douyin

字节青训营极简版抖音实现

## 外部依赖

-   ffmpeg: 用于视频压缩，从视频截取图片这两个功能。需要在部署的机器上正确安装 ffmpeg

## 编译

### C-GO

由于使用了 sqlite，所以需要启用 c-go

直接使用 `go build` 编译

## 配置文件

第一次运行会在可执行文件目录下生成 config.yaml 文件，一个配置文件如下所示

```yaml
douyin:
    # douyin 服务的监听地址
    serve-addr: 192.168.80.148:11451

    # 只支持 sqlite 和 mysql，如果使用 sqlite，dsn 就是数据库文件的名称
    sql:
        # type: sqlite
        # dsn: ./data.db
        type: mysql
        dsn: root:humxc@tcp(127.0.0.1:3306)/douyin?charset=utf8mb4&parseTime=True&loc=Local
    redis:
        addr: ""
        password: ""
        db: 0
    # /douyin/feed 中一次请求返回视频的最大数量，最大不会超过 30，如果 超过 30，会被强制设定为 30
    # 该值的详细定义见客户端文档
    feed-num: 30
    # 视频屠夫是进行视频压缩和获取封面的单元，见 videos.Butcher
    # Butcher 是多协程工作的，该值是同时进行视频处理的最大协程数量
    video-butcher-max-job: 16
# 这是存储服务，只用于存储文件
storage:
    # 存储数据的根文件夹
    data-dir: Data
    # 服务的监听地址，此地址必须为明确的 ip，不能是 ":11452"
    serve-addr: 192.168.80.148:11452
```
