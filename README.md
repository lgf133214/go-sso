## go 单点登录

### 使用流程
    
    - 1、获取应用的 app-id，目前为手动添加

    - 2、应用内如需登录，跳转至本服务下 /app-id/:app-id/login 页面，携带跳转前页面链接
        - GET https://xxx/app-id/:app-id/login?redirect=
            - 未登录（无 Cookie 或过期时间 < 5min）重新登录
                - POST https://xxx/app-id/:app-id/login 验证密码
                - 成功，设置 Cookie，跳转
            - 已登录，未过期，跳转至 app-redirect-url?st=&redirect=

    - 3、验证成功后携带 st（service ticket） 跳转至应用提供的 url（与 app-id 在一个表里，目前手动添加）

    - 4、应用程序收到 st 参数后，发送验证请求至本服务

    - 5、验证成功后，返回用户信息

    - 6、应用程序添加对应的 Cookie 返回客户端

### 