package sql

import "fmt"

type sql struct {
	text   string
}

func (s *sql) first(key, value string) *sql {
	s.text = fmt.Sprintf("(%s[%s])", value, key)
	return s
}

func (s *sql) and(key, value string) *sql {
	expr := fmt.Sprintf("(%s[%s])", value, key)
	s.text = fmt.Sprintf("(%s AND %s)", s.text, expr)
	return s
}

func (s *sql) or(key, value string) *sql {
	expr := fmt.Sprintf("(%s[%s])", value, key)
	s.text = fmt.Sprintf("(%s OR %s)", s.text, expr)
	return s
}


