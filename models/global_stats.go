package models

type GlobalSnippetsData struct {
	TotalSnippets int64 `json:"total_snippets"`
}

type GlobalSnippetsResponse struct {
	Success bool               `json:"success"`
	Data    GlobalSnippetsData `json:"data"`
}

type GlobalLanguagesData struct {
	UniqueLanguages []string `json:"unique_languages"`
	TotalLanguages  int64    `json:"total_languages"`
}

type GlobalLanguagesResponse struct {
	Success bool     `json:"success"`
	Data    []string `json:"data"`
}
