package consulclient

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	timeout, _ := time.ParseDuration("10s")
	_ = New("http://consul.service.example.consul:8500", timeout)
}
