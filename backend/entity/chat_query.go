package entity

type ChatQueryRequest struct {
	SessionID string `json:"session_id"`
	Query     string `json:"query"`
}
