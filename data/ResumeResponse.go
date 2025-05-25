package data

type Resume struct {
	Resume1 string `json:"resume1"`
	Resume2 string `json:"resume2"`
}

type ResumeResponse struct {
	Resume Resume `json:"resume"`
}
