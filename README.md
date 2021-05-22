# godaemon
帮助go程序在后台运行或者杀死在后台运行的go程序

## 使用

在main函数所在的文件中加上
```go
_ "github.com/markwinds/godaemon"
```

## 示例

- 代码

```go
package main

import (
	_ "github.com/markwinds/godaemon"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("hello, golang!\n"))
	})
	log.Fatalln(http.ListenAndServe(":7070", mux))
}
```

- 后台运行

```shell
#windows
main.exe -d
#linux
./main -d
```

- 查看程序是否正常运行：打开浏览器输入网址127.0.0.1:7070可以看到页面显示hello, golang!

- 杀死后台运行程序

```shell
#windows
main.exe -k
#linux
./main -k
```
