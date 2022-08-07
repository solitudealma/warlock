## 使用说明

```
# 克隆项目
git clone https://github.com/solitudealma/warlock.git

# 进入server文件夹
cd warlock

# 使用 go mod 并安装go依赖包
go generate

# 编译 
go build -o main main.go (windows编译命令为go build -o main.exe main.go )
chomd +x compress_game_js.sh && ./compress_game_js.sh

# 运行二进制
./main (windows运行命令为 main.exe)
```
