package model

type ManageRequest struct {
	Username string `json:"username"`
	Action   string `json:"action"`
}

func NewManageRequest() *ManageRequest {
	return &ManageRequest{}
}
