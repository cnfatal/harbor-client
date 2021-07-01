package client

import (
	"net/url"
	"strconv"
)

// Error same as https://github.com/goharbor/harbor/blob/4e1f6633afb824cd16341044a0e82f4f1f230cd2/src/lib/errors/errors.go#L32
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type CommonListOptions struct {
	// The page number
	// Default value: 1
	Page int `json:"page"`

	// The size of per page
	// Default value: 10
	Size int `json:"size"`
	// Query string to query resources.
	// Supported query patterns are "exact match(k=v)", "fuzzy match(k=~v)","range(k=[min~max])",
	// "list with union releationship(k={v1 v2 v3})" and"list with intersetion relationship(k=(v1 v2 v3))".
	// The value of range and list can be string(enclosed by " or '),integer or time(in format "2020-04-09 02:36:00").
	// All of these query patterns should be put in the query string "q=xxx" and splitted by ",". e.g. q=k1=v1,k2=~v2,k3=[min~max]
	Q string `json:"q"`
}

func (o *CommonListOptions) toQuery() url.Values {
	return url.Values{
		"page":      []string{strconv.Itoa(o.Page)},
		"page_size": []string{strconv.Itoa(o.Page)},
		"q":         []string{o.Q},
	}
}
