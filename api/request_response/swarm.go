package requestresponse

import "likeminds-pandemonium/api/constant"

type SwarmCreateWidgetRequestBody struct {
	ParentEntityID   string      `json:"parent_entity_id,omitempty"`
	ParentEntityType string      `json:"parent_entity_type,omitempty"`
	Metadata         interface{} `json:"metadata,omitempty"`
}

type SwarmCreateWidgetRequest struct {
	Headers     constant.ApiHeaders
	RequestBody SwarmCreateWidgetRequestBody
}
