package sales

import (
	"fmt"
	"strings"
)

// EvidencePassage is an exact excerpt returned by the approved-playbook
// search. The ID is the stable reference used by script sections.
type EvidencePassage struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// ScriptText is a piece of draft text and the playbook passages that support
// it. EvidenceIDs is mandatory for claims and objection responses.
type ScriptText struct {
	Text        string   `json:"text"`
	EvidenceIDs []string `json:"evidenceIds,omitempty"`
}

// ObjectionResponse pairs a customer concern with a response grounded in the
// approved playbook.
type ObjectionResponse struct {
	Objection string     `json:"objection"`
	Response  ScriptText `json:"response"`
}

// SalesScript is a reviewable draft, not an instruction to place a call or
// make an offer.
type SalesScript struct {
	Opening        ScriptText          `json:"opening"`
	Questions      []string            `json:"questions"`
	Objections     []ObjectionResponse `json:"objections"`
	Claims         []ScriptText        `json:"claims"`
	StopConditions []string            `json:"stopConditions"`
}

// ValidateStructure checks required sections and resolves evidence IDs against
// the supplied passages. It does not establish that a passage semantically
// supports its text; callers need a separate support gate before presenting a
// draft as grounded.
func (script SalesScript) ValidateStructure(evidence []EvidencePassage) error {
	passageIDs := make(map[string]struct{}, len(evidence))
	for _, passage := range evidence {
		if strings.TrimSpace(passage.ID) == "" || strings.TrimSpace(passage.Text) == "" {
			return ValidationError{Field: "playbook passages", Problem: "must include an ID and text"}
		}
		if _, exists := passageIDs[passage.ID]; exists {
			return ValidationError{Problem: "duplicate playbook passage ID"}
		}
		passageIDs[passage.ID] = struct{}{}
	}

	if strings.TrimSpace(script.Opening.Text) == "" {
		return ValidationError{Field: "opening text", Problem: "is required"}
	}
	if len(script.Questions) == 0 {
		return ValidationError{Problem: "at least one question is required"}
	}
	for index, question := range script.Questions {
		if strings.TrimSpace(question) == "" {
			return ValidationError{Problem: fmt.Sprintf("question %d is required", index+1)}
		}
	}

	for index, claim := range script.Claims {
		label := fmt.Sprintf("claim %d", index+1)
		if strings.TrimSpace(claim.Text) == "" {
			return ValidationError{Field: label, Problem: "text is required"}
		}
		if err := validateEvidenceIDs(label, claim.EvidenceIDs, passageIDs); err != nil {
			return err
		}
	}

	for index, objection := range script.Objections {
		label := fmt.Sprintf("objection response %d", index+1)
		if strings.TrimSpace(objection.Objection) == "" || strings.TrimSpace(objection.Response.Text) == "" {
			return ValidationError{Field: label, Problem: "must include an objection and response"}
		}
		if err := validateEvidenceIDs(label, objection.Response.EvidenceIDs, passageIDs); err != nil {
			return err
		}
	}

	if len(script.StopConditions) == 0 {
		return ValidationError{Problem: "at least one stop condition is required"}
	}
	for index, stopCondition := range script.StopConditions {
		if strings.TrimSpace(stopCondition) == "" {
			return ValidationError{Problem: fmt.Sprintf("stop condition %d is required", index+1)}
		}
	}

	return nil
}

func validateEvidenceIDs(label string, ids []string, known map[string]struct{}) error {
	if len(ids) == 0 {
		return ValidationError{Field: label, Problem: "must cite at least one playbook passage"}
	}
	for _, id := range ids {
		if strings.TrimSpace(id) == "" {
			return ValidationError{Field: label, Problem: "contains an empty playbook passage ID"}
		}
		if _, exists := known[id]; !exists {
			return ValidationError{Field: label, Problem: "cites an unknown playbook passage"}
		}
	}
	return nil
}
