package jsonl

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParser_ValidMessages tests parsing a valid JSONL file with well-formed messages
func TestParser_ValidMessages(t *testing.T) {
	jsonl := `{"role":"user","content":"Hello","tokens":2}
{"role":"assistant","content":"Hi there","tokens":3}
{"role":"user","content":"How are you?","tokens":4}`

	parser := New(strings.NewReader(jsonl))
	require.NotNil(t, parser)

	// First message
	msg1, err := parser.Next()
	require.NoError(t, err)
	require.NotNil(t, msg1)
	assert.Equal(t, "user", msg1.Role)
	assert.Equal(t, "Hello", msg1.Content)
	assert.Equal(t, int64(2), msg1.Tokens)

	// Second message
	msg2, err := parser.Next()
	require.NoError(t, err)
	require.NotNil(t, msg2)
	assert.Equal(t, "assistant", msg2.Role)
	assert.Equal(t, "Hi there", msg2.Content)
	assert.Equal(t, int64(3), msg2.Tokens)

	// Third message
	msg3, err := parser.Next()
	require.NoError(t, err)
	require.NotNil(t, msg3)
	assert.Equal(t, "user", msg3.Role)
	assert.Equal(t, "How are you?", msg3.Content)
	assert.Equal(t, int64(4), msg3.Tokens)

	// EOF after all messages
	msg4, err := parser.Next()
	assert.Nil(t, msg4)
	assert.Equal(t, io.EOF, err)
}

// TestParser_EmptyFile tests handling of an empty JSONL file
func TestParser_EmptyFile(t *testing.T) {
	parser := New(strings.NewReader(""))
	require.NotNil(t, parser)

	msg, err := parser.Next()
	assert.Nil(t, msg)
	assert.Equal(t, io.EOF, err)
}

// TestParser_SingleMessage tests parsing a file with a single message
func TestParser_SingleMessage(t *testing.T) {
	jsonl := `{"role":"user","content":"Single message","tokens":2}`

	parser := New(strings.NewReader(jsonl))
	msg, err := parser.Next()
	require.NoError(t, err)
	require.NotNil(t, msg)
	assert.Equal(t, "user", msg.Role)
	assert.Equal(t, "Single message", msg.Content)

	// EOF on next call
	msg2, err := parser.Next()
	assert.Nil(t, msg2)
	assert.Equal(t, io.EOF, err)
}

// TestParser_MalformedJSON tests handling of invalid JSON lines
func TestParser_MalformedJSON(t *testing.T) {
	jsonl := `{"role":"user","content":"Valid","tokens":1}
{invalid json without quotes
{"role":"assistant","content":"Also valid","tokens":2}`

	parser := New(strings.NewReader(jsonl))

	// First valid message
	msg1, err := parser.Next()
	require.NoError(t, err)
	require.NotNil(t, msg1)
	assert.Equal(t, "user", msg1.Role)

	// Malformed JSON should return error
	msg2, err := parser.Next()
	assert.Nil(t, msg2)
	assert.Error(t, err)
	assert.NotEqual(t, io.EOF, err)

	// Parser should continue after malformed line
	msg3, err := parser.Next()
	require.NoError(t, err)
	require.NotNil(t, msg3)
	assert.Equal(t, "assistant", msg3.Role)
}

// TestParser_BlankLines tests handling of blank lines
func TestParser_BlankLines(t *testing.T) {
	jsonl := `{"role":"user","content":"First","tokens":1}

{"role":"assistant","content":"Second","tokens":2}`

	parser := New(strings.NewReader(jsonl))

	msg1, err := parser.Next()
	require.NoError(t, err)
	require.NotNil(t, msg1)
	assert.Equal(t, "user", msg1.Role)

	// Blank line should be skipped or error
	// Depending on implementation, blank lines might be handled gracefully
	msg2, err := parser.Next()
	// Accept either skipping blank lines (returns next valid message)
	// or returning an error for blank line
	if err != nil {
		// If error, next message after should be valid
		msg3, err := parser.Next()
		require.NoError(t, err)
		require.NotNil(t, msg3)
		assert.Equal(t, "assistant", msg3.Role)
	} else if msg2 != nil {
		// If skipped, we got the second message
		assert.Equal(t, "assistant", msg2.Role)
	}
}

// TestParser_SpecialCharacters tests messages with special characters and escapes
func TestParser_SpecialCharacters(t *testing.T) {
	jsonl := `{"role":"user","content":"Message with \"quotes\" and newlines","tokens":5}
{"role":"assistant","content":"Special chars: !@#$%^&*()","tokens":6}`

	parser := New(strings.NewReader(jsonl))

	msg1, err := parser.Next()
	require.NoError(t, err)
	require.NotNil(t, msg1)
	assert.Contains(t, msg1.Content, "quotes")

	msg2, err := parser.Next()
	require.NoError(t, err)
	require.NotNil(t, msg2)
	assert.Contains(t, msg2.Content, "!@#$%^&*()")
}

