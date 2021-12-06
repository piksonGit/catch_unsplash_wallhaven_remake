//database包提供给爬虫一个保存到数据库的工具.
package database

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/piksonGit/pmongo/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IDB interface {
	SaveImgInfo(condition, item bson.M)
	UpdateImgInfo(condition, item bson.M)
	CheckDuplicate(condition bson.M) (num int64)
	SetConnectString(url, dbname, collectName string)
	AppendCategory(condition bson.M, category string) string
}

type MongoDB struct {
	collection db.Col
	fileName   string
}

func init() {

}

//SetConnectString是设置mongo的连接url
func (mongo *MongoDB) SetConnectString(url string, dbname string, collectionName string) {

	mongo.collection = db.Conn(url, dbname, collectionName)

}

//SaveImgInfo是把图片信息保存到mongo数据库里。
func (mongo *MongoDB) SaveImgInfo(condition, item bson.M) {

	var updateOptions *options.UpdateOptions
	updateOptions = &options.UpdateOptions{}
	updateOptions.SetUpsert(true)
	mongo.collection.UpdateOne(condition, item, updateOptions)

}

//UpdateImgInfo是更新图片的信息
func (mongo *MongoDB) UpdateImgInfo(condition, item bson.M) {
	mongo.collection.UpdateOne(condition, item)
}

//CheckDuplicate是查重
func (mongo *MongoDB) CheckDuplicate(condition bson.M) (num int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	num, err := mongo.collection.CountDocuments(ctx, condition)
	defer cancel()
	if err != nil {
		log.Println(err)
	}
	return num
}

//AppendCategory追加新的分类到已经存在的集合
func (mongo *MongoDB) AppendCategory(condition bson.M, category string) string {
	result := mongo.collection.FindOne(condition)
	oldCategory := result["category"].(string)

	if strings.Contains(oldCategory, category) {
		return oldCategory
	}

	newCategory := oldCategory + " " + category
	yes := true
	mongo.collection.UpdateOne(condition, bson.M{
		"$set": bson.M{
			"category": newCategory,
		},
	}, &options.UpdateOptions{
		Upsert: &yes,
	})
	return newCategory
}

//里式替换原则，开闭法则，迪米特法则，单一职责原则，接口隔离原则，
