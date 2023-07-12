package page

import (
	"github.com/pilagod/gorm-cursor-paginator/v2/paginator"
	proto "good/api/proto/pb"
)

// Request Request
type Request struct {
	NextPageToken *string `json:"next_page_token"`
	LastPageToken *string `json:"last_page_token"`
	Limit         *int    `json:"limit"`
	SortBy        []Sort  `json:"sort_by"`
}

// Sort Sort
type Sort struct {
	Key  string `json:"key"`
	Sort string `json:"sort"`
}

// ToPB ToPB
func (p *Request) ToPB() *proto.PageReq {
	res := &proto.PageReq{}
	if p.Limit != nil {
		l := int64(*p.Limit)
		res.Limit = &l
	}
	res.LastPageToken = p.LastPageToken
	res.NextPageToken = p.NextPageToken
	if p.SortBy != nil {
		res.SortBy = make([]*proto.Sort, 0)
		for _, sort := range p.SortBy {
			res.SortBy = append(res.SortBy, &proto.Sort{
				Key:  sort.Key,
				Sort: sort.Sort,
			})
		}
	}
	return res
}



// NewPBRequest NewPBRequest
func NewPBRequest(req *proto.PageReq) *Request {
	p := &Request{}
	if req != nil {
		p.NextPageToken = req.NextPageToken
		p.LastPageToken = req.LastPageToken
		if req.SortBy != nil {
			p.SortBy = make([]Sort, 0)
			for _, sort := range req.SortBy {
				p.SortBy = append(p.SortBy, Sort{
					Key:  sort.Key,
					Sort: sort.Sort,
				})
			}
		}

		if req.Limit != nil {
			l := int(*req.Limit)
			p.Limit = &l
		}
	}
	return p
}

// Result Result
type Result struct {
	NextPageToken *string `json:"next_page_token"`
	LastPageToken *string `json:"last_page_token"`
}

// CursorToResult CursorToResult
func CursorToResult(c paginator.Cursor) *Result {
	return &Result{
		NextPageToken: c.After,
		LastPageToken: c.Before,
	}
}

// ToPB ToPB
func (p *Result) ToPB() *proto.PageRes {
	res := &proto.PageRes{}
	if p.NextPageToken == nil {
		res.NextPageToken = ""
	} else {
		res.NextPageToken = *p.NextPageToken
	}

	if p.LastPageToken == nil {
		res.LastPageToken = ""
	} else {
		res.LastPageToken = *p.LastPageToken
	}

	return res
}

// NewPaginator NewPaginator
func NewPaginator(query *Request) *paginator.Paginator {
	opts := []paginator.Option{
		&paginator.Config{
			Keys:  []string{"ID"},
			Limit: 10,
			Order: paginator.ASC,
		},
	}

	if query.SortBy != nil {
		rules := make([]paginator.Rule, 0)
		for _, sort := range query.SortBy {
			order := paginator.DESC
			if sort.Sort == "ASC" || sort.Sort == "asc" {
				order = paginator.ASC
			}
			rules = append(rules, paginator.Rule{
				Key:   sort.Key,
				Order: order,
			})
		}

		opts = append(opts, paginator.WithRules(rules...))
	}
	if query.Limit != nil {
		opts = append(opts, paginator.WithLimit(*query.Limit))
	}
	if query.NextPageToken != nil {
		opts = append(opts, paginator.WithAfter(*query.NextPageToken))
	}
	if query.LastPageToken != nil {
		opts = append(opts, paginator.WithBefore(*query.LastPageToken))
	}
	return paginator.New(opts...)
}

