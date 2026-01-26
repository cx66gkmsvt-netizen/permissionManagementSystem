#!/bin/bash

# 设置环境变量 (Optional, defaults are in the binary)
# export DB_HOST=127.0.0.1
# export DB_PORT=3306
# export DB_USER=root
# export DB_PASSWORD=password
# export DB_NAME=user_center
# export SERVER_PORT=8080

# 赋予执行权限
chmod +x server

# 启动服务
./server
