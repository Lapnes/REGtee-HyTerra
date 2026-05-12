package pagination

import (
	"log"
	"math"
)

type Pagination struct {
	Limit    int    `form:"limit" cache:"optional"`
	Page     int    `form:"page" cache:"optional"`
	OrderBy  string `form:"order_by" cache:"optional"`
	OrderDir string `form:"order_dir" cache:"optional" validate:"omitempty,oneof=asc desc"`

	Offset    int
	TotalPage int
	TotalData int
}

const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaximumLimit = 100
)

func (p *Pagination) Paginate() {
	p.ValidatePagination()
	p.Offset = p.Limit * (p.Page - 1)
}

func (p *Pagination) SetToDefault() {
	log.Println(p)
	p.Page, p.Limit = DefaultPage, DefaultLimit
}

func (p *Pagination) ValidatePagination() {
	if p.Page < 1 || p.Limit < 1 {
		p.SetToDefault()
	}
	if p.Limit > MaximumLimit {
		p.Limit = MaximumLimit
	}
}

func (p *Pagination) SetTotalPage() {
	p.ValidatePagination()
	if p.TotalData > 0 && p.TotalData < p.Limit {
		p.Limit = p.TotalData
	}

	p.TotalPage = int(math.Ceil(float64(p.TotalData) / float64(p.Limit)))
}
