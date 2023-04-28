# Go makefile

# 常量设置
PROJECT_NAME=FiveCDNTurbo
PACKAGE_VERSION=V1.0

# Go 工具设置
GO=go
GOFMT=gofmt
GOBUILD=$(GO) build
GOCLEAN=$(GO) clean

# 源文件和目录
SRC=./cmd
BIN=./bin/$(PROJECT_NAME)

all: build

.PHONY: build
build:
	@echo "Building $(PROJECT_NAME) version $(PACKAGE_VERSION)..."
	$(GOBUILD) -o $(BIN) $(SRC)

.PHONY: clean
clean:
	@echo "Cleaning up..."
	$(GOCLEAN)
	rm -rf $(BIN)

.PHONY: fmt
fmt:
	@echo "Formatting source code..."
	$(GOFMT) -w $(SRC)

.PHONY: package
package:
	@echo "Packaging $(PROJECT_NAME) version $(PACKAGE_VERSION)..."
	sh ./package.sh $(PROJECT_NAME) $(PACKAGE_VERSION)