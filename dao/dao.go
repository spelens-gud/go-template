package dao

type PageList struct {
	// 当前页码
	Page int `form:"page" json:"page"`
	// 每页数量
	PageSize int `form:"page_size" json:"page_size"`
	// 总数量
	Total int64 `form:"total" json:"total"`
}
