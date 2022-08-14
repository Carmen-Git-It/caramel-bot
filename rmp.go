package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
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
	totalRating         int
	numRatings          int
	courses             []string
	totalRatingByCourse map[string]int
	numRatingsByCourse  map[string]int
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

	fmt.Println("Found ", professorBody.SearchResultsTotal, " results")
	fmt.Println("First professor found: ", professorBody.Professors[0].TFname, professorBody.Professors[0].TMiddleName, professorBody.Professors[0].TLname)

	// Pagination variables
	var page int = 1
	var remaining = 1
	var ratingQuery = fmt.Sprint("https://www.ratemyprofessors.com/paginate/professors/ratings?tid=", professorBody.Professors[0].TId, "&filter=&courseCode=&page=", page)

	// Start building result
	result.professorName = professorBody.Professors[0].TFname + " " + professorBody.Professors[0].TMiddleName + " " + professorBody.Professors[0].TLname
	result.numRatingsByCourse = make(map[string]int)
	result.totalRatingByCourse = make(map[string]int)

	// To try to make up for RMP's artificial skew
	ratingCount := make(map[int]int)

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

		// Incorporate ratings into result
		for _, rating := range ratingBody.Ratings {
			// Ignore rating if not helpful AND a 1.0 rating (not sure why, seems fucked up tbh)
			if rating.HelpCount >= rating.NotHelpCount || rating.ROverall > 2.0 {
				ratingCount[int(math.Ceil(rating.ROverall))]++
				result.totalRating += int(rating.RClarity) + int(rating.RHelpful)
				result.numRatings += 2
				if result.numRatingsByCourse[rating.RClass] == 0 {
					result.courses = append(result.courses, rating.RClass)
				}
				result.totalRatingByCourse[rating.RClass] += int(rating.RClarity) + int(rating.RHelpful)
				result.numRatingsByCourse[rating.RClass] += 2
			}
		}
	}

	fmt.Println("Number of 5s: ", ratingCount[5], " Number of ratings total: ", (result.numRatings / 2))

	fmt.Println("Number of 5s: ", ratingCount[5], " Number of 4s: ", ratingCount[4], " Number of 3s: ", ratingCount[3], " Number of 2s: ", ratingCount[2], " Number of 1s: ", ratingCount[1])

	// If the professor has a ton of 5* ratings, remove all 1* ratings
	if ratingCount[5] > (result.numRatings/2)-ratingCount[5] {
		fmt.Println("Professor is a 5 star professor")
		result.totalRating -= ratingCount[1] * 2
		result.numRatings -= ratingCount[1] * 2
		result.totalRating -= ratingCount[2] * 2
		result.numRatings -= ratingCount[2] * 2
	}

	fmt.Println("Overall rating of ", result.professorName, ": ", float64(result.totalRating)/float64(result.numRatings), " from ", result.numRatings/2, " ratings")

	for _, course := range result.courses {
		fmt.Println("Rating for ", course, ": ", float64(result.totalRatingByCourse[course])/float64(result.numRatingsByCourse[course]), " from ", result.numRatingsByCourse[course]/2, " ratings")
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
