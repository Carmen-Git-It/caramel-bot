package commands

import "testing"

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

func generateTestProfessor() RMPResult {
	var testProfessor RMPResult

	// TODO: Put logic here!
	return testProfessor
}
