#!/bin/bash

PROJECT=Turbo
VERSION=$1
BUILD_DIR=./build
PKG_DIR=${BUILD_DIR}/${PROJECT}-${VERSION}
TAR_FILE=${BUILD_DIR}/${PROJECT}-${VERSION}.tar.gz

# 操作系统和架构列表
OS_ALL="linux windows"
ARCH_ALL="386 amd64 arm arm64"

# 创建打包目录
mkdir -p ${PKG_DIR}

# 循环编译
for OS in ${OS_ALL}; do
    for ARCH in ${ARCH_ALL}; do
        # 根据操作系统和架构生成文件名
        FILENAME=${PROJECT}_${OS}_${ARCH}

        # 设置 GOOS 和 GOARCH 环境变量
        export GOOS=${OS}
        export GOARCH=${ARCH}

        # 编译可执行文件
        echo "Building ${FILENAME}..."
        go build -o ${PKG_DIR}/${FILENAME} ./cmd

        # 复制配置文件到打包目录
        cp ./config.toml ${PKG_DIR}/

        # 创建 tar 包
        tar -czf ${BUILD_DIR}/${PROJECT}-${OS}-${ARCH}-${VERSION}.tar.gz -C ${BUILD_DIR} ${PROJECT}-${VERSION}

        # 清理打包目录
        rm -rf ${PKG_DIR}
    done
done