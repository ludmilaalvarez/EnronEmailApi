package models

type EmailResponse struct {
	Hits HitsInfo `json:"hits"`
}

type HitsInfo struct {
	Hits []EmailHit `json:"hits"`
}

type EmailHit struct {
	Source Email `json:"_source"`
}
