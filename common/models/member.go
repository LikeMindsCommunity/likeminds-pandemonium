package models

type Member struct {
	SDKClientInfo SDKClientInfo `json:"sdk_client_info"`
}

type SDKClientInfo struct {
	UUID string `json:"uuid"`
}
