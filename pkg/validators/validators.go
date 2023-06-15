package validators

import "github.com/open-policy-agent/frameworks/constraint/pkg/externaldata"

type Validator interface {
	ValidKey(string, []externaldata.Item) []externaldata.Item
}

func WriteValidKey(key string, results []externaldata.Item) []externaldata.Item {
	results = append(results, externaldata.Item{
		Key:   key,
		Value: key + "_valid",
	})
	return results
}

func WriteInvalidKey(key string, results []externaldata.Item) []externaldata.Item {
	results = append(results, externaldata.Item{
		Key:   key,
		Error: key + "_invalid",
	})
	return results
}
