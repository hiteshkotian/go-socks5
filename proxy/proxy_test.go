package proxy

import (
	"os"
	"testing"
)

func TestDaemonNewInvalidID(t *testing.T) {
	daemon, err := DaemonNew("", "")
	if err == nil {
		t.Errorf("Error should be returned for invalid ID")
		return
	}

	if daemon != nil {
		t.Errorf("Daemon returned for invalid parameters should be null")
		return
	}
}

func TestDaemonNewDefaultProxyID(t *testing.T) {
	daemon, err := DaemonNew("1010", "")
	expected, err := os.Hostname()
	if err != nil {
		t.Errorf("Error should not be returned")
		return
	}

	if daemon == nil {
		t.Errorf("Invalid daemon provided")
		return
	}

	if expected != daemon.proxyID {
		t.Errorf("Proxy ID expected : %s but was set to %s", expected, daemon.proxyID)
	}
}

func TestDaemonNewWithProxyID(t *testing.T) {
	expected := "proxy-id-str"
	daemon, err := DaemonNew("1010", expected)
	if err != nil {
		t.Errorf("Error should not be returned")
		return
	}

	if daemon == nil {
		t.Errorf("Invalid daemon provided")
		return
	}

	if expected != daemon.proxyID {
		t.Errorf("Proxy ID expected : %s but was set to %s", expected, daemon.proxyID)
	}
}

func TestDaemonNew(t *testing.T) {
	expected := "proxy-id-str"
	expectedID := "1010"
	daemon, err := DaemonNew(expectedID, expected)
	if err != nil {
		t.Errorf("Error should not be returned")
		return
	}

	if daemon == nil {
		t.Errorf("Invalid daemon provided")
		return
	}

	if expected != daemon.proxyID {
		t.Errorf("Proxy ID expected : %s but was set to %s", expected, daemon.proxyID)
	}

	if expectedID != daemon.id {
		t.Errorf("ID expected : %s but was set to %s", expectedID, daemon.id)
	}

	if daemon.state != Initializing {
		t.Errorf("State not set to Initializing")
	}
}
