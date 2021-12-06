package store

import (
	"fmt"

	"github.com/gocolly/colly"
	"github.com/piksonGit/pmongo/db"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type Istore interface {
	Save(url string) string
}

type QiniuStore struct {
	filepath     string
	bucketname   string
	bucketdomain string
	docname      string
}

//setFilePath 设置文件在七牛的前缀
func (qiniu *QiniuStore) setFilePath(filepath string) {
	qiniu.filepath = filepath
}

//setBucketName 设置存储的bucket
func (qiniu *QiniuStore) setBucketName(bucketname string) {
	qiniu.bucketname = bucketname
}

//setBucketDomain设置bucket的域名。
func (qiniu *QiniuStore) setBucketDomain(bucketdomain string) {
	qiniu.bucketdomain = bucketdomain
}

//setDocName设置文件名
func (qiniu *QiniuStore) setDocName(docname string) {
	qiniu.docname = docname
}

//init初始化设置前缀，bucket名，bucket域名，文件名。
func (qiniu *QiniuStore) init(filepath, bucketname, bucketdomain, docname string) {

	qiniu.setBucketDomain(bucketdomain)
	qiniu.setBucketName(bucketname)
	qiniu.setDocName(docname)
	qiniu.setFilePath(filepath)

}

//Save方法是保存到七牛的方法，在传递到save之前需要对一些参数进行初始化，否则无法正确上传到七牛
func (qiniu *QiniuStore) Save(url string) {
	qiniu.save(url)
}

//适配器
func (qiniu *QiniuStore) save(url string) (info string) {
	config := db.ReadJson("./config.json")
	filename := qiniu.filepath + qiniu.docname
	accessKey := config["qiniu_key"]
	secretKey := config["qiniu_secret"]
	mac := qbox.NewMac(accessKey, secretKey)
	cfg := storage.Config{
		UseHTTPS: true,
	}

	bucketManager := storage.NewBucketManager(mac, &cfg)
	bucket := qiniu.bucketname
	resURL := url
	fetchRet, err := bucketManager.Fetch(resURL, bucket, filename)
	if err != nil {
		fmt.Println("fetch error,", err)
	}

	return qiniu.bucketdomain + fetchRet.Key
}

type FileStore struct {
	R *colly.Response
}

func (filestore *FileStore) Save(filepath string) {
	filestore.R.Save(filepath)
}
