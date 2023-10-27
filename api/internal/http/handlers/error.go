package handlers

type HttpError struct {
	Error     string `json:"error"`
	Code      int    `json:"code"`
	Temporary bool   `json:"temporary"`
}
