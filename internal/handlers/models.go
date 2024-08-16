package handlers

type healthzGetResponse struct {
	ServiceName string `json:"serviceName,omitempty"`
	Version     string `json:"version,omitempty"`
	Status      string `json:"status,omitempty"`
}

type response struct {
	Message string
}
