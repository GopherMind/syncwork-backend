package models

type Proposal struct {
	TaskID      string `json:"task_id"`
	UserID      string `json:"user_id"`
	CoverLetter string `json:"cover_letter"`
	Status      string `json:"status"`
}
