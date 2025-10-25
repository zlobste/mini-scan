package scanning_test

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/zlobste/mini-scan/pkg/scanning"
)

func TestDecodeMessage_V1Format(t *testing.T) {
	response := "hello world"
	encoded := base64.StdEncoding.EncodeToString([]byte(response))

	scan := scanning.Scan{
		Ip:          "192.168.1.1",
		Port:        80,
		Service:     "HTTP",
		Timestamp:   1000,
		DataVersion: scanning.V1,
		Data: map[string]any{
			"response_bytes_utf8": encoded,
		},
	}

	data, _ := json.Marshal(scan)
	decoded, err := scanning.DecodeMessage(data)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if decoded.ResponseText != response {
		t.Errorf("expected %q, got %q", response, decoded.ResponseText)
	}
	if decoded.IP != "192.168.1.1" {
		t.Errorf("expected IP 192.168.1.1, got %s", decoded.IP)
	}
}

func TestDecodeMessage_V2Format(t *testing.T) {
	response := "hello world"

	scan := scanning.Scan{
		Ip:          "192.168.1.2",
		Port:        443,
		Service:     "HTTPS",
		Timestamp:   2000,
		DataVersion: scanning.V2,
		Data: map[string]any{
			"response_str": response,
		},
	}

	data, _ := json.Marshal(scan)
	decoded, err := scanning.DecodeMessage(data)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if decoded.ResponseText != response {
		t.Errorf("expected %q, got %q", response, decoded.ResponseText)
	}
}

func TestDecodeMessage_InvalidJSON(t *testing.T) {
	_, err := scanning.DecodeMessage([]byte("invalid json"))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestDecodeMessage_MissingResponseField(t *testing.T) {
	scan := scanning.Scan{
		Ip:          "192.168.1.1",
		Port:        80,
		Service:     "HTTP",
		Timestamp:   1000,
		DataVersion: scanning.V1,
		Data: map[string]any{
			"wrong_field": "value",
		},
	}

	data, _ := json.Marshal(scan)
	_, err := scanning.DecodeMessage(data)

	if err == nil {
		t.Error("expected error for missing response field")
	}
}
