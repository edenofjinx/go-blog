package dbrepo

import (
	"bitbucket.org/julius_liaudanskis/go-blog/config"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// paginate pagination from url params
func paginate(r *http.Request, config *config.AppConfig) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		params := r.URL.Query()
		page := getPage(params, config)
		limit := getLimit(params, config)
		order := getOrder(params)

		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit).Order(fmt.Sprintf("id %s", order))
	}
}

// getPage gets page number from url params
func getPage(params url.Values, config *config.AppConfig) int {
	page := params.Get("page")
	if page == "" {
		page = "1"
	}
	p, err := strconv.Atoi(page)
	if err != nil {
		config.ErrorLog.Println(err)
	}
	if p <= 0 {
		p = 1
	}
	return p
}

// getLimit gets limit from url params
func getLimit(params url.Values, config *config.AppConfig) int {
	limit := params.Get("limit")
	if limit == "" {
		limit = "10"
	}
	l, err := strconv.Atoi(limit)
	if err != nil {
		config.ErrorLog.Println(err)
	}
	switch {
	case l > 100:
		l = 100
	case l <= 0:
		l = 10
	}
	return l
}

// getOrder gets order from url params
func getOrder(params url.Values) string {
	order := params.Get("order")
	o := strings.ToUpper(order)
	if o == "ASC" || o == "DESC" {
		return o
	}
	return "ASC"
}
