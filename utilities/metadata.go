package utilities

import (
	"capstone/entities"
	"strconv"
)

func GetMetadata(pageParam, limitParam string) *entities.Metadata {
	page, err := strconv.Atoi(pageParam)
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 10
	}
	return &entities.Metadata{
		Page:  page,
		Limit: limit,
	}
}

func GetFullMetadata(pageParam string, limitParam string, sortParam string, orderParam string, searchParam string) *entities.MetadataFull {
	page, err := strconv.Atoi(pageParam)
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 10
	}

	if sortParam == "" {
		sortParam = "id"
	}

	if orderParam == "" {
		orderParam = "asc"
	}

	if searchParam == "" {
		searchParam = ""
	}

	return &entities.MetadataFull{
		Page:   page,
		Limit:  limit,
		Sort:   sortParam,
		Order:  orderParam,
		Search: searchParam,
	}
}