package data

type Resume struct {
	Result string `json:"result"`
	Source string `json:"source"`
}

type Resumes struct {
	Resume1 Resume `json:"resume1"`
	Resume2 Resume `json:"resume2"`
}

type ResumeResponse struct {
	Resume Resumes `json:"resume"`
}
