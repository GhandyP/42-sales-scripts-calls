package playbook

import (
	"encoding/json"
	"testing"

	"github.com/GhandyP/42-sales-scripts-calls/internal/sales"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/toolconfirmation"
)

type runnableSearchTool interface {
	Run(agent.Context, any) (map[string]any, error)
}

type noConfirmationContext struct {
	agent.ContextMock
}

func (*noConfirmationContext) ToolConfirmation() *toolconfirmation.ToolConfirmation {
	return nil
}

func TestNewSearchToolRequiresIndex(t *testing.T) {
	got, err := NewSearchTool(nil)
	if err == nil || err.Error() != "playbook index is required" {
		t.Fatalf("NewSearchTool(nil) error = %v, want playbook index validation error", err)
	}
	if got != nil {
		t.Fatalf("NewSearchTool(nil) returned a tool: %#v", got)
	}
}

func TestSearchToolRunsAgainstItsRequestScopedIndex(t *testing.T) {
	billingIndex, err := NewIndex("Billing refunds require manager approval within five business days.")
	if err != nil {
		t.Fatalf("NewIndex(billing) error = %v, want nil", err)
	}
	deliveryIndex, err := NewIndex("Delivery rescheduling requires one day notice.")
	if err != nil {
		t.Fatalf("NewIndex(delivery) error = %v, want nil", err)
	}

	billingTool, err := NewSearchTool(billingIndex)
	if err != nil {
		t.Fatalf("NewSearchTool(billing) error = %v, want nil", err)
	}
	if billingTool.Name() != "search_playbook" {
		t.Fatalf("tool name = %q, want search_playbook", billingTool.Name())
	}
	if billingTool.Description() == "" {
		t.Fatal("tool description is empty")
	}
	if billingTool.IsLongRunning() {
		t.Fatal("tool IsLongRunning() = true, want false")
	}

	deliveryTool, err := NewSearchTool(deliveryIndex)
	if err != nil {
		t.Fatalf("NewSearchTool(delivery) error = %v, want nil", err)
	}

	billingResults, err := invokeSearchTool(t, billingTool, "refunds", 1)
	if err != nil {
		t.Fatalf("billing tool invocation error = %v, want nil", err)
	}
	wantBilling := []sales.EvidencePassage{{
		ID:   "playbook-001",
		Text: "Billing refunds require manager approval within five business days.",
	}}
	if len(billingResults) != len(wantBilling) || billingResults[0] != wantBilling[0] {
		t.Fatalf("billing tool results = %#v, want %#v", billingResults, wantBilling)
	}

	deliveryResults, err := invokeSearchTool(t, deliveryTool, "refunds", 1)
	if err != nil {
		t.Fatalf("delivery tool invocation error = %v, want nil", err)
	}
	if len(deliveryResults) != 0 {
		t.Fatalf("delivery tool leaked another request's playbook: %#v", deliveryResults)
	}

	deliveryResults, err = invokeSearchTool(t, deliveryTool, "rescheduling", 1)
	if err != nil {
		t.Fatalf("delivery search invocation error = %v, want nil", err)
	}
	wantDelivery := []sales.EvidencePassage{{
		ID:   "playbook-001",
		Text: "Delivery rescheduling requires one day notice.",
	}}
	if len(deliveryResults) != len(wantDelivery) || deliveryResults[0] != wantDelivery[0] {
		t.Fatalf("delivery tool results = %#v, want %#v", deliveryResults, wantDelivery)
	}
}

func TestSearchToolPropagatesQueryValidation(t *testing.T) {
	index, err := NewIndex("Approved billing terms.")
	if err != nil {
		t.Fatalf("NewIndex() error = %v, want nil", err)
	}
	searchTool, err := NewSearchTool(index)
	if err != nil {
		t.Fatalf("NewSearchTool() error = %v, want nil", err)
	}

	results, err := invokeSearchTool(t, searchTool, " ! \n", 0)
	if err == nil || err.Error() != "search query must contain at least one term" {
		t.Fatalf("tool invocation error = %v, want query validation error", err)
	}
	if len(results) != 0 {
		t.Fatalf("tool returned results for invalid query: %#v", results)
	}
}

func invokeSearchTool(t *testing.T, searchTool tool.Tool, query string, limit int) ([]sales.EvidencePassage, error) {
	t.Helper()

	runnable, ok := searchTool.(runnableSearchTool)
	if !ok {
		t.Fatalf("tool %T does not expose the ADK runtime Run method", searchTool)
	}
	response, err := runnable.Run(&noConfirmationContext{}, map[string]any{
		"query": query,
		"limit": limit,
	})
	if err != nil {
		return nil, err
	}

	encoded, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("marshal tool response: %v", err)
	}
	var result struct {
		Passages []sales.EvidencePassage `json:"passages"`
	}
	if err := json.Unmarshal(encoded, &result); err != nil {
		t.Fatalf("decode tool response: %v", err)
	}
	return result.Passages, nil
}
