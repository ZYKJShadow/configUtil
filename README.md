# configUtil
Go语言解析配置文件的工具

## 使用说明
将文件夹放到`GO_PATH/src/`下，在头部`import`导入,

## 使用示例
配置文件示例：
```ini
[mysql]
address=localhost
port=3306
username=root
password=root

[redis]
host=localhost
port=6379
password=root
database=0
```
### main函数中使用
```go
//根据配置文件的名称，并取好别名，要跟配置文件的名称对应
type Mysql struct {
	Address  string `ini:"address"`
	Port     int    `ini:"port"`
	Username string `ini:"username"`
	Password string `ini:"password"`
}

type Redis struct {
	Host     string `ini:"host"`
	Port     int    `ini:"port"`
	Password string `ini:"password"`
	Database int    `ini:"database"`
}

//最后集体封装到一个struct中
type Config struct {
	Mysql `ini:"mysql"`
	Redis `ini:"redis"`
}

func main() {
  //将封装好的struct传入
var config Config
//参数1：配置文件的路径，指的是GO_PATH/src/下的路径，例如GO_PATH/src/goProject/config.ini，则传入如下所示
//参数2：结构体指针，一定要传指针
configUtil.LoadConfig("goProject/config.ini", &config)
}
```

