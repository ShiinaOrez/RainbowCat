package crad

import (
	"context"
	"github.com/ShiinaOrez/RainbowCat/crawlers/crad/sql"
	_const "github.com/ShiinaOrez/kylin/const"
	"github.com/anaskhan96/soup"
	"github.com/mozillazg/request"
	"net/http"
	"strconv"
	"strings"
)

type Result struct {
	Articles []Article `json:"articles"`
	Total    string    `json:"total"`
}

type Article struct {
	ID       string   `json:"id"`
	Link     string   `json:"link"`
	Topic    string   `json:"topic"`
	Authors  string   `json:"authors"`
	DOI      string   `json:"doi"`
	Abstract string   `json:"abstract"`
}

func checkFailed(err error, notifyCh *chan int) bool {
	if err != nil {
		*notifyCh<- _const.StatusFailed
	}
	return err != nil
}

func Proc(ctx context.Context, notifyCh *chan int) {
	var (
		kw   string
		page int
	)
	kw = ctx.Value("keyword").(string)
	page, _ = strconv.Atoi(ctx.Value("page").(string))
	c := new(http.Client)
	req := request.NewRequest(c)
	searchSQL := sql.BuildAllSearchSQL(kw).DefaultSource().Page(page)
	data := strings.NewReader(searchSQL.SQL())
	resp, err := req.PostForm("http://crad.ict.ac.cn/CN/article/advancedSearchResult.do", data)
	defer resp.Body.Close()
	if checkFailed(err, notifyCh) {
		return
	}
	text, err := resp.Text()
	if checkFailed(err, notifyCh) {
		return
	}
	doc := soup.HTMLParse(text)
	result := Result{}
	divs := doc.Find("form", "id", "AbstractList").
		FindAll("div", "class", "noselectrow")
	result.Articles = make([]Article, len(divs))
	for index, div := range divs {
		result.Articles[index].ID = div.Attrs()["id"][3:]
		tr := div.Find("table").
			Find("tbody").
			Find("tr")

		link := tr.Find("a", "class", "txt_biaoti")
		authorList := tr.Find("div", "class", "authorList").Find("span")
		doi := tr.Find("span", "class", "abs_njq").Find("a")
		abstract := tr.Find("div", "id", "Abstract"+result.Articles[index].ID)
		result.Articles[index].Link = link.Attrs()["href"]
		result.Articles[index].Topic = link.Text()
		result.Articles[index].Authors = authorList.Text()
		if doi.Pointer != nil {
			result.Articles[index].DOI = doi.Attrs()["href"]
		}
		if abstract.Pointer != nil {
			result.Articles[index].Abstract = abstract.FullText()
		}
	}

	var pageTotal string
	if len(divs) < 30 {
		pageTotal = "共"+strconv.Itoa(len(divs))+"条记录"
	} else {
		ul := doc.FindAll("ul", "class", "page_ul_two")
		for _, u := range ul {
			lis := u.FindAll("li")
			for _, li := range lis {
				if v, ok := li.Attrs()["class"]; ok {
					if v == "no_background page-total" {
						pageTotal = li.Text()
					}
				}
			}
		}
	}

	result.Total = pageTotal
	*notifyCh<- _const.StatusSuccess
	return
}