package app

import (
	"testing"

	ft "github.com/openshift/ansible-service-broker/pkg/fusortest"
)

func TestBrokerConfig(t *testing.T) {
	config, err := CreateConfig("../../etc/test.config.yaml")
	if err != nil {
		t.Fatal(err)
	}

	ft.AssertNotNil(t, config, "message is nil")
	ft.AssertFalse(t, config.Broker.EnableAPBBind, "enable should be false")
}
