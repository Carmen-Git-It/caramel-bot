package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
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
	Ratings   []Rating
	Remaining int
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
	RClass            string
	RComments         string
	RDate             string
	REasy             float64
	REasyString       string
	RErrorMsg         string
	RHelpful          int
	RInterest         string
	ROverall          float64
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

type RMPResult struct {
	professorName       string
	totalRating         string
	numRatings          string
	courses             []string
	totalRatingByCourse map[string]float64
	numRatingsByCourse  map[string]int
	wouldTakeAgain      string
	levelOfDifficulty   string
	ratingDistribution  map[int]int
}

func QueryProfessor(professor string) (RMPResult, error) {
	var result RMPResult

	var professorQuery string = fmt.Sprint("https://www.ratemyprofessors.com/filter/professor/?&page=1&filter=teacherlastname_sort_s+asc&query=", url.QueryEscape(professor), "&queryoption=TEACHER&queryBy=schoolId&sid=1497")
	fmt.Println("Checking query: ", professorQuery)

	// Query the API for the professor name given
	professorResponse, err := http.Get(professorQuery)
	if err != nil {
		fmt.Println("Error contacting RateMyProfessors API server")
		fmt.Println(err)
		return result, errors.New("Error contacting RateMyProfessors API server")
	}
	defer professorResponse.Body.Close()

	var professorBody ProfessorQueryResult

	// Decode the JSON response into a struct
	err = json.NewDecoder(professorResponse.Body).Decode(&professorBody)
	if err != nil {
		fmt.Println("Error decoding JSON")
		return result, errors.New("Error decoding JSON")
	}

	// If no professors were found, return an error
	if professorBody.SearchResultsTotal == 0 {
		return result, errors.New("Professor not found")
	}

	fmt.Println("Found ", professorBody.SearchResultsTotal, " results")
	fmt.Println("First professor found: ", professorBody.Professors[0].TFname, professorBody.Professors[0].TMiddleName, professorBody.Professors[0].TLname)

	result.professorName = professorBody.Professors[0].TFname + " " + professorBody.Professors[0].TMiddleName + " " + professorBody.Professors[0].TLname

	var profId int = professorBody.Professors[0].TId
	var scrapeURL string = fmt.Sprint("https://www.ratemyprofessors.com/ShowRatings.jsp?tid=", profId)

	// Create web scraper
	c := colly.NewCollector(
		colly.AllowedDomains("www.ratemyprofessors.com"),
	)

	c.OnHTML("div[class]", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("class"), "RatingValue__Numerator") {
			result.totalRating = e.Text
		}
		if strings.Contains(e.Attr("class"), "FeedbackItem__FeedbackNumber") {
			if strings.Contains(e.Text, "%") {
				result.wouldTakeAgain = e.Text
			} else {
				result.levelOfDifficulty = e.Text
			}
		}
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if e.Attr("href") == "#ratingsList" {
			result.numRatings = strings.Split(e.Text, " ")[0]
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit(scrapeURL)

	// Pagination variables
	var page int = 1
	var remaining = 1
	var ratingQuery = fmt.Sprint("https://www.ratemyprofessors.com/paginate/professors/ratings?tid=", professorBody.Professors[0].TId, "&filter=&courseCode=&page=", page)

	// Start building result
	result.professorName = professorBody.Professors[0].TFname + " " + professorBody.Professors[0].TLname
	result.numRatingsByCourse = make(map[string]int)
	result.totalRatingByCourse = make(map[string]float64)

	// To try to make up for RMP's artificial skew
	result.ratingDistribution = make(map[int]int)

	// Loop until no ratings are left
	for ok := true; ok; ok = remaining > 0 {
		fmt.Println("Checking query: ", ratingQuery)
		ratingResponse, err := http.Get(ratingQuery)

		// Query the API for ratings
		if err != nil {
			fmt.Println("Error contacting RateMyProfessors API server")
			fmt.Println(err)
			return result, errors.New("Error contacting RateMyProfessors API server")
		}
		defer ratingResponse.Body.Close()

		// Decode the JSON response into a struct
		var ratingBody RatingQueryResult
		err = json.NewDecoder(ratingResponse.Body).Decode(&ratingBody)
		if err != nil {
			fmt.Println("Error decoding JSON")
			fmt.Println(err)
			return result, errors.New("Error decoding JSON")
		}

		// Check remaining ratings
		remaining = ratingBody.Remaining
		fmt.Println(ratingBody.Remaining, " results remaining")
		page++
		ratingQuery = fmt.Sprint("https://www.ratemyprofessors.com/paginate/professors/ratings?tid=", professorBody.Professors[0].TId, "&filter=&courseCode=&page=", page)

		if len(ratingBody.Ratings) > 0 {
			// Incorporate ratings into result
			for _, rating := range ratingBody.Ratings {
				// Ignore rating if not helpful AND a 1.0 rating (not sure why, seems fucked up tbh)
				result.ratingDistribution[int(math.Ceil(rating.ROverall))]++
				if result.numRatingsByCourse[rating.RClass] == 0 {
					result.courses = append(result.courses, rating.RClass)
				}
				result.totalRatingByCourse[rating.RClass] += float64(rating.RClarity) + float64(rating.RHelpful)
				result.numRatingsByCourse[rating.RClass] += 2
			}
		} else {
			return result, errors.New("No ratings found")
		}
	}

	for _, course := range result.courses {
		result.totalRatingByCourse[course] = result.totalRatingByCourse[course] / float64(result.numRatingsByCourse[course])
		result.numRatingsByCourse[course] /= 2
	}

	return result, nil
}

func stringSliceContains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
