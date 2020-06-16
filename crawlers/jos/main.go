package main

import (
	"encoding/json"
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/mozillazg/request"
	"net/http"
	"strings"
)

const (
	Keyword = "key_word"
	EnKeyword = "en_key_word"
	Title = "title"
	EnTitle = "en_title"
	Username = "user_real_name"
	EnUsername = "pin_yin_name"
	Institution = "cn_institution"
	EnInstitution = "en_institution"
	Abstract = "abstract"
	EnAbstract = "en_abstract"
	FundProject = "fund_project"
	DOI = "doi"
	ColumnName = "column_name"
	FileNo = "file_no"


	EVENTTARGET = "__EVENTTARGET"
	EVENTARGUMENT = "__EVENTARGUMENT"
	VIEWSTATE = "__VIEWSTATE"
	VIEWSTATEGENERATOR = "__VIEWSTATEGENERATOR"
	VIEWSTATEENCRYPTED = "__VIEWSTATEENCRYPTED"
)

type Result struct {
	Articles []Article `json:"articles"`
	Total    string    `json:"total"`
}

type Article struct {
	No        string   `json:"no"`
	Title     string   `json:"title"`
	Link      string   `json:"link"`
	Authors   string   `json:"authors"`
	Year      string   `json:"year"`
	Volume    string   `json:"volume"`
	Period    string   `json:"period"`
	PageDuring string  `json:"page_during"`
}

func main() {
	c := new(http.Client)
    var (
		formMap = map[string]string {
			"KeyList": Keyword,
			"Key": "协同工作",
			"StartYearList": "0",
			"EndYearList": "2020",
			"to": "",
		}
	)
    text, _ := soup.Get("http://www.jos.org.cn/jos/ch/reader/key_query.aspx")
	doc := soup.HTMLParse(text)
	inputs := doc.Find("body").FindAll("input")
	for _, input := range inputs {
		if value, ok := input.Attrs()["value"]; !ok {
			continue
		} else {
			if input.Attrs()["name"] != "go"{
				formMap[input.Attrs()["name"]] = value
			}
		}
	}
	req := request.NewRequest(c)
	req.Headers["Content-Type"] = "application/x-www-form-urlencoded"
	resp, _ := req.PostForm("http://www.jos.org.cn/jos/ch/reader/key_query.aspx", formMap)
	defer resp.Body.Close()

	newText, _ := resp.Text()
	newDoc := soup.HTMLParse(newText)
	result := Result{}
	trs := newDoc.Find("table", "id", "DataGrid1").FindAll("tr")
	for _, tr := range trs {
		if _, ok := tr.Attrs()["align"]; ok {
			continue
		}
		article := Article{}
		tds := tr.FindAll("td")
		if len(tds) < 5 {
			continue
		}
		article.No = trimStr(tds[0].FullText())
		article.Title = tds[1].Find("a").Text()
		article.Authors = trimStr(tds[2].FullText())
		yearsInfoString := trimStr(tds[3].FullText())
		article.Year = strings.Split(yearsInfoString, ",")[0]
		yearsInfoString = yearsInfoString[len(article.Year)+1:]
		article.Volume = strings.Split(yearsInfoString, "(")[0]
		yearsInfoString = yearsInfoString[len(article.Volume)+1:]
		article.Period = strings.Split(yearsInfoString, ")")[0]
		yearsInfoString = yearsInfoString[len(article.Period)+2:]
		article.PageDuring = yearsInfoString
		article.Link = fmt.Sprintf("http://www.jos.org.cn/jos/ch/reader/view_abstract.aspx?flag=%d&file_no=%s&journal_id=%s", 1, article.No, "jos")
		result.Articles = append(result.Articles, article)
	}
	bytes, _ := json.MarshalIndent(result, "", "    ")
	str := strings.Replace(string(bytes), "\\u0026", "&", -1)
	fmt.Println(str)
}

func trimStr(str string) string {
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\t", "", -1)
	return str
}