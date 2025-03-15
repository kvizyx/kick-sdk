package urloptional

import (
	"net/url"
	"strings"

	"github.com/glichtv/kick-sdk/optional"
)

type (
	Value  = []optional.Optional[string]
	Values map[string]Value
)

func Single(value string) Value {
	return Value{
		optional.From(value),
	}
}

func SingleOptional(value optional.Optional[string]) Value {
	return Value{value}
}

func Many(values []string) Value {
	origin := make(Value, len(values))

	for index, value := range values {
		origin[index] = optional.From(value)
	}

	return origin
}

func Join(values []string, separator string) Value {
	return Value{
		optional.From(
			strings.Join(values, separator),
		),
	}
}

func (v Values) Encode() string {
	values := make(url.Values, len(v))

	for key, candidates := range v {
		for _, candidate := range candidates {
			value, set := candidate.Value()
			if !set {
				continue
			}

			values.Add(key, value)
		}
	}

	return values.Encode()
}
