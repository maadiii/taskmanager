package utils

import (
	"fmt"
	"strings"
)

type Pagination struct {
	lastID string
	limit  int
	order  string
	query  *strings.Builder

	args []any
}

func (p *Pagination) LastID(id string) *Pagination {
	p.lastID = id

	return p
}

func (p *Pagination) Limit(limit int) *Pagination {
	p.limit = limit

	return p
}

func (p *Pagination) Desc() *Pagination {
	p.order = "DESC"

	return p
}

func (p *Pagination) ASC() *Pagination {
	p.order = "ASC"

	return p
}

func (p *Pagination) Select(tableName, columnName string, otherColumns ...string) *Pagination {
	builder := new(strings.Builder)
	builder.WriteString("SELECT ")
	fmt.Fprintf(builder, " %s ", columnName)
	if len(otherColumns) > 0 {
		builder.WriteString(strings.Join(otherColumns, ", "))
	}
	builder.WriteString(" FROM ")
	builder.WriteString(tableName)

	p.query = builder

	return p
}

func (p *Pagination) Paginate(query string, args ...any) (string, []any) {
	if query != "" {
		fmt.Fprintf(p.query, " WHERE %s", query)
		p.args = append(p.args, args...)
	}

	if p.lastID != "" {
		fmt.Fprintf(p.query, " AND id > $%d", len(p.args)+1)
		p.args = append(p.args, p.lastID)
	}

	if p.order != "" {
		fmt.Fprintf(p.query, " ORDER BY id $%d", len(p.args)+1)
		p.args = append(p.args, p.order)
	}

	if p.limit > 0 {
		fmt.Fprintf(p.query, " LIMIT $%d", len(p.args)+1)
		p.args = append(p.args, p.limit)
	}

	return p.query.String(), p.args
}
