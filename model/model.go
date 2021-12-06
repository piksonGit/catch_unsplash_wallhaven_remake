package model

type Item struct {
	Id             string       `json:"id"`
	Width          int          `json:"width"`
	Height         int          `json:"height"`
	Altdescription string       `json:"alt_description"`
	Categories     []string     `json:"categories"`
	Color          string       `json:"color"`
	Totallikes     int          `json:"total_likes"`
	Urls           Downloadurls `json:"urls"`
}
type Downloadurls struct {
	Full    string `json:"full"`
	Raw     string `json:"raw"`
	Regular string `json:"regular"`
	Small   string `json:small`
	Thumb   string `json:thumb`
}
type Result struct {
	Total      int    `json:"total"`
	Totalpages int    `json:"total_pages"`
	Results    []Item `json:"results"`
}
