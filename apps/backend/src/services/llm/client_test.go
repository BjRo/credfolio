package llm

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	apiKey := "test-key"
	client := NewClient(apiKey)
	if client == nil {
		t.Error("expected non-nil client")
	}
	if client.model != "gpt-4o" {
		t.Errorf("expected model 'gpt-4o', got %s", client.model)
	}
}

