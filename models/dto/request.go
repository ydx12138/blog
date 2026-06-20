package dto

/*
	type Articles struct {
		Page     int    `json:"page"`
		PageSize int    `json:"page_size"`
		Sort     string `json:"sort"` //默认日期排序，也可按浏览，点赞，评论数排序
		Keyword  string `json:"keywords"`
	}
*/
type PageQuery struct {
	Page int `form:"page" binding:"required,gte=1"`
}

type ArticleKeyWord struct {
	Keyword string `form:"keyword" binding:"required"`
}

type AdminLogin struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type UserRegister struct {
	Email      string `form:"email" json:"email" binding:"required,email"`
	Password   string `form:"password" json:"password" binding:"required,min=6,max=10"`
	Repassword string `form:"re_password" json:"re_password" binding:"required,min=6,max=10,eqfield=Password"`
	Nickname   string `form:"nickname" json:"nickname"`
}

type UserLogin struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required"`
}

type CreateCommentReq struct {
	ArticleID uint64 `form:"article_id" json:"article_id" binding:"required"`
	Content   string `form:"content" json:"content" binding:"required"`
	ParentID  uint64 `form:"parent_id" json:"parent_id"`
}

type CreateArticleReq struct {
	Title       string `json:"title" binding:"required"`
	Summary     string `json:"summary"`
	Content     string `json:"content"`
	ContentType int8   `json:"content_type"`
	Cover       string `json:"cover"`
	CategoryID  uint64 `json:"category_id"`
	Tags        string `json:"tags"`
	Status      int8   `json:"status"`
}

type UpdateArticleReq struct {
	ID          uint64 `json:"id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Summary     string `json:"summary"`
	Content     string `json:"content"`
	ContentType int8   `json:"content_type"`
	Cover       string `json:"cover"`
	CategoryID  uint64 `json:"category_id"`
	Tags        string `json:"tags"`
	Status      int8   `json:"status"`
}

type PageQueryWithSize struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

type CategoryArticlesQuery struct {
	CategoryID uint64 `form:"category_id" binding:"required"`
	Page       int    `form:"page"`
}

type CommentListQuery struct {
	ArticleID uint64 `form:"article_id" binding:"required"`
	Page      int    `form:"page"`
}

type AdminArticleQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Keyword  string `form:"keyword"`
	Status   int8   `form:"status"`
}

type AdminCommentQuery struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

type UserStatusReq struct {
	ID     uint64 `json:"id" binding:"required"`
	Status uint64 `json:"status" binding:"required"`
}

type IDReq struct {
	ID uint64 `form:"id" uri:"id" json:"id" binding:"required"`
}
