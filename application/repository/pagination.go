package repository

import (
	"context"
	"net/url"
	"strconv"

	"gorm.io/gorm"
)

func Paginate(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		if ctx == nil {
			return db.Offset(0).Limit(10)
		}

		q, ok := ctx.Value("values").(url.Values)
		if !ok {
			return db.Offset(0).Limit(10)
		}

		page, _ := strconv.Atoi(q.Get("page"))
		if page == 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(q.Get("page_size"))
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
