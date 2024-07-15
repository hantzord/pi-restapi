package base

import "capstone/entities"

type BaseMetadataFullSuccessResponse struct {
	Status   bool                   `json:"status"`
	Message  string                 `json:"message"`
	Metadata *entities.MetadataFull `json:"metadata"`
	Data     any                    `json:"data"`
}

func NewMetadataFullSuccessResponse(message string, metadata *entities.MetadataFull, data any) *BaseMetadataFullSuccessResponse {
	return &BaseMetadataFullSuccessResponse{
		Status:   true,
		Message:  message,
		Metadata: metadata,
		Data:     data,
	}
}