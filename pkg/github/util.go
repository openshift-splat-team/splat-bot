package github

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/sets"
	"net/http"
	"strings"
)

var knownAuthTypes = sets.New[string]("bearer", "basic", "negotiate")

// toCurl is a slightly adjusted copy of https://github.com/kubernetes/kubernetes/blob/74053d555d71a14e3853b97e204d7d6415521375/staging/src/k8s.io/client-go/transport/round_trippers.go#L339
func toCurl(r *http.Request) string {
	headers := ""
	for key, values := range r.Header {
		for _, value := range values {
			headers += fmt.Sprintf(` -H %q`, fmt.Sprintf("%s: %s", key, maskAuthorizationHeader(key, value)))
		}
	}

	return fmt.Sprintf("curl -k -v -X%s %s '%s'", r.Method, headers, r.URL.String())
}

// maskAuthorizationHeader masks credential content from authorization headers
// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Authorization
func maskAuthorizationHeader(key string, value string) string {
	if !strings.EqualFold(key, "Authorization") {
		return value
	}
	if len(value) == 0 {
		return ""
	}
	var authType string
	if i := strings.Index(value, " "); i > 0 {
		authType = value[0:i]
	} else {
		authType = value
	}
	if !knownAuthTypes.Has(strings.ToLower(authType)) {
		return "<masked>"
	}
	if len(value) > len(authType)+1 {
		value = authType + " <masked>"
	} else {
		value = authType
	}
	return value
}
