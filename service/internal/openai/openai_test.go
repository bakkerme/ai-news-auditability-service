package openai

import "testing"

func TestPreprocessJSON(t *testing.T) {
	client := &Client{}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic json with code blocks",
			input:    "```json\n{\"key\": \"value\"}\n```",
			expected: "{\"key\": \"value\"}",
		},
		{
			name:     "json with think tags",
			input:    "<think>Let me process this</think>```json\n{\"key\": \"value\"}\n```",
			expected: "{\"key\": \"value\"}",
		},
		{
			name:     "multiple think tags",
			input:    "<think>First thought</think>```json\n{\"data\": 123}\n```<think>Another thought</think>",
			expected: "{\"data\": 123}",
		},
		{
			name:     "think tags without json markers",
			input:    "<think>Some thought</think>{\"raw\": \"json\"}",
			expected: "{\"raw\": \"json\"}",
		},
		{
			name:     "nested json content",
			input:    "```json\n{\"outer\": {\"inner\": \"value\"}}\n```",
			expected: "{\"outer\": {\"inner\": \"value\"}}",
		},
		{
			name:     "think tags with whitespace",
			input:    "<think>\n  Processing...\n</think>```json\n{\"key\": \"value\"}\n```",
			expected: "{\"key\": \"value\"}",
		},
		{
			name:     "no json markers",
			input:    "{\"plain\": \"json\"}",
			expected: "{\"plain\": \"json\"}",
		},
		{
			name:     "json on new line",
			input:    "```\njson\n{\"key\": \"value\"}\n```",
			expected: "{\"key\": \"value\"}",
		},
		{
			name:     "just json word on line",
			input:    "json\n{\"key\": \"value\"}\n```",
			expected: "{\"key\": \"value\"}",
		},
		{
			name:     "json with carriage return newlines",
			input:    "```\r\njson\r\n{\"key\": \"value\"}\r\n```",
			expected: "{\"key\": \"value\"}",
		},
		{
			name:     "json with indented line",
			input:    "  json\n{\"key\": \"value\"}\n```",
			expected: "{\"key\": \"value\"}",
		},
		{
			name:     "json with tabs and spaces",
			input:    "\t  json  \n{\"key\": \"value\"}\n```",
			expected: "{\"key\": \"value\"}",
		},
		{
			name:     "empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "only think tags",
			input:    "<think>Just thinking</think>",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.PreprocessJSON(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestPreprocessYAML(t *testing.T) {
	client := &Client{}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic yaml with code blocks",
			input:    "```yaml\nkey: value\n```",
			expected: "key: value",
		},
		{
			name:     "yaml with think tags",
			input:    "<think>Let me process this</think>```yaml\nkey: value\n```",
			expected: "key: value",
		},
		{
			name:     "multiple think tags",
			input:    "<think>First thought</think>```yaml\ndata: 123\n```<think>Another thought</think>",
			expected: "data: 123",
		},
		{
			name:     "think tags without yaml markers",
			input:    "<think>Some thought</think>key: value",
			expected: "key: value",
		},
		{
			name:     "nested yaml content",
			input:    "```yaml\nouter:\n  inner: value\n```",
			expected: "outer:\n  inner: value",
		},
		{
			name:     "think tags with whitespace",
			input:    "<think>\n  Processing...\n</think>```yaml\nkey: value\n```",
			expected: "key: value",
		},
		{
			name:     "no yaml markers",
			input:    "plain: yaml",
			expected: "plain: yaml",
		},
		{
			name:     "yaml on new line",
			input:    "```\nyaml\nkey: value\n```",
			expected: "key: value",
		},
		{
			name:     "just yaml word on line",
			input:    "yaml\nkey: value\n```",
			expected: "key: value",
		},
		{
			name:     "yaml with carriage return newlines",
			input:    "```\r\nyaml\r\nkey: value\r\n```",
			expected: "key: value",
		},
		{
			name:     "yaml with indented line",
			input:    "  yaml\nkey: value\n```",
			expected: "key: value",
		},
		{
			name:     "yaml with tabs and spaces",
			input:    "\t  yaml  \nkey: value\n```",
			expected: "key: value",
		},
		{
			name:     "user example format",
			input:    "yaml\nid: \"t3_1kebauw\"\ntitle: \"Qwen AI Platform Shipping Announcement Sparks Community Speculation\"\noverview: \"The image analysis reveals a promotional graphic\"\n```",
			expected: "id: \"t3_1kebauw\"\ntitle: \"Qwen AI Platform Shipping Announcement Sparks Community Speculation\"\noverview: \"The image analysis reveals a promotional graphic\"",
		},
		{
			name:     "empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "only think tags",
			input:    "<think>Just thinking</think>",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.PreprocessYAML(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestNewWithSafeTimeouts(t *testing.T) {
	client := NewWithSafeTimeouts("http://localhost:8080", "test-key", "test-model")
	
	// Verify the client was created
	if client == nil {
		t.Fatal("expected client to be created, got nil")
	}
	
	// Verify the model is set correctly
	if client.model != "test-model" {
		t.Errorf("expected model to be 'test-model', got %q", client.model)
	}
	
	// Verify the safe timeout configuration is used
	if client.retry.MaxTotalTimeout != SafeOpenAIRetryConfig.MaxTotalTimeout {
		t.Errorf("expected MaxTotalTimeout to be %v, got %v", SafeOpenAIRetryConfig.MaxTotalTimeout, client.retry.MaxTotalTimeout)
	}
	
	if client.retry.MaxRetries != SafeOpenAIRetryConfig.MaxRetries {
		t.Errorf("expected MaxRetries to be %d, got %d", SafeOpenAIRetryConfig.MaxRetries, client.retry.MaxRetries)
	}
}
