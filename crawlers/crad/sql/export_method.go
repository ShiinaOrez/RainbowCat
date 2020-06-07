package sql

import "strconv"

func (s *sql) SQL() string {
	return "searchSQL="+s.text
}

func (s *sql) DefaultSource() *sql {
	return s.and("Journal", "1")
}

func (s *sql) Page(p int) *sql {
	return s.and("Pager", strconv.Itoa(p))
}

func BuildAllSearchSQL(keyword string) *sql {
	s := &sql{}
	return s.first("Title", keyword).
		or("Abstract", keyword).
		or("Keyword", keyword).
		or("Author", keyword).
		or("AuthorCompany", keyword).
		or("DOI", keyword)
}

func BuildKeywordSearchSQL(kw string) *sql {
	s := &sql{}
	return s.first("Keyword", kw)
}