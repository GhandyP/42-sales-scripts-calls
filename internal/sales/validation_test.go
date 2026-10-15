package sales

import (
	"reflect"
	"strings"
	"testing"
)

func TestCallBriefValidateRequiresAllInputs(t *testing.T) {
	valid := CallBrief{
		CustomerFacts:    "private-sentinel customer facts",
		Objective:        "prepare for a discovery call",
		ApprovedPlaybook: "private-sentinel approved playbook text",
	}

	tests := []struct {
		name  string
		brief CallBrief
		want  string
	}{
		{
			name: "customer facts are blank",
			brief: CallBrief{
				CustomerFacts:    " \t ",
				Objective:        valid.Objective,
				ApprovedPlaybook: valid.ApprovedPlaybook,
			},
			want: "customer facts are required",
		},
		{
			name: "objective is blank",
			brief: CallBrief{
				CustomerFacts:    valid.CustomerFacts,
				Objective:        " \n ",
				ApprovedPlaybook: valid.ApprovedPlaybook,
			},
			want: "call objective is required",
		},
		{
			name: "approved playbook is blank",
			brief: CallBrief{
				CustomerFacts:    valid.CustomerFacts,
				Objective:        valid.Objective,
				ApprovedPlaybook: " \t ",
			},
			want: "approved playbook text is required",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			before := test.brief
			err := test.brief.Validate()
			if err == nil {
				t.Fatal("Validate() error = nil, want required-field error")
			}
			if got := err.Error(); got != test.want {
				t.Fatalf("Validate() error = %q, want %q", got, test.want)
			}
			if strings.Contains(err.Error(), "private-sentinel") {
				t.Fatalf("Validate() error exposed input content: %q", err)
			}
			if !reflect.DeepEqual(test.brief, before) {
				t.Fatalf("Validate() mutated rejected brief: got %#v, want %#v", test.brief, before)
			}
		})
	}
}

func TestCallBriefValidateEnforcesRuneLimits(t *testing.T) {
	tests := []struct {
		name  string
		max   int
		field string
		set   func(*CallBrief, string)
	}{
		{
			name:  "customer facts",
			max:   MaxCustomerFactsRunes,
			field: "customer facts",
			set:   func(brief *CallBrief, value string) { brief.CustomerFacts = value },
		},
		{
			name:  "call objective",
			max:   MaxObjectiveRunes,
			field: "call objective",
			set:   func(brief *CallBrief, value string) { brief.Objective = value },
		},
		{
			name:  "approved playbook",
			max:   MaxApprovedPlaybookRunes,
			field: "approved playbook text",
			set:   func(brief *CallBrief, value string) { brief.ApprovedPlaybook = value },
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			brief := validCallBrief()
			test.set(&brief, strings.Repeat("界", test.max))
			if err := brief.Validate(); err != nil {
				t.Fatalf("Validate() rejected exactly %d Unicode characters: %v", test.max, err)
			}

			tooLong := strings.Repeat("界", test.max+1)
			test.set(&brief, tooLong)
			before := brief
			err := brief.Validate()
			want := test.field + " must be at most " + itoa(test.max) + " characters"
			if err == nil {
				t.Fatal("Validate() error = nil, want length error")
			}
			if got := err.Error(); got != want {
				t.Fatalf("Validate() error = %q, want %q", got, want)
			}
			if strings.Contains(err.Error(), "界") {
				t.Fatalf("Validate() error exposed input content: %q", err)
			}
			if !reflect.DeepEqual(brief, before) {
				t.Fatalf("Validate() mutated rejected brief: got %#v, want %#v", brief, before)
			}
		})
	}
}

func TestSalesScriptValidateRequiresCitationsForClaimsAndObjections(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*SalesScript)
		want   string
	}{
		{
			name: "claim without evidence",
			mutate: func(script *SalesScript) {
				script.Claims[0].EvidenceIDs = nil
			},
			want: "claim 1 must cite at least one playbook passage",
		},
		{
			name: "objection response without evidence",
			mutate: func(script *SalesScript) {
				script.Objections[0].Response.EvidenceIDs = nil
			},
			want: "objection response 1 must cite at least one playbook passage",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			script, evidence := validGroundedScript()
			test.mutate(&script)
			before := script
			err := script.ValidateStructure(evidence)
			if err == nil {
				t.Fatal("ValidateStructure() error = nil, want citation error")
			}
			if got := err.Error(); got != test.want {
				t.Fatalf("Validate() error = %q, want %q", got, test.want)
			}
			if !reflect.DeepEqual(script, before) {
				t.Fatalf("Validate() mutated rejected script: got %#v, want %#v", script, before)
			}
		})
	}
}

