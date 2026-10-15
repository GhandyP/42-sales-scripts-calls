package sales

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	MaxCustomerFactsRunes    = 4000
	MaxObjectiveRunes        = 1000
	MaxApprovedPlaybookRunes = 20000
)

// CallBrief contains only the information supplied for one draft request.
// CallBrief values are transient and are never persisted by this package.
type CallBrief struct {
	CustomerFacts    string `json:"customerFacts"`
	Objective        string `json:"objective"`
	ApprovedPlaybook string `json:"approvedPlaybook"`
}

// ValidationError describes a safe, user-displayable validation failure. It
// contains field names and rule text, never the rejected input value.
type ValidationError struct {
	Field   string
	Problem string
}

func (err ValidationError) Error() string {
	if err.Field == "" {
		return err.Problem
	}
	return err.Field + " " + err.Problem
}

// Validate rejects incomplete or excessively large request fields without
// changing the caller's brief or including its contents in an error.
func (brief CallBrief) Validate() error {
	fields := []struct {
		name  string
		value string
		limit int
	}{
		{name: "customer facts", value: brief.CustomerFacts, limit: MaxCustomerFactsRunes},
		{name: "call objective", value: brief.Objective, limit: MaxObjectiveRunes},
		{name: "approved playbook text", value: brief.ApprovedPlaybook, limit: MaxApprovedPlaybookRunes},
	}

	for _, field := range fields {
		if strings.TrimSpace(field.value) == "" {
			problem := "is required"
			if field.name == "customer facts" {
				problem = "are required"
			}
			return ValidationError{Field: field.name, Problem: problem}
		}
		if utf8.RuneCountInString(field.value) > field.limit {
			return ValidationError{
				Field:   field.name,
				Problem: fmt.Sprintf("must be at most %d characters", field.limit),
			}
		}
	}
	return nil
}
