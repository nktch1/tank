package parser

import (
	"context"
	"strings"
)

type Parser interface {
	ParseSearchResponse(ctx context.Context, rawResponse []byte) Response
}

func getRootDomain(domain string) string {
	domain = strings.ToLower(domain)

	parts := strings.Split(domain, ".")
	if len(parts) < 3 {
		return domain
	}

	if _, ok := tlds[strings.Join(parts[len(parts)-2:], ".")]; ok {
		return strings.Join(parts[len(parts)-3:], ".")
	}

	return strings.Join(parts[len(parts)-2:], ".")
}