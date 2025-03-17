package optional

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptional_From(t *testing.T) {
	t.Parallel()

	var (
		expected = "test"
		value    = From[string](expected)
	)

	assert.Equal(t, true, value.set)
	assert.Equal(t, expected, value.value)
}

func TestOptional_FromPtr(t *testing.T) {
	t.Parallel()

	stub := "test"

	tests := []struct {
		name  string
		value *string
		set   bool
	}{
		{
			name:  "Value is not nil",
			value: &stub,
			set:   true,
		},
		{
			name:  "Value is nil",
			value: nil,
			set:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := FromPtr[string](test.value)
			assert.Equal(t, test.set, result.set)

			if test.value != nil {
				expected := *test.value
				assert.Equal(t, expected, result.value)
			}
		})
	}
}

func TestOptional_Value(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value Optional[string]
	}{
		{
			name: "Value is not set",
			value: Optional[string]{
				value: "",
				set:   false,
			},
		},
		{
			name: "Value is set",
			value: Optional[string]{
				value: "test",
				set:   true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, set := test.value.Value()

			assert.Equal(t, test.value.set, set)
			assert.Equal(t, test.value.value, result)
		})
	}
}

func TestOptional_IsSet(t *testing.T) {
	var (
		valueSet   = Optional[string]{set: true}
		valueUnset = Optional[string]{set: false}
	)

	assert.Equal(t, true, valueSet.IsSet())
	assert.Equal(t, false, valueUnset.IsSet())
}

func TestOptional_MarshalJSON(t *testing.T) {
	stub := "test"

	tests := []struct {
		name         string
		input        any
		expectedJSON string
	}{
		{
			name:         "Test",
			input:        struct{ Value Optional[string] }{Value: From("test")},
			expectedJSON: `{"Value": "test"}`,
		},
		{
			name:         "Test pointer",
			input:        struct{ Value Optional[*string] }{Value: From(&stub)},
			expectedJSON: `{"Value": "test"}`,
		},
		{
			name:         "Optional integer",
			input:        struct{ Value Optional[int] }{},
			expectedJSON: `{"Value": null}`,
		},
		{
			name:         "Optional optional",
			input:        struct{ Value Optional[Optional[uint64]] }{Value: From(From(uint64(1)))},
			expectedJSON: `{"Value": 1}`,
		},
		{
			name: "Optional struct",
			input: struct {
				Value Optional[struct{ Inner string }]
			}{
				Value: From(struct{ Inner string }{}),
			},
			expectedJSON: `{"Value": {"Inner": ""}}`,
		},
		{
			name: "Optional omitted struct",
			input: struct {
				Value Optional[struct{ Inner string }]
			}{},
			expectedJSON: `{"Value": null}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, err := json.Marshal(test.input)
			assert.NoError(t, err)
			assert.JSONEq(t, test.expectedJSON, string(data))
		})
	}
}

func TestOptional_UnmarshalJSON(t *testing.T) {
	var values struct {
		String        Optional[string]
		OmittedString Optional[string]
		StringPointer Optional[*string]
		NullInt       Optional[int]
	}

	err := json.Unmarshal([]byte(`
	{
		"String": "test",
		"StringPointer": "test pointer",
		"NullInt": null
	}`), &values)
	assert.NoError(t, err)

	strValue, _ := values.String.Value()
	assert.Equal(t, "test", strValue)

	assert.True(t, values.String.IsSet())
	assert.False(t, values.OmittedString.IsSet())
	assert.True(t, values.StringPointer.IsSet())

	var (
		ptrValue, _    = values.StringPointer.Value()
		strPtrValue, _ = values.StringPointer.Value()
	)

	if assert.NotNil(t, ptrValue) {
		assert.EqualValues(t, "test pointer", *strPtrValue)
	}

	assert.False(t, values.NullInt.IsSet())

	nullInt, _ := values.NullInt.Value()
	assert.Zero(t, nullInt)
}
