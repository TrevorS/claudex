// ABOUTME: Message type for JSONL serialization and deserialization
package jsonl

// Message represents a single message line in a JSONL file
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Tokens  int64  `json:"tokens"`
}
