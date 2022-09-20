package commands

import (
	"strings"
	"testing"
)

func TestQueryProfessor(t *testing.T) {
	// Test 1: Should not find any value.
	prof, err := QueryProfessor("Testing Not Found")

	if err == nil {
		t.Errorf("Got %v, should have received error / nil.", prof)
	}

	// Test 2: should find a value.
	prof, err = QueryProfessor("David")

	if err != nil {
		t.Errorf("Got %v as an error, should have received a professor.", err)
	}

	// Test 3: Empty string.
	prof, err = QueryProfessor("")

	if err == nil {
		t.Errorf("Got %v, should have received an error / nil.", prof)
	}
}

func TestCompareOverallRating(t *testing.T) {
	betterProf := RMPResult{
		professorName: "Hannibal of Carthage",
		totalRating:   "4.2",
		numRatings:    "100",
		courses:       []string{"COM101", "ULI101", "IPC144"},
		totalRatingByCourse: map[string]float64{
			"COM101": 4.2,
			"IPC144": 2.4,
			"ULI101": 4.5,
		},
		numRatingsByCourse: map[string]int{
			"COM101": 25,
			"IPC144": 12,
			"ULI101": 34,
		},
		wouldTakeAgain:    "80",
		levelOfDifficulty: "3",
		ratingDistribution: map[int]int{
			1: 10,
			2: 5,
			3: 15,
			4: 20,
			5: 25,
		},
		topTags: []string{"Lovely Human", "The Best Around", "Easy Grader"},
		rmpURL:  "www.google.ca",
	}

	worseProf := RMPResult{
		professorName: "Chaddius Maximus",
		totalRating:   "2.4",
		numRatings:    "20",
		courses:       []string{"OOP244", "ULI101", "IPC144"},
		totalRatingByCourse: map[string]float64{
			"OOP244": 2.5,
			"IPC144": 2.8,
			"ULI101": 2.1,
		},
		numRatingsByCourse: map[string]int{
			"OOP244": 10,
			"IPC144": 12,
			"ULI101": 14,
		},
		wouldTakeAgain:    "80",
		levelOfDifficulty: "3",
		ratingDistribution: map[int]int{
			1: 15,
			2: 5,
			3: 5,
			4: 10,
			5: 5,
		},
		topTags: []string{"Lovely Human", "The Best Around", "Easy Grader"},
		rmpURL:  "www.google.ca",
	}
	result := CompareOverallRating(betterProf, worseProf)
	if !strings.Contains(result, betterProf.professorName) {
		t.Errorf("Error: Better professor should have be chosen.")
	}
}
