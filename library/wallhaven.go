package library

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/piksonGit/catch_unsplash_wallhaven_remake/store"
	"go.mongodb.org/mongo-driver/bson"
)

type Wallhaven struct {
	c_detail *colly.Collector
	c_image  *colly.Collector
	c_login  *colly.Collector
	Catcher
}

func (wall *Wallhaven) SetLimit() {
	wall.c.Limit(&colly.LimitRule{
		DomainRegexp: `wallhaven\.cc`,
		//Delay:        10 * time.Second,
		RandomDelay: 10 * time.Second,
		Parallelism: 12,
	})
}

func (wall *Wallhaven) GetRequestMethod() func(r *colly.Request) {
	return func(r *colly.Request) {
		r.Headers.Set("referer", "https://wallhaven.cc/")
		r.Headers.Set("cookie", config["wallhaven_cookie"])
	}
}

//GetHTMLMEthod分析出图片链接和图片详情页链接，然后分发给不同的收集器去处理
func (wall *Wallhaven) GetHTMLMethod() func(e *colly.HTMLElement) {
	return func(e *colly.HTMLElement) {

		id := e.Attr("data-wallpaper-id")
		idhead := id[0:2]
		fmt.Println("catching,", idhead)
		imageurl := fmt.Sprintf("https://w.wallhaven.cc/full/%s/wallhaven-%s.jpg", idhead, id)
		detailurl := fmt.Sprintf("https://wallhaven.cc/w/%s", id)
		ctx := colly.NewContext()
		ctx.Put("detailurl", detailurl)
		ctx.Put("id", id)
		wall.c_image.Request("GET", imageurl, nil, ctx, nil)

	}
}
func (wall *Wallhaven) GetDetailResponseMethod() func(r *colly.Response) {

	return func(r *colly.Response) {

	}
}

func (wall *Wallhaven) GetDetailHTMLMethod() func(e *colly.HTMLElement) {
	return func(e *colly.HTMLElement) {
		wh := e.ChildText("h3")
		wh_arr := strings.Split(wh, "x")
		tags := e.ChildText(`div[data-storage-id="showcase-tags"]`)
		id := e.Request.Ctx.Get("id")
		filename := fmt.Sprintf("%d.jpg", id)
		info := bson.M{
			"filename":       filename,
			"altdescription": tags,
			"category":       "anime",
			"width":          wh_arr[0],
			"height":         wh_arr[1],
			"source":         "wallhaven.cc",
			"pageurl":        e.Request.URL,
		}
		fmt.Println("获取详情回调运行", id)
		idb.SaveImgInfo(bson.M{
			"id": id,
		}, bson.M{"$set": info})
		//imgTable.InsertOne(info)
	}
}
func (wall *Wallhaven) GetImageScrapMethod() func(r *colly.Response) {
	return func(r *colly.Response) {
		wall.c_detail.Request("GET", r.Ctx.Get("detailurl"), nil, r.Ctx, nil)

	}
}
func (wall *Wallhaven) GetImageResponseMethod() func(r *colly.Response) {
	return func(r *colly.Response) {
		id := r.Ctx.Get("id")
		path := config["wallhaven_filestore_path"]
		filepath := fmt.Sprintf(path+"/%s.jpg", id)
		Istore := &store.FileStore{r}
		Istore.Save(filepath)
		//r.Save(filepath)

	}
}

func (wall *Wallhaven) GetOnError() func(r *colly.Response, e error) {
	return func(r *colly.Response, e error) {
		log.Println(e.Error(), r.Request.URL)
	}
}

func (wall *Wallhaven) Catch(url string, filename string, storagename string) {

	wall.SetCollector()
	var err error
	wall.c_image, err = wall.GetClone()
	if err != nil {
		fmt.Println(err)
	}
	wall.c_detail, err = wall.GetClone()
	if err != nil {
		fmt.Println(err)
	}
	wall.c_login, err = wall.GetClone()
	if err != nil {
		fmt.Println(err)
	}
	wall.SetLimit()
	wall.c.OnRequest(wall.GetRequestMethod())
	wall.c.OnHTML(`figure[data-wallpaper-id]`, wall.GetHTMLMethod())
	wall.c_image.OnResponse(wall.GetImageResponseMethod())
	wall.c_image.OnResponse(wall.GetImageScrapMethod())
	//wall.c_detail.OnHTML(`.sidebar-content`, wall.GetDetailHTMLMethod())
	wall.c.Visit(url)
	//	wall.c.Wait()

}
