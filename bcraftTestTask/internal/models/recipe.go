package models

type Recipe struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Info        string `json:"info"`
	Ingredients string `json:"ingredients"`
	Steps       []Step `json:"steps"`
}

type Step struct {
	Id   int64  `json:"id"`
	Info string `json:"info"`
	Time int64  `json:"time"`
}
