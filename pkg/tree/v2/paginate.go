package tree

import (
	"math"
)

func Paginate(items []string, page, pageSize int) []string {
	if page < 1 || pageSize < 1 {
		return []string{}
	}

	totalItems := len(items)
	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))

	if page > totalPages {
		return []string{}
	}

	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize

	if endIndex > totalItems {
		endIndex = totalItems
	}

	return items[startIndex:endIndex]
}
