package where

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	defaultLimit = -1
)

type Tenant struct {
	Key       string
	ValueFunc func(ctx context.Context) string
}

var registeredTenant Tenant

type Where interface {
	where(db *gorm.DB) *gorm.DB
}

type Query struct {
	Query interface{}

	Args []interface{}
}

type Option func(*Options)

type Options struct {
	Offset int `json:"offset"`

	Limit int `json:"limit"`

	Filter map[any]any

	Clauses []clause.Expression

	Queries []Query
}

func WithOffset(offset int64) Option {
	return func(whr *Options) {
		if offset < 0 {
			offset = 0
		}
		whr.Offset = int(offset)
	}
}

func WithLimit(limit int64) Option {
	return func(whr *Options) {
		if limit <= 0 {
			limit = defaultLimit
		}
		whr.Limit = int(limit)
	}
}

func WithPage(page int, pageSize int) Option {
	return func(whr *Options) {
		if page == 0 {
			page = 1
		}
		if pageSize == 0 {
			pageSize = defaultLimit
		}
		whr.Offset = (page - 1) * pageSize
		whr.Limit = pageSize
	}
}

func WithFilter(filter map[any]any) Option {
	return func(whr *Options) {
		whr.Filter = filter
	}
}

func WithClauses(conds ...clause.Expression) Option {
	return func(whr *Options) {
		whr.Clauses = append(whr.Clauses, conds...)
	}
}

func WithQuery(query interface{}, args ...interface{}) Option {
	return func(whr *Options) {
		whr.Queries = append(whr.Queries, Query{Query: query, Args: args})
	}
}

func NewWhere(opts ...Option) *Options {
	whr := &Options{
		Offset:  0,
		Limit:   defaultLimit,
		Filter:  map[any]any{},
		Clauses: make([]clause.Expression, 0),
	}

	for _, opt := range opts {
		opt(whr)
	}
	return whr
}

func (whr *Options) O(offset int) *Options {
	if offset <= 0 {
		offset = 0
	}
	whr.Offset = offset
	return whr
}

func (whr *Options) L(limit int) *Options {
	if limit <= 0 {
		limit = defaultLimit
	}
	whr.Limit = limit
	return whr
}

func (whr *Options) P(page int, pageSize int) *Options {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = defaultLimit
	}

	whr.Offset = (page - 1) * pageSize
	whr.Limit = pageSize
	return whr
}

func (whr *Options) C(conds ...clause.Expression) *Options {
	whr.Clauses = append(whr.Clauses, conds...)
	return whr
}

func (whr *Options) Q(query interface{}, args ...interface{}) *Options {
	whr.Queries = append(whr.Queries, Query{Query: query, Args: args})
	return whr
}

func (whr *Options) T(ctx context.Context) *Options {
	if registeredTenant.Key != "" && registeredTenant.ValueFunc != nil {
		whr.F(registeredTenant.Key, registeredTenant.ValueFunc(ctx))
	}
	return whr
}

func (whr *Options) F(kvs ...any) *Options {
	if len(kvs)%2 != 0 {
		return whr
	}

	for i := 0; i < len(kvs); i = i + 2 {
		key := kvs[i]
		value := kvs[i+1]
		whr.Filter[key] = value
	}
	return whr
}

func (whr *Options) Where(db *gorm.DB) *gorm.DB {
	for _, query := range whr.Queries {
		conds := db.Statement.BuildCondition(query.Query, query.Args)
		whr.Clauses = append(whr.Clauses, conds...)
	}

	return db.Where(whr.Filter).Clauses(whr.Clauses...).Offset(whr.Offset).Limit(whr.Limit)
}

func O(offset int) *Options {
	return NewWhere().O(offset)
}

func L(limit int) *Options {
	return NewWhere().L(limit)
}

func P(page int, pageSize int) *Options {
	return NewWhere().P(page, pageSize)
}

func C(conds ...clause.Expression) *Options {
	return NewWhere().C(conds...)
}

func T(ctx context.Context) *Options {
	return NewWhere().F(registeredTenant.Key, registeredTenant.ValueFunc)
}

func F(kvs ...any) *Options {
	return NewWhere().F(kvs...)
}

func RegisterTenant(key string, valueFunc func(ctx context.Context) string) {
	registeredTenant = Tenant{
		Key:       key,
		ValueFunc: valueFunc,
	}
}
