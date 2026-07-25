package urlpath

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var assigmentRegexp = regexp.MustCompile(`/(\w+)(?:=(\w+))?`)

func Parse(path string) *Path {
	matches := assigmentRegexp.FindAllStringSubmatch(path, -1)
	segments := make([]Segment, len(matches))

	for i, m := range matches {
		segments[i] = Segment{m[1], m[2]}
	}

	return &Path{segments: segments}
}

type Path struct {
	segments []Segment
}

func (p *Path) String() string {
	var b strings.Builder
	for _, s := range p.segments {
		b.WriteString(s.String())
	}
	return b.String()
}

func (p *Path) Set(k, v string) {
	for i, s := range p.segments {
		if s.Key == k {
			p.segments[i].Value = v
			return
		}
	}

	p.segments = append(p.segments, Segment{k, v})
}

type Segment struct {
	Key, Value string
}

func (s Segment) String() string {
	if s.Value == "" {
		return "/" + url.PathEscape(s.Key)
	}
	return "/" + url.PathEscape(fmt.Sprintf("%s=%s", s.Key, s.Value))
}