func TestSalesScriptValidateRejectsUnknownCitation(t *testing.T) {
	script, evidence := validGroundedScript()
	script.Claims[0].EvidenceIDs = []string{"not-in-playbook"}
	before := script

	err := script.ValidateStructure(evidence)
	if err == nil {
		t.Fatal("ValidateStructure() error = nil, want unknown-citation error")
	}
	if got, want := err.Error(), "claim 1 cites an unknown playbook passage"; got != want {
		t.Fatalf("Validate() error = %q, want %q", got, want)
	}
	if strings.Contains(err.Error(), "not-in-playbook") {
		t.Fatalf("Validate() error exposed citation input: %q", err)
	}
	if !reflect.DeepEqual(script, before) {
		t.Fatalf("Validate() mutated rejected script: got %#v, want %#v", script, before)
	}
}

func TestSalesScriptValidateRejectsMalformedEvidence(t *testing.T) {
	script, evidence := validGroundedScript()
	evidence = append(evidence, EvidencePassage{ID: "source-1", Text: "duplicate source"})

	err := script.ValidateStructure(evidence)
	if err == nil {
		t.Fatal("ValidateStructure() error = nil, want duplicate-evidence error")
	}
	if got, want := err.Error(), "duplicate playbook passage ID"; got != want {
		t.Fatalf("Validate() error = %q, want %q", got, want)
	}
}

func TestSalesScriptValidateRequiresQuestionsAndStopConditions(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*SalesScript)
		want   string
	}{
		{
			name:   "questions missing",
			mutate: func(script *SalesScript) { script.Questions = nil },
			want:   "at least one question is required",
		},
		{
			name:   "question blank",
			mutate: func(script *SalesScript) { script.Questions = []string{" \t "} },
			want:   "question 1 is required",
		},
		{
			name:   "stop conditions missing",
			mutate: func(script *SalesScript) { script.StopConditions = nil },
			want:   "at least one stop condition is required",
		},
		{
			name:   "stop condition blank",
			mutate: func(script *SalesScript) { script.StopConditions = []string{" \n "} },
			want:   "stop condition 1 is required",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			script, evidence := validGroundedScript()
			test.mutate(&script)
			before := script

			err := script.ValidateStructure(evidence)
			if err == nil {
				t.Fatal("ValidateStructure() error = nil, want required-section error")
			}
			if got := err.Error(); got != test.want {
				t.Fatalf("Validate() error = %q, want %q", got, test.want)
			}
			if !reflect.DeepEqual(script, before) {
				t.Fatalf("Validate() mutated rejected script: got %#v, want %#v", script, before)
			}
		})
	}
}

func TestSalesScriptValidateAcceptsGroundedDraft(t *testing.T) {
	script, evidence := validGroundedScript()
	if err := script.ValidateStructure(evidence); err != nil {
		t.Fatalf("ValidateStructure() error = %v, want nil", err)
	}
}

func validCallBrief() CallBrief {
	return CallBrief{
		CustomerFacts:    "The customer uses a small remote team.",
		Objective:        "Understand the team's approved requirements.",
		ApprovedPlaybook: "The product supports shared workspaces.",
	}
}

func validGroundedScript() (SalesScript, []EvidencePassage) {
	return SalesScript{
		Opening:   ScriptText{Text: "Thanks for making time."},
		Questions: []string{"What does your team need to coordinate?"},
		Objections: []ObjectionResponse{{
			Objection: "Will this support a remote team?",
			Response: ScriptText{
				Text:        "The approved playbook says shared workspaces are supported.",
				EvidenceIDs: []string{"source-1"},
			},
		}},
		Claims: []ScriptText{{
			Text:        "The product supports shared workspaces.",
			EvidenceIDs: []string{"source-1"},
		}},
		StopConditions: []string{"Stop if the required capability is not present in the approved playbook."},
	}, []EvidencePassage{{ID: "source-1", Text: "The product supports shared workspaces."}}
}

func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	var digits [20]byte
	index := len(digits)
	for value > 0 {
		index--
		digits[index] = byte('0' + value%10)
		value /= 10
	}
	return string(digits[index:])
}
