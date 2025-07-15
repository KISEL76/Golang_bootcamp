package common

type SuccessResponse struct {
	Change int    `json:"change"`
	Thanks string `json:"thanks"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type CandyOrder struct {
	Money      int    `json:"money"`
	CandyType  string `json:"candyType"`
	CandyCount int    `json:"candyCount"`
}
