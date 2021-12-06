package step

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var paramPattern = regexp.MustCompile(`{((?:int)|(?:string))}`)

type stepMatcher struct {
	re    *regexp.Regexp
	types []string
}

func newStepMatcher(pattern string) (*stepMatcher, error) {
	all := paramPattern.FindAllStringIndex(pattern, -1)

	sb := new(strings.Builder)

	last := 0

	types := []string{}

	sb.WriteString("^")

	for _, sl := range all {
		if sl[0] > last {
			sb.WriteString(pattern[last:sl[0]])
		}

		t := pattern[sl[0]+1 : sl[1]-1]
		types = append(types, t)

		switch t {
		case "int":
			sb.WriteString(`(-?\d+)`)
		case "string":
			sb.WriteString(`"((?:[^"]|(?:\\"))*)"`)
		default:
			return nil, fmt.Errorf("unknown parameter type %q", t)
		}

		last = sl[1]

	}

	if last < len(pattern) {
		sb.WriteString(pattern[last:])
	}

	sb.WriteString("$")

	matchRegexp := sb.String()

	mat, err := regexp.Compile(matchRegexp)

	if err != nil {
		return nil, fmt.Errorf("while parsing regexp: %w", err)
	}

	return &stepMatcher{
		re:    mat,
		types: types,
	}, nil
}

func (m stepMatcher) match(txt string) ([]interface{}, error) {
	sm := m.re.FindStringSubmatch(txt)
	if sm == nil {
		return nil, errNotMatching
	}

	if len(sm)-1 != len(m.types) {
		return nil, fmt.Errorf("something went wrong matching types (%d) and groups (%d)", len(m.types), len(sm)-1)
	}

	params := make([]interface{}, len(m.types))

	for i, st := range sm[1:] {
		t := m.types[i]
		switch t {
		case "int":
			v, err := strconv.ParseInt(st, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("while parsing int %q: %w", st, err)
			}
			params[i] = int(v)
		case "string":
			params[i] = strings.ReplaceAll(st, "\\\"", "\"")
		default:
			return nil, fmt.Errorf("unknown parameter type %q", t)
		}
	}

	return params, nil

}

func Match(pattern, txt string) ([]interface{}, error) {
	matcher, err := newStepMatcher(pattern)
	if err != nil {
		return nil, err
	}
	return matcher.match(txt)
}
