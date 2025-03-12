package urlext

import (
	"github.com/glichtv/kick-kit/optional"
	"net/url"
	"strings"
)

type (
	Value  = optional.Optional[string]
	Values map[string]Value
)

func Single(value string) Value {
	return optional.From(value)
}

func Join(values []string, separator string) Value {
	return optional.From(
		strings.Join(values, separator),
	)
}

func (ov Values) Encode() string {
	values := make(url.Values, len(ov))

	for key, candidate := range ov {
		value, set := candidate.Value()
		if set {
			values.Add(key, value)
		}
	}

	return values.Encode()
}
