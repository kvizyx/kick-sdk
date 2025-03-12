package urlext

import (
	"github.com/glichtv/kick-kit/optional"
	"net/url"
	"strings"
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

func (ov Values) Encode() string {
	values := make(url.Values, len(ov))

	for key, candidates := range ov {
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
