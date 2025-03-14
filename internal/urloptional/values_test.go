package urloptional

import (
	"net/url"
	"testing"

	"github.com/glichtv/kick-sdk/optional"
	"github.com/stretchr/testify/assert"
)

func TestSingle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected Value
	}{
		{
			name:     "Simple string",
			input:    "test",
			expected: Value{optional.From("test")},
		},
		{
			name:     "Empty string",
			input:    "",
			expected: Value{optional.From("")},
		},
		{
			name:     "String with spaces",
			input:    "test test",
			expected: Value{optional.From("test test")},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Single(test.input)
			assert.Equal(t, test.expected, result)

			if len(result) > 0 {
				value, set := result[0].Value()
				assert.True(t, set)
				assert.Equal(t, test.input, value)
			}
		})
	}
}

func TestMany(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    []string
		expected Value
	}{
		{
			name:  "Multiple strings",
			input: []string{"test", "test", "test"},
			expected: Value{
				optional.From("test"),
				optional.From("test"),
				optional.From("test"),
			},
		},
		{
			name:     "No values in slice",
			input:    []string{},
			expected: Value{},
		},
		{
			name:     "Single string",
			input:    []string{"test"},
			expected: Value{optional.From("test")},
		},
		{
			name:  "Strings with spaces",
			input: []string{"test test", "test test"},
			expected: Value{
				optional.From("test test"),
				optional.From("test test"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Many(test.input)
			assert.Equal(t, test.expected, result)

			for index, opt := range result {
				value, set := opt.Value()
				assert.True(t, set)
				assert.Equal(t, test.input[index], value)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       []string
		separator   string
		expected    Value
		joinedValue string
	}{
		{
			name:        "Join with separator",
			input:       []string{"test", "test", "test"},
			separator:   ",",
			expected:    Value{optional.From("test,test,test")},
			joinedValue: "test,test,test",
		},
		{
			name:        "Join empty slice",
			input:       []string{},
			separator:   ",",
			expected:    Value{optional.From("")},
			joinedValue: "",
		},
		{
			name:        "Join single value",
			input:       []string{"test"},
			separator:   ",",
			expected:    Value{optional.From("test")},
			joinedValue: "test",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Join(test.input, test.separator)
			assert.Equal(t, test.expected, result)

			if len(result) > 0 {
				value, set := result[0].Value()
				assert.True(t, set)
				assert.Equal(t, test.joinedValue, value)
			}
		})
	}
}

func TestValues_Encode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    Values
		expected string
	}{
		{
			name: "Single value",
			input: Values{
				"key": Single("value"),
			},
			expected: "key=value",
		},
		{
			name: "Multiple values for same key",
			input: Values{
				"key": Many([]string{"value1", "value2"}),
			},
			expected: "key=value1&key=value2",
		},
		{
			name: "Multiple keys",
			input: Values{
				"key1": Single("value1"),
				"key2": Single("value2"),
			},
			expected: "key1=value1&key2=value2",
		},
		{
			name: "Empty values",
			input: Values{
				"key": Value{},
			},
			expected: "",
		},
		{
			name: "Mixed values",
			input: Values{
				"single": Single("test"),
				"many":   Many([]string{"test", "test"}),
				"joined": Join([]string{"test", "test"}, ","),
			},
			expected: "joined=test%2Ctest&many=test&many=test&single=test",
		},
		{
			name: "Values needing URL encoding",
			input: Values{
				"space":   Single("test test"),
				"special": Single("test&test=test"),
			},
			expected: "space=test+test&special=test%26test%3Dtest",
		},
		{
			name: "Unset optional value",
			input: Values{
				"key": Value{optional.FromPtr[string](nil)},
			},
			expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.input.Encode()
			assert.Equal(t, test.expected, result)

			if result != "" {
				parsed, err := url.ParseQuery(result)
				assert.NoError(t, err)
				assert.NotNil(t, parsed)
			}
		})
	}
}
