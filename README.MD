# FiveCDN 

FiveCDN 所使用的全套系统（加速系统-Turbo 用户操作页面（User_Actions））

## Turbo
加速系统的主要程序，基于 Go Gin 框架，支持分路径代理 路径黑白名单 黑白名单自动更新 等基础功能

### 获取
您有两种获取主程序的方式
#### 1.构建
在 Bash Golang Git 等环境存在时，
先 git clone https://github.com/ahmr-bot/FiveCDN.git
然后 cd Turbo
直接使用 ./package.sh 进行构建
#### 2.直接下载
前往<https://github.com/ahmr-bot/FiveCDN/releases> 直接下载
支持 Windows/Linux 386 amd64 arm arm64 架构

### 配置
```toml
[server]
host = "localhost"
port = 8080
name = "广东移动"
powered_by = "Five CDN Turbo V1.0"
# 服务器相关

[whitelist]
enabled = true
url = "https://mecdn.mcserverx.com/gh/ahmr-bot/MECDN-WhiteList/master/list.txt"
update_interval = "5m"
# 白名单相关

# 配置需要代理的域名和对应的路径
[[proxy]]
domain = "a.com"
path = "/a"
[[proxy]]
domain = "b.com"
path = "/b/article"

[debug]
debug= false
# 调试模式
```
### 启动
使用  ./Turbo-xxxxx -c config.toml  即可启动

## User_Actions
该程序是 白名单提交、URL 刷新（使用 cdnfly api）、人工审核 三部分组成
基于 Python Flask 框架
### 获取
直接克隆仓库并打开 User_Action 文件夹即可 

### 配置
main.py:
``` python
# 限制单IP每日10次
MAX_REQUESTS_PER_DAY = 10

# 用于生成API签名的密钥对
API_KEY = 'xxxxxxxxx'
API_SECRET = 'xxxxxxxxxxx'

# 缓存刷新API的URL
API_URL = 'https://xxxxxxxxxxxx/v1/jobs'
```
其他文件无需配置

### 使用 
main.py / audit.py 直接运行即可
main.py：白名单提交、URL 刷新（使用 cdnfly api）
audit.py：人工审核程序