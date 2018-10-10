package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"time"
)

const (
	AudioSamplingRateMP3  = "22050"
	//myfolder = "/Users/junzhang/Desktop/temp/go/"
	myfolder = "/home/pdt_test_caoyang02/oss/"

)
var bucket *oss.Bucket
var todayFolder string

func main() {
	//创建下目录，今天的放在今天的
	mkTodayDir()
	//println(strings.Replace("audio/000c8725-db0d-473d-bb1d-d81edbe3dfaf.amr","amr","mp3",1))
	bucket = bucketo()
	getAndDownload()
	fileNames := listfile()
	for _,filename := range fileNames {
		//println(filename)
		convertToMp3(filename)
	}

}

//创建当天的目录
func mkTodayDir() {
	today := time.Now().Format("20060102")
	_,err2 := os.Stat(myfolder+today)
	if os.IsNotExist(err2) {
		err3 := os.Mkdir(myfolder+today,os.ModePerm)
		if err3 == nil {
			println("mkdir"+myfolder+today+" success..")
		}else {
			println("mkdir"+myfolder+today+" fail..")
		}

	}else {
		println("dir "+myfolder+today+" exists.")
	}
	todayFolder = myfolder+today +"/"

}


func convertToMp3(filenames ...string) error {
	var tofilename string
	fromfilename := filenames[0]

	switch len(filenames) {
	case 1:
		tofilename = filenames[0]
		break
	case 2:
		tofilename = filenames[1]
		break
	default:
		tofilename = filenames[0]
	}
	comm := exec.Command("ffmpeg", "-i", fromfilename+".amr", "-ar", AudioSamplingRateMP3, tofilename+".mp3")
	//判断mp3文件是否存在
	fileInfo,err1 := os.Stat(tofilename+".mp3")
	if fileInfo != nil && err1 == nil {
		println(tofilename+".mp3"+"文件存在，不转换...")
		return nil
	}else {
		println(fromfilename+".amr"+"转换....")
		if err := comm.Run(); err != nil {
			return err
		}
		println("=>"+tofilename+".mp3"+"转换成功。")
	}

	println("=>"+tofilename+".mp3"+"开始上传...")
	//上传本地文件。
	arr := strings.Split(tofilename,"/")
	fileName := arr[len(arr)-1]

	// 判断文件是否存在。
	isExist, err := bucket.IsObjectExist("audio/"+fileName+".mp3")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	if isExist {
		println("audio/"+fileName+".mp3 文件存在,不上传.")
	}else {
		err1 = bucket.PutObjectFromFile("audio/"+fileName+".mp3", tofilename+".mp3",oss.Progress(&OssProgressListener{}))
		if err1 != nil {
			fmt.Println("Error:", err1)
			os.Exit(-1)
		}
		fmt.Println("upload Completed.")
	}



	return nil
}

func listfile() []string {
	files,_ := ioutil.ReadDir(todayFolder)
	var fileArray []string
	for _,file := range files {
		if file.IsDir(){
			continue
		}else {
			fileName := strings.Split(file.Name(),".")[0]
			filenameSuffix := strings.Split(file.Name(),".")[1]

			if len(fileName)>0 && filenameSuffix == "amr" {
				fileArray = append(fileArray,todayFolder+fileName)
			}
		}
	}
	return fileArray
}

// 定义进度条监听器。
type OssProgressListener struct {
}
// 定义进度变更事件处理函数。
func (listener *OssProgressListener) ProgressChanged(event *oss.ProgressEvent) {
	switch event.EventType {
	case oss.TransferStartedEvent:
		fmt.Printf("Transfer Started, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	case oss.TransferDataEvent:
		fmt.Printf("\rTransfer Data, ConsumedBytes: %d, TotalBytes %d, %d%%.",
			event.ConsumedBytes, event.TotalBytes, event.ConsumedBytes*100/event.TotalBytes)
	case oss.TransferCompletedEvent:
		fmt.Printf("\nTransfer Completed, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	case oss.TransferFailedEvent:
		fmt.Printf("\nTransfer Failed, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	default:
	}
}

func bucketo() *oss.Bucket  {
	// 创建OSSClient实例。
	client, err := oss.New("xxx", "xxx", "xxx")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	// 获取存储空间。
	bucket, err := client.Bucket("paifenle-tars")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	return bucket
}

func getAndDownload() {
	d, _ := time.ParseDuration("-24h")
	// 分页列举包含指定前缀的文件。每页列举80个。
	prefix := oss.Prefix("audio/")
	marker := oss.Marker("")
	for {
		lsRes, err := bucket.ListObjects(oss.MaxKeys(80), marker, prefix)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(-1)
		}
		prefix = oss.Prefix(lsRes.Prefix)
		marker = oss.Marker(lsRes.NextMarker)
		// 打印结果。
		for _, object := range lsRes.Objects {
			//列举所有aac文件：audio/00005d64-783c-40da-8f90-05b534a2b353.aac
			objectName := object.Key
			if objectName != "audio/" && strings.Contains(objectName,".amr") {
				if object.LastModified.Format("20060102") != time.Now().Add(d).Format("20060102") {
					println("跳过...")
					continue
				}

				fmt.Println("-Object:", objectName)
				objectNameSuffix := strings.Split(objectName,"/")[1]

				// 下载文件到本地文件。
				err = bucket.GetObjectToFile(objectName, todayFolder+objectNameSuffix,oss.Progress(&OssProgressListener{}))
				if err != nil {
					fmt.Println("Error:", err)
					os.Exit(-1)
				}
				fmt.Println("Transfer Completed.")

			}

		}

		if !lsRes.IsTruncated {
			break
		}
	}

}
