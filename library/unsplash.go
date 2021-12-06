package library

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gocolly/colly"
	"github.com/piksonGit/catch_unsplash_wallhaven_remake/model"
	"github.com/piksonGit/catch_unsplash_wallhaven_remake/store"
	"go.mongodb.org/mongo-driver/bson"
)

type Unsplash struct {
	c_detail *colly.Collector
	Catcher
}

func (up *Unsplash) SetLimit() {
	up.c.Limit(&colly.LimitRule{
		DomainRegexp: `unsplash\.com`,
		//Delay:        10 * time.Second,
		RandomDelay: 10 * time.Second,
		Parallelism: 12,
	})
}
func (up *Unsplash) GetRequestMethod() func(r *colly.Request) {
	return func(r *colly.Request) {
		r.Headers.Set("referer", "https://unsplash.com/")
	}
}

//GetResponseMethod是分析unsplash的目录结构，然后循环调用下载
func (up *Unsplash) GetResponseMethod() func(e *colly.Response) {

	return func(r *colly.Response) {
		var result model.Result
		err := json.Unmarshal(r.Body, &result)
		if len(r.Body) < 10 {
			os.Exit(1)
		}
		if err != nil {
			fmt.Println("到头了")
			fmt.Println(err)
			os.Exit(1)
		}
		for _, v := range result.Results {
			ctx := colly.NewContext()
			ctx.Put("info", v)
			up.c_detail.Request("GET", v.Urls.Raw, nil, ctx, nil)
		}
	}
}

//断点续传，排重。
//store类别
//保存在本地或者七牛，如果是保存到七牛，可以直接调用一个方法直接上传到七牛。保存到本地就是r.save()，不论是本地还是七牛，都得从数据库中记录下文件的信息。
func (up *Unsplash) GetDetailResponseMethod() func(e *colly.Response) {
	return func(r *colly.Response) {

		info := r.Ctx.GetAny("info")

		//fmt.Println(info)
		data, err := bson.Marshal(info)
		if err != nil {
			fmt.Println(err)
		}
		bm := bson.M{}

		err1 := bson.Unmarshal(data, bm)
		if err1 != nil {
			fmt.Println(err1)
		}
		id := bm["id"]
		fmt.Println("即将保存的图片的id是", id)
		//mongo查重。

		condition := bson.M{
			"id": bm["id"],
		}
		duplicateNum := idb.CheckDuplicate(condition)
		r.Request.Ctx.Put("documentcount", duplicateNum)
		r.Request.Ctx.Put("currentid", id)
		if duplicateNum == 0 {
			path := fmt.Sprintf("/Users/peterq/code/golang/imagess/%s.jpg", id)
			filename := fmt.Sprintf("%s.jpg", id)
			r.Request.Ctx.Put("filename", filename)
			Istore := &store.FileStore{r}
			Istore.Save(path)
			//r.Save(path) //或者可以选择直接保存到七牛。

		} else {
			fmt.Println("id为此的文件已经被保存，这个时候需要在数据库里判断是不是需要在category字段添加新的分类词")
		}
	}
}
func (up *Unsplash) getOnScrapMethod() func(*colly.Response) {

	return func(r *colly.Response) {
		info := r.Ctx.GetAny("info")
		data, err := bson.Marshal(info)
		if err != nil {
			fmt.Println(err)
		}
		num := r.Request.Ctx.GetAny("documentcount")
		id := r.Request.Ctx.Get("currentid")
		bm := bson.M{}
		err1 := bson.Unmarshal(data, bm)
		var newCategory string
		if err1 != nil {
			fmt.Println("bson解码出错:", err1)
		}
		if num.(int64) != 0 {
			condition := bson.M{
				"id": id,
			}
			idb.AppendCategory(condition, config["unsplash_category"])
		} else {
			condition := bm
			newCategory = config["unsplash_category"]
			upItem := bson.M{
				"$set": bson.M{
					"filename": r.Request.Ctx.Get("filename"),
					"category": newCategory,
					"path":     newCategory,
				},
			}
			idb.SaveImgInfo(condition, upItem)
		}

	}

}
func (up *Unsplash) Catch(url string, fileName string, storageName string) {
	up.SetCollector()
	var err error
	up.c_detail, err = up.GetClone()
	if err != nil {
		fmt.Println(err)
	}
	up.SetLimit()
	//Visit 方法一定要在设置回调方法之后再调用。
	up.c.OnRequest(up.GetRequestMethod())
	up.c.OnResponse(up.GetResponseMethod())
	up.c_detail.OnResponse(up.GetDetailResponseMethod())
	up.c_detail.OnScraped(up.getOnScrapMethod())
	up.c.Visit(url)

	//up.c.Error()
}
