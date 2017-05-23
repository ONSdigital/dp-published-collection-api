package publishedcollection

type PublishedCollection struct {
	CollectionID     string          `json:"id"`
	CollectionName   string          `json:"name"`
	PublishDate      string          `json:"publishDate"`
	PublishStartDate string          `json:"publishStartDate"`
	PublishEndDate   string          `json:"publishEndDate"`
	Results          []PublishedItem `json:"publishResults"`
}

type PublishedItem struct {
	Duration int64  `json:"duration"`
	Uri      string `json:"uri"`
	Size     int64  `json:"size"`
}
