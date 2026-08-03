# 博客后端

## 微信小程序登录部署

微信小程序使用 `POST /api/wechat/login` 将 `wx.login` 返回的临时 `code` 发送给后端。后端向微信 `code2Session` 服务换取 openid，首次登录自动创建用户并签发项目现有的 JWT。

### 配置

在服务器的 `settings.yaml` 中添加微信小程序 AppID，AppSecret 不要写入仓库：

```yaml
wechat:
  app_id: "你的微信小程序 AppID"
  app_secret: ""
```

通过服务器环境变量提供 AppSecret：

```env
WECHAT_APP_SECRET=你的微信小程序AppSecret
```

`.env.example` 仅为变量名示例，不能填入真实密钥后提交。服务端不会保存或返回微信的 `session_key`。

### 发布前检查

1. 后端首次带此版本启动时，`AutoMigrate` 会为 `user` 表增加唯一的 `wechat_open_id` 字段。
2. 在微信公众平台将生产环境 HTTPS 域名添加到“request 合法域名”。微信小程序正式版不能请求 HTTP 地址。
3. 当前 IP 地址仅用于本地开发调试；购买域名并部署 HTTPS 后，需要同步修改小程序 `services/config.js` 的 `BASE_URL`。
4. 使用真实设备完成一次评论：首次评论应触发微信登录并自动建号，随后评论复用本地 JWT；点赞不需要登录。
