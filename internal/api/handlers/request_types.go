package handlers

type ScrapedDataRequest struct {
	Title           string   `json:"title"`
	Headers         []Header `json:"headers"`
	Paragraphs      []string `json:"paragraphs"`
	Links           []Link   `json:"links"`
	MetaDescription *string  `json:"metaDescription"`
	MainContent     *string  `json:"mainContent"`
}

type Header struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Link struct {
	Text string `json:"text"`
	Href string `json:"href"`
}


type ValidateKeyRequest struct {
	APIKey string `json:"apiKey"`
}