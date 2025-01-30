package requestresponse

type SwarmCreateWidgetRequestBody struct {
	ParentEntityID   string      `json:"parent_entity_id,omitempty"`
	ParentEntityType string      `json:"parent_entity_type,omitempty"`
	Metadata         interface{} `json:"metadata,omitempty"`
}

type SwarmRequestHeaders struct {
	MemberID     string `json:"x-member-id,omitempty"`
	APIKey       string `json:"x-api-key,omitempty"`
	PlatformType string `json:"x-platform-type,omitempty"`
}

type SwarmCreateWidgetRequest struct {
	Headers     SwarmRequestHeaders
	RequestBody SwarmCreateWidgetRequestBody
}
