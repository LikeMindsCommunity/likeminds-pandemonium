package requestresponse

import "github.com/LikeMindsCommunity/likeminds-pandemonium/api/constant"

type SwarmCreateWidgetRequestBody struct {
	ParentEntityID   string      `json:"parent_entity_id,omitempty"`
	ParentEntityType string      `json:"parent_entity_type,omitempty"`
	Metadata         interface{} `json:"metadata,omitempty"`
}

type SwarmCreateWidgetRequest struct {
	Headers     constant.ApiHeaders
	RequestBody SwarmCreateWidgetRequestBody
}

type SwarmWidgetResponse struct {
	ID               string                 `json:"_id"`
	CreatedByLM      bool                   `json:"created_by_lm"`
	ParentEntityID   string                 `json:"parent_entity_id"`
	ParentEntityType string                 `json:"parent_entity_type"`
	MetaData         map[string]interface{} `json:"metadata"`
	LMMeta           map[string]interface{} `json:"_lm_meta"`
	CommunityId      int                    `json:"community_id"`
	CreatedAt        *int64                 `json:"created_at"`
	UpdatedAt        *int64                 `json:"updated_at"`
}
