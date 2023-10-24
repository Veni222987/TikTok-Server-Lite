一、项目介绍
抖声提供了包括视频流、视频投稿和个人主页等基础的短视频APP功能。在扩展功能方面，互动方向支持用户点赞和评论列表查看；而社交方向涵盖了关注、粉丝列表以及消息系统等功能。该项目采用了Gin+GORM框架来构建服务。数据库方面选择了MySQL，并且还利用了Redis来缓存token信息。通过这些选择，项目实现了基础服务，并且还包括了社交和互动功能。

二、项目实现
3.1 技术选型与相关开发文档
3.1.1 场景分析
场景一：浏览视频 Feed 流
- 角色：所有用户
- 主要目标：浏览最新的视频内容
- 细节：用户在应用首页查看按投稿时间排列的视频 Feed 流，浏览视频、点赞、进入个人主页等。
场景二：发布视频
- 角色：登录用户
- 主要目标：上传自己的视频并分享给其他用户
- 细节：用户在个人主页点击发布按钮，录制或上传视频，添加描述和标签，将视频信息显示在个人主页和首页 Feed 流。
场景三：互动和社交
- 角色：登录用户
- 主要目标：与其他用户进行互动和社交交流
- 细节：用户可以点赞、评论视频，查看喜欢列表，关注其他用户，查看消息，与关注的用户进行私信交流。
场景四：查看个人主页
- 角色：所有用户
- 主要目标：查看用户的基本信息和投稿列表
- 细节：用户点击自己头像进入个人主页，查看个人信息、关注数、粉丝数以及发布的视频列表。
场景五：管理关注和粉丝列表
- 角色：登录用户
- 主要目标：查看和管理关注的用户和粉丝
- 细节：用户查看关注列表和粉丝列表，进入用户主页，取消关注用户，管理自己的关注关系。
场景六：发送消息
- 角色：登录用户
- 主要目标：与其他用户进行私信交流
- 细节：用户进入消息页面，点击用户头像进入聊天页面，发送文本消息、表情等与对方进行私信交流。
3.1.2 需要解决的问题：
1. 用户认证和安全性： 如何确保用户的身份认证以及数据的安全性，以防止未经授权的访问和信息泄露？
2. 视频上传和存储： 如何处理用户上传的视频，包括视频的压缩、存储、转码和展示，以及避免因视频存储而导致的性能问题？
3. 用户互动实时性： 如何实现用户点赞、评论、关注等互动行为的实时性，以保证用户体验？
4. 服务器和数据库扩展性： 如何设计服务器架构和数据库，以便在用户数量增加时保持稳定的性能？
3.1.3 前提假设：
1. 用户设备和网络连接： 用户使用的设备和网络连接具有足够的性能和带宽，以便流畅地浏览和上传视频。
2. 用户合法性： 用户提供的信息和内容是合法的，没有侵犯版权、隐私等问题。
3. 数据库和服务器可用性： 数据库和服务器稳定可靠，以保证系统的正常运行和数据的安全。
4. 用户体验期望： 用户对于浏览、互动、上传视频等操作有一定的耐心，但也期望在合理的时间内完成操作。
3.1.4 技术选型与相关文档
- 编程语言： Go语言
- Web框架： Gin
- 数据库： MySQL
- 缓存： Redis
- 对象存储： 使用云存储服务，阿里云 OSS 
相关开发文档：
- Go语言（Golang）： 官方网站是 https://golang.org/ 。
- Gin Web框架： 官方文档和代码可以在 GitHub 上找到：https://github.com/gin-gonic/gin 。
- MySQL数据库： 官方网站是 https://www.mysql.com/ 
- Redis缓存： 官方网站是 https://redis.io/ 。
3.2 架构设计
本项目采用了单体应用架构，在单体应用架构中，整个应用的所有功能模块都集成在一个单一的应用中。这种架构适用于小型项目和快速开发。
前端： 使用现成的抖声apk。
后端：
- 使用Go语言，搭配Gin框架构建后端服务，处理前端请求和业务逻辑。
- 实现用户认证和授权，确保安全的用户访问和操作。
- 设计并实现数据库模型，使用MySQL存储用户信息、视频信息、互动数据等。
- 实现用户互动、消息系统、关注功能等。
数据：使用MySQL作为主要的关系型数据库，存储用户数据、视频信息、互动数据等。
缓存： 使用Redis作为缓存层，提高数据的读取速度。
文件存储： 使用阿里云OSS存储用户上传的视频，处理视频存储和转码，减轻服务器负担。
3.3 项目代码介绍
项目的代码结构如下：
├─controller/
├─middleware/
├─model/
├─service/
├─sql/
├─test/
└─utils/
└─main.go
└─router.go
└─config.yaml
1. controller/：这个目录包含应用程序的控制器代码，负责处理用户请求并响应。
其中负责聊天的函数代码如下：
```go
func SendMessage(ctx *gin.Context) {
    actionType := ctx.Query("action_type")
    if len(ctx.Query("content")) == 0 {
       ctx.JSON(400, gin.H{
          "status_code": 1,
          "status_msg":  "content不能为空",
       })
       return
    }
    if actionType == "1" {
       uid := service.GetIdByToken(ctx.Query("token"))
       toUserId, err := strconv.ParseInt(ctx.Query("to_user_id"), 10, 64)
       if err != nil {
          panic(err)
       }
       content := ctx.Query("content")
       //插入数据库
       messageLog := model.Message{OriginUserId: uid, DestinationUserId: toUserId, Content: content, CreateDate: time.Now().UnixNano() / 1e6}
       if res := service.Db.Create(&messageLog); res.Error != nil {
          panic(res.Error)
       }
       ctx.JSON(200, gin.H{
          "status_code": 0,
          "status_msg":  "发送成功",
       })
    } else {
       ctx.JSON(400, gin.H{
          "status_code": 1,
          "status_msg":  "actionType错误",
       })
    }
```
controller使用了gin.Context对象处理HTTP请求的上下文信息，通过获取Context中的请求信息，经过业务代码，将返回结果也写入Context中。
3. middleware/：这个目录包含中间件代码，用于在到达控制器之前或之后执行一些操作，比如身份验证、日志记录等。
检测请求参数中的token可以注册使用如下中间件：
```go
func QueryAuthMiddleWare() gin.HandlerFunc {
    return func(ctx *gin.Context) {
       token := ctx.Query("token")
       if service.IsTokenExist(token) {
          service.RedisClient.Set(token, service.RedisClient.Get(token).Result, 86400000000000)
          ctx.Next()
       } else {
          fmt.Println("无效的token")
          ctx.AbortWithStatusJSON(401, gin.H{
             "error": "无效的Token",
          })
          return
       }
    }
}
```
4. model/：这个目录包含应用程序的模型代码，用于表示应用程序的数据结构和操作。
比如用户信息的结构体如下：
```go
type User struct {
    Id              int64  `gorm:"id" json:"id"`                             // 用户id
    Name            string `gorm:"name" json:"name"`                         // 用户名称
    FollowCount     int    `gorm:"follow_count" json:"follow_count"`         // 关注总数
    FollowerCount   int    `gorm:"follower_count" json:"follower_count"`     // 粉丝总数
    Avatar          string `gorm:"avatar" json:"avatar"`                     // 用户头像
    BackgroundImage string `gorm:"background_image" json:"background_image"` // 用户个人页顶部大图
    Signature       string `gorm:"signature" json:"signature"`               // 个人简介
    TotalFavorited  int    `gorm:"total_favorited" json:"total_favorited"`   // 获赞数量
    WorkCount       int    `gorm:"work_count" json:"work_count"`             // 作品数
    FavoriteCount   int    `gorm:"favorite_count" json:"favorite_count"`     // 喜欢数
}

func (*User) TableName() string {
    return "user"
}
```
5. service/：这个目录包含其他服务代码，包括数据库，Redis，OSS等配置和操作函数。
判断token是否存在的代码：
```go
// IsTokenExist 判断token是否存在
func IsTokenExist(token string) bool {
    result, err := RedisClient.Exists(token).Result()
    if err != nil {
       panic(err)
    }
    if result == 0 {
       return false
    }
    return true
}
```
6. sql/：这个目录包含用于数据库操作的SQL脚本，记录数据库表的基本结构和信息，这里就不举例了。
7. test/：这个目录包含应用程序的测试代码。测试用于验证应用程序的各个部分是否按照预期工作，确保代码的质量和稳定性。
8. utils/：这个目录包含一些工具函数或实用程序代码，通常用于辅助其他部分的开发，比如时间转换，GORM表结构生成等工具。
9. main.go：这是主程序文件，程序的入口点。它负责启动应用程序并配置其他模块。
10. router.go：这个文件包含应用程序的路由配置。
11. config.yaml：这个文件包含应用程序的配置信息，比如数据库连接、端口号、日志级别等。配置示例如下：
```yaml
database:
  driver: 
  host: 
  port:
  username: 
  password: 
  database:

redis:
  addr:
  password:
  DB:

OSS:
  endPoint:
  accessKey:
  accessSecret:
```
项目的全部代码可以在仓库中查看，这里就不再一一举例了。
