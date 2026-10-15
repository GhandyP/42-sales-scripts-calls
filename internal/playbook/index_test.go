package playbook

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/GhandyP/42-sales-scripts-calls/internal/sales"
)

func TestNewIndexCreatesStableExactPassagesAndReturnsCopies(t *testing.T) {
	source := " \tPricing starts at $50.\n\nRefunds require approval. \n"
	index, err := NewIndex(source)
	if err != nil {
		t.Fatalf("NewIndex() error = %v, want nil", err)
	}

	first, err := index.Search("PRICING", 1)
	if err != nil {
		t.Fatalf("Search() error = %v, want nil", err)
	}
	wantFirst := []sales.EvidencePassage{{ID: "playbook-001", Text: "Pricing starts at $50."}}
	if !reflect.DeepEqual(first, wantFirst) {
		t.Fatalf("Search(PRICING) = %#v, want %#v", first, wantFirst)
	}

	first[0].Text = "caller mutation"
	again, err := index.Search("pricing", 1)
	if err != nil {
		t.Fatalf("Search() after result mutation error = %v, want nil", err)
	}
	if !reflect.DeepEqual(again, wantFirst) {
		t.Fatalf("Search() exposed mutable index data: got %#v, want %#v", again, wantFirst)
	}

	second, err := index.Search("refunds", 1)
	if err != nil {
		t.Fatalf("Search(refunds) error = %v, want nil", err)
	}
	wantSecond := []sales.EvidencePassage{{ID: "playbook-002", Text: "Refunds require approval."}}
	if !reflect.DeepEqual(second, wantSecond) {
		t.Fatalf("Search(refunds) = %#v, want %#v", second, wantSecond)
	}
}

func TestNewIndexSplitsLongPassagesAtRuneLimit(t *testing.T) {
	source := strings.Repeat("offer ", 200) + " finalterm"
	index, err := NewIndex(source)
	if err != nil {
		t.Fatalf("NewIndex() error = %v, want nil", err)
	}

	got, err := index.Search("offer", 5)
	if err != nil {
		t.Fatalf("Search() error = %v, want nil", err)
	}
	if len(got) != 2 {
		t.Fatalf("Search() returned %d passages, want 2: %#v", len(got), got)
	}

	cursor := 0
	for i, passage := range got {
		wantID := fmt.Sprintf("playbook-%03d", i+1)
		if passage.ID != wantID {
			t.Errorf("passage %d ID = %q, want %q", i, passage.ID, wantID)
		}
		if gotRunes := utf8.RuneCountInString(passage.Text); gotRunes > 1000 {
			t.Errorf("passage %q has %d runes, want at most 1000", passage.ID, gotRunes)
		}
		relative := strings.Index(source[cursor:], passage.Text)
		if relative < 0 {
			t.Errorf("passage %q is not an exact ordered excerpt of the source", passage.ID)
			continue
		}
		cursor += relative + len(passage.Text)
	}
}

func TestIndexSearchRanksDistinctTermsAndKeepsSourceOrderOnTies(t *testing.T) {
	source := "Alpha beta launch policy.\n\nAlpha launch policy.\n\nBeta alpha offer.\n\nAlpha alpha."
	index, err := NewIndex(source)
	if err != nil {
		t.Fatalf("NewIndex() error = %v, want nil", err)
	}

	got, err := index.Search("ALPHA beta", 5)
	if err != nil {
		t.Fatalf("Search() error = %v, want nil", err)
	}
	assertPassageIDs(t, got, []string{"playbook-001", "playbook-003", "playbook-002", "playbook-004"})
}

