package pagination

import (
	"net/http"
	"strconv"
)

// GetPagination get pagination from request.
func GetPagination(r *http.Request) (page int, limit int, ok bool) {
	var err error
	limitStr := r.URL.Query().Get("limit")
	pageStr := r.URL.Query().Get("page")

	if limitStr == "" && pageStr == "" {
		return
	}
	ok = true

	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	page, err = strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	return
}
