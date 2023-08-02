package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/open-policy-agent/frameworks/constraint/pkg/externaldata"
)

const (
	apiVersion = "externaldata.gatekeeper.sh/v1beta1"
	kind       = "ProviderResponse"
)

// sendResponse sends back the response to Gatekeeper.
func SendResponse(results *[]externaldata.Item, systemErr string, w http.ResponseWriter) {
	response := externaldata.ProviderResponse{
		APIVersion: apiVersion,
		Kind:       kind,
		Response: externaldata.Response{
			Idempotent: true, // mutation requires idempotent results
		},
	}

	if results != nil {
		response.Response.Items = *results
	} else {
		response.Response.SystemError = systemErr
	}

	fmt.Printf("sending response response: %v\n", response)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Printf("unable to encode response, error: %s\n", err.Error())

		os.Exit(1)
	}
}

// Return containerName, containerImage
func GetContainerNameAndImageFromKey(externalProviderKey string) (string, string) {
	splited := strings.Split(externalProviderKey, "|")
	return splited[0], splited[1]
}