func TestIndexSearchDefaultsAndBoundsResultCount(t *testing.T) {
	source := strings.Join([]string{
		"launch policy one.",
		"launch policy two.",
		"launch policy three.",
		"launch policy four.",
		"launch policy five.",
		"launch policy six.",
	}, "\n\n")
	index, err := NewIndex(source)
	if err != nil {
		t.Fatalf("NewIndex() error = %v, want nil", err)
	}

	defaultResults, err := index.Search("launch", 0)
	if err != nil {
		t.Fatalf("Search() with default limit error = %v, want nil", err)
	}
	if len(defaultResults) != 3 {
		t.Fatalf("Search() with default limit returned %d passages, want 3", len(defaultResults))
	}

	maximumResults, err := index.Search("launch", 5)
	if err != nil {
		t.Fatalf("Search() with maximum limit error = %v, want nil", err)
	}
	if len(maximumResults) != 5 {
		t.Fatalf("Search() with maximum limit returned %d passages, want 5", len(maximumResults))
	}
	assertPassageIDs(t, maximumResults, []string{"playbook-001", "playbook-002", "playbook-003", "playbook-004", "playbook-005"})

	for _, limit := range []int{-1, 6} {
		t.Run(fmt.Sprintf("invalid limit %d", limit), func(t *testing.T) {
			results, err := index.Search("launch", limit)
			if err == nil || err.Error() != "search limit must be between 1 and 5" {
				t.Fatalf("Search() error = %v, want search limit validation error", err)
			}
			if len(results) != 0 {
				t.Fatalf("Search() returned results for invalid limit: %#v", results)
			}
		})
	}
}

func TestIndexSearchRejectsBlankQueryAndReturnsEmptyOnMiss(t *testing.T) {
	index, err := NewIndex("Policy details are approved.")
	if err != nil {
		t.Fatalf("NewIndex() error = %v, want nil", err)
	}

	results, err := index.Search(" \n ! ", 0)
	if err == nil || err.Error() != "search query must contain at least one term" {
		t.Fatalf("Search() error = %v, want query validation error", err)
	}
	if len(results) != 0 {
		t.Fatalf("Search() returned results for blank query: %#v", results)
	}

	results, err = index.Search("unrelated", 0)
	if err != nil {
		t.Fatalf("Search() for no-match query error = %v, want nil", err)
	}
	if len(results) != 0 {
		t.Fatalf("Search() for no-match query returned %#v, want no passages", results)
	}
}

func TestNewIndexRejectsMissingAndOversizedPlaybooks(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "missing playbook",
			input: " \n\t ",
			want:  "approved playbook text is required",
		},
		{
			name:  "oversized playbook",
			input: strings.Repeat("x", sales.MaxApprovedPlaybookRunes+1),
			want:  "approved playbook text must be at most 20000 characters",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			index, err := NewIndex(test.input)
			if err == nil || err.Error() != test.want {
				t.Fatalf("NewIndex() error = %v, want %q", err, test.want)
			}
			if index != nil {
				t.Fatalf("NewIndex() returned an index for rejected input: %#v", index)
			}
		})
	}
}

func TestIndexesRemainScopedToTheirPlaybook(t *testing.T) {
	billing, err := NewIndex("Billing approval takes two days.")
	if err != nil {
		t.Fatalf("NewIndex(billing) error = %v, want nil", err)
	}
	delivery, err := NewIndex("Delivery approval takes one day.")
	if err != nil {
		t.Fatalf("NewIndex(delivery) error = %v, want nil", err)
	}

	billingResults, err := billing.Search("billing", 3)
	if err != nil {
		t.Fatalf("billing.Search() error = %v, want nil", err)
	}
	assertPassageIDs(t, billingResults, []string{"playbook-001"})
	if billingResults[0].Text != "Billing approval takes two days." {
		t.Fatalf("billing result = %q, want only the billing playbook", billingResults[0].Text)
	}

	deliveryResults, err := delivery.Search("delivery", 3)
	if err != nil {
		t.Fatalf("delivery.Search() error = %v, want nil", err)
	}
	assertPassageIDs(t, deliveryResults, []string{"playbook-001"})
	if deliveryResults[0].Text != "Delivery approval takes one day." {
		t.Fatalf("delivery result = %q, want only the delivery playbook", deliveryResults[0].Text)
	}
	crossResults, err := billing.Search("delivery", 3)
	if err != nil {
		t.Fatalf("billing.Search(delivery) error = %v, want nil", err)
	}
	if len(crossResults) != 0 {
		t.Fatalf("billing index leaked results from another playbook: %#v", crossResults)
	}
}

func assertPassageIDs(t *testing.T, passages []sales.EvidencePassage, want []string) {
	t.Helper()
	if len(passages) != len(want) {
		t.Fatalf("got %d passages, want %d: %#v", len(passages), len(want), passages)
	}
	for i, passage := range passages {
		if passage.ID != want[i] {
			t.Errorf("passage[%d].ID = %q, want %q", i, passage.ID, want[i])
		}
	}
}