// TestParser_LargeFile tests parsing a larger JSONL file
func TestParser_LargeFile(t *testing.T) {
	var buf bytes.Buffer
	count := 100
	for i := 0; i < count; i++ {
		role := "user"
		if i%2 == 0 {
			role = "assistant"
		}
		buf.WriteString(`{"role":"` + role + `","content":"Message ` + string(rune('0'+i%10)) + `","tokens":` + string(rune('1'+(i%5))) + `}` + "\n")
	}

	parser := New(&buf)
	msgCount := 0
	for {
		msg, err := parser.Next()
		if err == io.EOF {
			break
		}
		require.NoError(t, err, "iteration %d", msgCount)
		require.NotNil(t, msg)
		msgCount++
	}

	assert.Equal(t, count, msgCount)
}

// TestParser_StreamingInterface tests that parser properly implements streaming interface
func TestParser_StreamingInterface(t *testing.T) {
	jsonl := `{"role":"user","content":"Line 1","tokens":1}
{"role":"assistant","content":"Line 2","tokens":2}
{"role":"user","content":"Line 3","tokens":3}`

	parser := New(strings.NewReader(jsonl))

	// Verify parser can be called repeatedly
	messages := []*Message{}
	for {
		msg, err := parser.Next()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		messages = append(messages, msg)
	}

	assert.Equal(t, 3, len(messages))
	assert.Equal(t, "user", messages[0].Role)
	assert.Equal(t, "assistant", messages[1].Role)
	assert.Equal(t, "user", messages[2].Role)
}

// TestParser_MissingFields tests handling of JSON with missing required fields
// Note: Go's json.Unmarshal doesn't error on missing fields, it just leaves them empty
// This test verifies that behavior
func TestParser_MissingFields(t *testing.T) {
	jsonl := `{"content":"Missing role field","tokens":1}
{"role":"user","content":"Has all fields","tokens":2}`

	parser := New(strings.NewReader(jsonl))

	// First line missing role field - will unmarshal but role will be empty
	msg1, err := parser.Next()
	require.NoError(t, err)
	require.NotNil(t, msg1)
	assert.Equal(t, "", msg1.Role) // Role is empty string when not provided

	// Second message with all fields
	msg2, err := parser.Next()
	require.NoError(t, err)
	require.NotNil(t, msg2)
	assert.Equal(t, "user", msg2.Role)
}

// TestParser_ExtraFields tests handling of JSON with extra fields
func TestParser_ExtraFields(t *testing.T) {
	jsonl := `{"role":"user","content":"Message","tokens":2,"extra":"field","another":123}
{"role":"assistant","content":"Response","tokens":3,"timestamp":"2024-01-01"}`

	parser := New(strings.NewReader(jsonl))

	msg1, err := parser.Next()
	require.NoError(t, err)
	require.NotNil(t, msg1)
	assert.Equal(t, "user", msg1.Role)
	assert.Equal(t, "Message", msg1.Content)

	msg2, err := parser.Next()
	require.NoError(t, err)
	require.NotNil(t, msg2)
	assert.Equal(t, "assistant", msg2.Role)
}

// TestParser_WithByteReader tests parser with io.Reader implementation
func TestParser_WithByteReader(t *testing.T) {
	jsonl := `{"role":"user","content":"Test","tokens":1}
{"role":"assistant","content":"OK","tokens":2}`

	reader := bytes.NewReader([]byte(jsonl))
	parser := New(reader)
	require.NotNil(t, parser)

	msg1, err := parser.Next()
	require.NoError(t, err)
	require.NotNil(t, msg1)

	msg2, err := parser.Next()
	require.NoError(t, err)
	require.NotNil(t, msg2)

	msg3, err := parser.Next()
	assert.Nil(t, msg3)
	assert.Equal(t, io.EOF, err)
}

// TestParser_WhitespaceHandling tests handling of whitespace in JSON
func TestParser_WhitespaceHandling(t *testing.T) {
	jsonl := `  {"role":"user","content":"Indented","tokens":1}
{"role":"assistant", "content":"Normal", "tokens": 2}
	{"role":"user","content":"Tabs","tokens":3}	`

	parser := New(strings.NewReader(jsonl))

	// First message with leading whitespace
	msg, err := parser.Next()
	// Leading/trailing whitespace handling depends on implementation
	// Should either skip whitespace or error on it
	if err != nil {
		// If error on leading whitespace, continue
		msg, err = parser.Next()
	}
	require.NoError(t, err)
	require.NotNil(t, msg)
}
