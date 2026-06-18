package request

/*
	type Articles struct {
		Page     int    `json:"page"`
		PageSize int    `json:"page_size"`
		Sort     string `json:"sort"` //默认日期排序，也可按浏览，点赞，评论数排序
		Keyword  string `json:"keywords"`
	}
*/
type PageQuery struct {
	Page int `form:"page" binding:"required;gte=1"`
}

type ArticleKeyWord struct {
	Keyword string `form:"keyword" binding:"required"`
}
