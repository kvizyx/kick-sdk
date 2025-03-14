package optional

import (
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
