package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	// "github.com/gocolly/colly"
)

type ProfessorQueryResult struct {
	Professors         []Professor
	SearchResultsTotal int
	Remaining          int
	Type               string
}

type Professor struct {
	TDept            string
	TSid             string
	Institution_name string
	TFname           string
	TMiddleName      string
	TLname           string
	TId              int
	TNumRatings      int
	Rating_class     string
	ContentType      string
	CategoryType     string
	Overall_Rating   string
}

type RatingQueryResult struct {
}

type Rating struct {
	Attendence        string
	ClarityColor      string
	EasyColor         string
	HelpColor         string
	HelpCount         int
	Id                int
	NotHelpCount      int
	OnlineClass       string
	Quality           string
	RClarity          int
	RComments         string
	RDate             string
	REasy             int
	REasyString       string
	RErrorMsg         string
	RHelpful          int
	RInterest         string
	ROverall          int
	ROverallString    string
	RStatus           int
	RTextBookUse      string
	RTimestamp        int
	RWouldTakeAgain   string
	SId               int
	TakenForCredit    string
	Teacher           string
	TeacherGrade      string
	TeacherRatingRags []string
	UnUsefulGrouping  string
	UsefulGrouping    string
}

func ScrapeProfessor(professor string) (Professor, error) {
	// c := colly.NewCollector(
	// // colly.AllowedDomains("ratemyprofessors.com"),
	// )

	// // DEBUG Read every div lol
	// c.OnHTML("div", func(e *colly.HTMLElement) {
	// 	fmt.Println(e.Text)
	// 	if strings.Contains(e.Text, "Loading...") {
	// 		fmt.Println("Loading")
	// 		c.Visit(baseURL + url.QueryEscape(professor))
	// 	}
	// })

	// c.OnHTML("a[href]", func(e *colly.HTMLElement) {

	// 	fmt.Println("Checking a with href: ", e.Attr("href"))
	// 	if strings.Contains(e.Attr("href"), "ShowRatings") {
	// 		fmt.Println("Found Professor Link to: ", e.Attr("href"))

	// 		c.Visit(e.Attr("href"))
	// 	}
	// })

	// c.OnHTML("div[class]", func(e *colly.HTMLElement) {
	// 	if strings.Contains(e.Attr("class"), "RatingValue__Numerator") {
	// 		fmt.Println("Found Rating Numerator: ", e.Text)
	// 	}
	// })

	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Scraping: ", r.URL)

	// })

	// c.Visit(baseURL + url.QueryEscape(professor))

	var result Professor = Professor{}

	var queryString string = fmt.Sprint("https://www.ratemyprofessors.com/filter/professor/?&page=1&filter=teacherlastname_sort_s+asc&query=", url.QueryEscape(professor), "&queryoption=TEACHER&queryBy=schoolId&sid=1497")

	fmt.Println("Checking query: ", queryString)

	response, err := http.Get(queryString)

	if err != nil {
		fmt.Println("Error contacting RateMyProfessors API server")
		fmt.Println(err)
		return result, errors.New("Error contacting RateMyProfessors API server")
	}

	defer response.Body.Close()

	var body ProfessorQueryResult

	err = json.NewDecoder(response.Body).Decode(&body)
	if err != nil {
		fmt.Println("Error decoding JSON")
		return result, errors.New("Error decoding JSON")
	}

	fmt.Println("Found ", body.SearchResultsTotal, " results")
	fmt.Println("First professor found: ", body.Professors[0].TFname, body.Professors[0].TMiddleName, body.Professors[0].TLname)

	return result, nil
}
