package views

type HealthResponse struct {
	LumiServerResponse Health `json:"lumi_server"`
	LumiDBResponse     Health `json:"lumi_db"`
	LumiCacheResponse  Health `json:"lumi_cache"`
	LumiAuthResponse   Health `json:"lumi_auth"`
}

type Health struct {
	IsUp    bool   `json:"is_up"`
	Message string `json:"message"`
}
