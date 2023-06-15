package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"mesaglio/gatekeeper-external-data-provider-test/pkg/utils"
	"mesaglio/gatekeeper-external-data-provider-test/pkg/validators/cosign"
	"mesaglio/gatekeeper-external-data-provider-test/pkg/validators/naming"

	"github.com/open-policy-agent/frameworks/constraint/pkg/externaldata"
)

func Handler(w http.ResponseWriter, req *http.Request) {
	// only accept POST requests
	if req.Method != http.MethodPost {
		utils.SendResponse(nil, "only POST is allowed", w)
		return
	}

	// read request body
	requestBody, err := io.ReadAll(req.Body)
	if err != nil {
		utils.SendResponse(nil, fmt.Sprintf("unable to read request body: %v", err), w)
		return
	}

	fmt.Printf("received request body: %s\n", string(requestBody[:]))

	// parse request body
	var providerRequest externaldata.ProviderRequest
	err = json.Unmarshal(requestBody, &providerRequest)
	if err != nil {
		utils.SendResponse(nil, fmt.Sprintf("unable to unmarshal request body: %v", err), w)
		return
	}

	results := make([]externaldata.Item, 0)
	nv := naming.NamingValidator{}
	cv, err := cosign.New("")
	if err != nil {
		fmt.Printf("ERROR: cant initialize cosgin service: %s", err.Error())
		utils.SendResponse(nil, fmt.Sprintf("initialize cosgin service: %s", err.Error()), w)
		return
	}
	// iterate over all keys
	for _, key := range providerRequest.Request.Keys {
		// Providers should add a caching mechanism to avoid extra calls to external data sources.

		// following checks are for testing purposes only
		// check if key contains "_systemError" to trigger a system error
		if strings.HasSuffix(key, "_systemError") {
			utils.SendResponse(nil, "testing system error", w)
			return
		}

		results = nv.ValidKey(key, results)
		results = cv.ValidKey(key, results)

	}
	utils.SendResponse(&results, "", w)
}
