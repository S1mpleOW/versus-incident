package dto

type IncidentBody struct {
	Logs string `json:"logs"`
}

type BitbucketBuildPayload struct {
	Type        string `json:"type"`
	Key         string `json:"key"`
	State       string `json:"state"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
}
