package response

type PaginationResponse struct {
	Limit   int  `json:"limit" example:"20"`
	Offset  int  `json:"offset" example:"0"`
	Total   *int `json:"total,omitempty" example:"150"`
	HasMore bool `json:"has_more" example:"true"`
}

func NewPaginationResponse(limit, offset int, total *int) PaginationResponse {
	pagination := PaginationResponse{
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}

	if total != nil {
		pagination.HasMore = offset+limit < *total
	}

	return pagination
}
