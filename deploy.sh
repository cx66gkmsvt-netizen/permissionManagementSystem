#!/bin/bash

# 用户中心系统 Docker 部署脚本
# 使用方法: ./deploy.sh [build|start|stop|restart|logs|clean]

set -e

APP_NAME="user-center"

# 构建镜像
build() {
    echo "🔨 Building Docker images..."
    docker-compose build --no-cache
    echo "✅ Build completed!"
}

# 启动服务
start() {
    echo "🚀 Starting services..."
    docker-compose up -d
    echo "✅ Services started!"
    echo ""
    echo "📍 Access the application:"
    echo "   Frontend: http://localhost:8081"
    echo "   Backend API: http://localhost:8080/api"
    echo ""
    docker-compose ps
}

# 停止服务
stop() {
    echo "🛑 Stopping services..."
    docker-compose down
    echo "✅ Services stopped!"
}

# 重启服务
restart() {
    stop
    start
}

# 查看日志
logs() {
    docker-compose logs -f
}

# 清理
clean() {
    echo "🧹 Cleaning up..."
    docker-compose down -v --rmi all
    echo "✅ Cleanup completed!"
}

# 帮助信息
usage() {
    echo "Usage: $0 {build|start|stop|restart|logs|clean}"
    echo ""
    echo "Commands:"
    echo "  build    - Build Docker images"
    echo "  start    - Start all services"
    echo "  stop     - Stop all services"
    echo "  restart  - Restart all services"
    echo "  logs     - View logs"
    echo "  clean    - Remove containers, volumes, and images"
}

# 主逻辑
case "$1" in
    build)
        build
        ;;
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        restart
        ;;
    logs)
        logs
        ;;
    clean)
        clean
        ;;
    *)
        usage
        exit 1
        ;;
esac
