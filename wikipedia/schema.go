package wikipedia

type RArticleResp struct {
	Query RQuery `json:"query"`
}
type RQuery struct {
	Pages map[string]Page `json:"pages"`
}
type Page struct {
	Title string `json:"title"`
}
