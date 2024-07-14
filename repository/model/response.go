package model

type BaseResponse struct {
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

type BaseModel struct {
	Items        interface{} `json:"items"`
	TotalItem    int         `json:"total_item"`
	TotalPage    int         `json:"total_page"`
	FilteredItem int         `json:"filtered_item"`
	FilteredPage int         `json:"filtered_page"`
}

type BaseRequest struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}
