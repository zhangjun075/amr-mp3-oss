# amr-mp3-oss

## 本程序主要是从阿里云下载固定前缀的amr文件
* 下载文件
* 转换文件
* 上传文件

## Mac 下编译 Linux 和 Windows 64位可执行程序
```$xslt
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go
```
## 日期格式化
```$xslt
fmt.Println(time.Now().Format("2006-01-02 15:04:05")) # 这是个奇葩,必须是这个时间点, 据说是go诞生之日, 记忆方法:6-1-2-3-4-5
```
