package historyRequest

type History struct {
}

type ListHistory struct {
	UserID    int64 `json:"-"`
	Code      string `json:"code,omitempty" form:"code"`
	Status    string `json:"status,omitempty" form:"status"`
	Page      int64  `json:"-" form:"page"`
	Limit     int64  `json:"-" form:"limit"`
	OrderBy   string `json:"-" form:"order_by"`
	OrderDesc bool   `json:"-" form:"order_desc"`
}
