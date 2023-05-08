package pagination

import (
	"fmt"
	"strconv"

	"github.com/kondohiroki/go-boilerplate/internal/logger"
	"go.uber.org/zap"
)

type PaginationDTI struct {
	Page     string `json:"page" validate:"required"`
	PerPage  string `json:"perPage" validate:"required"`
	SortBy   string `json:"sortBy" validate:"required"`
	SortDesc string `json:"sortDesc" validate:"required"`
}

type PaginationDTO struct {
	Total   int  `json:"total"`
	Limit   int  `json:"limit"`
	Page    int  `json:"page"`
	HasMore bool `json:"has_more"`
}

func ConvertPaginationToStrSql(pag *PaginationDTI) (string, error) {
	resultSql := ""

	//set s
	if pag.SortBy != "" {
		resultSql += " ORDER BY " + pag.SortBy
		if pag.SortDesc == "true" {
			resultSql += " DESC"
		}
	}

	if pag.Page != "" && pag.PerPage != "" {
		limit, err := strconv.Atoi(pag.PerPage)
		if err != nil {
			logger.Log.Error("convert string to int error.", zap.Error(err))
			return "", err
		}

		page, err := strconv.Atoi(pag.Page)
		if err != nil {
			logger.Log.Error("convert string to int error.", zap.Error(err))
			return "", err
		}

		offset := (page - 1) * limit
		resultSql += fmt.Sprintf(` LIMIT %d OFFSET %d`, limit, offset)
	}

	return resultSql, nil
}

func GetResponsePagination(pagDTI *PaginationDTI, total int) (PaginationDTO, error) {

	var resPag PaginationDTO

	resPag.Total = total
	if pagDTI.Page != "" && pagDTI.PerPage != "" {
		page, err := strconv.Atoi(pagDTI.Page)
		if err != nil {
			logger.Log.Error("convert string to int error.", zap.Error(err))
			return resPag, err
		}
		perPage, err := strconv.Atoi(pagDTI.PerPage)
		if err != nil {
			logger.Log.Error("convert string to int error.", zap.Error(err))
			return resPag, err
		}
		resPag.Page = page
		resPag.Limit = perPage
		resPag.HasMore = isHasMore(page, perPage, total)
	}
	return resPag, nil
}

func isHasMore(page int, limit int, total int) bool {
	return total > (page * limit)
}
