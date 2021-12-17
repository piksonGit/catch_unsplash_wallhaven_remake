package main

import (
	"fmt"
	"io/ioutil"
	"strconv"

	lib "github.com/piksonGit/catch_unsplash_wallhaven_remake/library"
	"github.com/piksonGit/pmongo/db"
)

func main() {

	/* 	var catcher lib.ICatcher
	   	catcher = new(lib.Unsplash)
	   	config := db.ReadJson("./config.json")
	   	page := 1
	   	url := fmt.Sprintf("https://unsplash.com/napi/search/photos?query=%s&per_page=20&page=%s&xp=", config["unsplash_category"], page)
	   	catcher.Catch(url, "", "") */

	var catcher lib.ICatcher
	catcher = new(lib.Wallhaven)
	config := db.ReadJson("./config.json")
	var i int
	//此处替换为config["wallhaven_last_page_num"]
	endnum, err := strconv.Atoi(config["wallhaven_last_page_num"])
	if err != nil {
		fmt.Println(err)
	}
	for i = 1; i <= endnum; i++ {
		url := fmt.Sprintf(config["wallhaven_catch_url"]+`&page=%d`, i)
		fmt.Println(url)
		catcher.Catch(url, "", "")
	}
	defer recordPage(i)

}
func recordPage(i int) {
	istring := strconv.Itoa(i)
	err := ioutil.WriteFile("./page.pid", []byte(istring), 0777)
	if err != nil {
		fmt.Println(err)
	}
}
