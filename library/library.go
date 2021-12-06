package library

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/piksonGit/catch_unsplash_wallhaven_remake/database"
	"github.com/piksonGit/pmongo/db"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

type ICatcher interface {
	Catch(url string, fileName string, storageName string)
	Store()
	SetStore()
}
type ImgInfo struct {
}
type Catcher struct {
	c *colly.Collector
}

var idb database.IDB
var config map[string]string

func init() {
	idb = new(database.MongoDB)
	config = db.ReadJson("./config.json")
	idb.SetConnectString(config["mongo_connect_string"], config["mongo_db_name"], config["mongo_collection_name"])
}

func (cat *Catcher) Catch(url string, fileName string, storageName string) {

	cat.c.Visit(url)

}

func (cat *Catcher) Store() {}

func (cat *Catcher) SetStore() {
	fmt.Println("Set Store")
}
func (cat *Catcher) GetClone() (*colly.Collector, error) {

	clone := cat.c.Clone()
	if cat.c == nil {
		return clone, errors.New("原始收集器未定义")
	} else {
		return clone, nil
	}
}

func (cat *Catcher) GetCollector() *colly.Collector {
	return cat.c
}

func (cat *Catcher) SetCollector() {

	cat.c = colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
		colly.IgnoreRobotsTxt(),
		colly.MaxBodySize(100000000),
	)
	cat.c.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   20 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})
	extensions.RandomUserAgent(cat.c)
}
