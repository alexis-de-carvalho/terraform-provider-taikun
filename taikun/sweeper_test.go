package taikun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/itera-io/taikungoclient/client"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func sharedConfig() (interface{}, error) {
	email := os.Getenv("TAIKUN_EMAIL")
	password := os.Getenv("TAIKUN_PASSWORD")

	if email == "" {
		return nil, fmt.Errorf("TAIKUN_EMAIL must be set for acceptance tests")
	}
	if password == "" {
		return nil, fmt.Errorf("TAIKUN_PASSWORD must be set for acceptance tests")
	}

	return &apiClient{
		client:              client.Default,
		email:               email,
		password:            password,
		useKeycloakEndpoint: false,
	}, nil
}