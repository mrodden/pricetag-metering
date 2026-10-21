package k8s

import (
	"encoding/json"
	"testing"
)

// TestOverlayDocParsing verifies the routing-overlay JSON structure the
// intelligent_route gateways deploy (praxis-routing-openai /
// praxis-routing-unified ConfigMaps) parses into overlayDoc correctly.
func TestOverlayDocParsing(t *testing.T) {
	// Shape captured from the live prod overlays (2026-10-21).
	raw := `{
		"selection_policy": {"mode": "deterministic"},
		"candidates": [
			{
				"kind": "inference_model",
				"name": "rits/deepseek-ai/DeepSeek-V4.1-Flash-1017",
				"site": "enmaas",
				"cluster": "cb-glm",
				"catalog": ["openai"],
				"owned_by": "curvebender"
			},
			{
				"kind": "inference_model",
				"name": "claude-haiku-5-5",
				"site": "enmaas",
				"cluster": "vertex",
				"provider": "vertex",
				"target_model": "vertex/claude-haiku-5-5",
				"paths": ["/v1/messages", "/v1/messages/count_tokens"],
				"catalog": ["anthropic"],
				"owned_by": "vertex"
			}
		]
	}`
	var doc overlayDoc
	if err := json.Unmarshal([]byte(raw), &doc); err != nil {
		t.Fatalf("parse overlay JSON: %v", err)
	}
	if len(doc.Candidates) != 2 {
		t.Fatalf("want 2 candidates, got %d", len(doc.Candidates))
	}
	ds := doc.Candidates[0]
	if ds.Name != "rits/deepseek-ai/DeepSeek-V4.1-Flash-1017" {
		t.Errorf("candidate[0].Name = %q", ds.Name)
	}
	if ds.Cluster != "cb-glm" {
		t.Errorf("candidate[0].Cluster = %q, want cb-glm", ds.Cluster)
	}
	if ds.OwnedBy != "curvebender" {
		t.Errorf("candidate[0].OwnedBy = %q, want curvebender", ds.OwnedBy)
	}
	haiku := doc.Candidates[1]
	if haiku.Provider != "vertex" {
		t.Errorf("candidate[1].Provider = %q, want vertex", haiku.Provider)
	}
	if len(haiku.Catalog) != 1 || haiku.Catalog[0] != "anthropic" {
		t.Errorf("candidate[1].Catalog = %v, want [anthropic]", haiku.Catalog)
	}
}

// TestOverlayDocParsingEmpty ensures a missing or empty candidates array
// parses without error (the overlay reader skips such CMs silently).
func TestOverlayDocParsingEmpty(t *testing.T) {
	for _, raw := range []string{`{}`, `{"candidates": []}`, `{"candidates": null}`} {
		var doc overlayDoc
		if err := json.Unmarshal([]byte(raw), &doc); err != nil {
			t.Fatalf("parse %q: %v", raw, err)
		}
		if len(doc.Candidates) != 0 {
			t.Errorf("%q: want 0 candidates, got %d", raw, len(doc.Candidates))
		}
	}
}
