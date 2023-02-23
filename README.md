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
# douyin 服务的监听地址
serve-addr: 192.168.80.148:11451
douyin:
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
    # 头像和背景的设置，在注册用户时会从数组里随机选择一个文件
    # 在访问时，会分别请求 /storage/avatars/:file 和 /storage/backgrounds/:file
    # 所以在这里设置后，在 storage 对应的文件夹里要创建对应的文件，否则无法访问
    avatars:
        - 1.jpg
        - 2.jpg
        - 3.jpg
    backgrounds:
        - 1.jpg
        - 2.jpg
        - 3.jpg
# 这是存储服务，只用于存储文件
storage:
    # 存储数据的根文件夹
    data-dir: Data
    # 服务的监听地址，用于拼接 url，见 service/storage.go 105 行
    pre-url: example.com:11451
    # 请求 /storage 时会计算 token + 文件路径 的 md5 对文件链接进行混淆
    # 例如请求文件 videos/a.mp4 而在请求时可能就会变成 6468737561696b766....
    # 从而不会将存储的文件路径暴露在外
    token: 这个值就是随便写，留空也行
```
