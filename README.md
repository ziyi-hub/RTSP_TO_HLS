# RTSP_TO_HLS

## 安装GoLang环境
### 配置国内镜像：
- go env -w GOPROXY=https://goproxy.cn,direct
### 切换到路径为 RTSP_TO_HLS/golang 文件夹下运行：
- <code> go run main.go http.go stream.go config.go </code> 

## 安装 node 和 npm
### 在 back 和 front 含有 package.json 路径下里运行 
- <code> npm i </code>
- 前端front运行 npm serve
- 后端back运行 node server.js

## postman 获取 token
- 后端成功运行后 postman GET <code>localhost:3000/login</code>
- body / raw / json 输入，随便一个名字
<code>
{
    "name": "Ziyi"
}
</code>
- postman 会返回一个 
"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7fSwiaWF0IjoxNzQ5NzAyOTE0LCJleHAiOjE3NDk3MDM1MTR9.X5S973bFmu_bwszuOnGFXLOI2rXCrloJX2-uzmf7fiw"

- 复制 eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7fSwia 到golang 的 config.go 的141行 "Bearer eyJhbGciOiJIUzI1NiIsI" 替换一下 Bearer 后面的token
- 然后wins运行 <code> go env -w GOPROXY=https://goproxy.cn,direct
           go run main.go http.go stream.go config.go </code>

- mac 运行 go run *.go 开启转流golang

### 注意开启顺序
- back -> golang -> front

### 全部开启后 golang 转流地址用于验证视频流能否使用
- http://localhost:8083/player/JSESSIONID=86F633899E8C2CF1577717

### 全部开启golang 生成的hls地址
- http://localhost:8083/play/hls/JSESSIONID=86F633899E8C2CF1577717/index.m3u8

### 前端VUE集成 
- http://localhost:8080/
- 集成位置为APP.VUE



