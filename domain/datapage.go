package domain

type DataPage struct {
	Status      bool        `json:"status"`
	Message     string      `json:"message"`
	Total       int64       `json:"total"`
	Pages       int         `json:"pages"`
	CurrentPage uint        `json:"current_page"`
	Limit       uint        `json:"per_page"`
	Data        interface{} `json:"data"`
}
