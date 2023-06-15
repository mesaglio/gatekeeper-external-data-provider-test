package naming

import (
	"mesaglio/gatekeeper-external-data-provider-test/pkg/validators"
	"strings"

	"github.com/open-policy-agent/frameworks/constraint/pkg/externaldata"
)

type NamingValidator struct {
}

func (nv NamingValidator) ValidKey(key string, results []externaldata.Item) []externaldata.Item {
	// check if key contains "error_" to trigger an error
	if strings.HasPrefix(key, "error_") {
		results = validators.WriteInvalidKey(key, results)
	} else if !strings.HasSuffix(key, "_valid") {
		// valid key will have "_valid" appended as return value
		results = validators.WriteValidKey(key, results)
	}
	return results
}
