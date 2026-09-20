package structs

import resp "link-service/internal/lib/api/response"

type Request struct {
	URL   string `json:"url", validate:"required, url"`
	Alias string `json:"alias", omitempty`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}