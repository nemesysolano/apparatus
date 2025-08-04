package service

type ApparatusQuestion struct {
	Prompt string `json:"prompt"`
	UserID string `json:"user-id"`
}

type ApparatusAnswer struct {
	Answer string  `json:"answer"`
	Score  float32 `json:"score"`
}
