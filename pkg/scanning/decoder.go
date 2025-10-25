package scanning

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type DecodedScan struct {
	IP           string
	Port         uint32
	Service      string
	Timestamp    int64
	ResponseText string
}

func DecodeMessage(data []byte) (DecodedScan, error) {
	var scan Scan
	if err := json.Unmarshal(data, &scan); err != nil {
		return DecodedScan{}, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	response := ""
	switch scan.DataVersion {
	case V1:
		v1Data, ok := scan.Data.(map[string]any)
		if !ok {
			return DecodedScan{}, fmt.Errorf("invalid V1 data format")
		}

		encoded, ok := v1Data["response_bytes_utf8"].(string)
		if !ok {
			return DecodedScan{}, fmt.Errorf("missing or invalid response_bytes_utf8")
		}

		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return DecodedScan{}, fmt.Errorf("failed to decode base64: %w", err)
		}
		response = string(decoded)

	case V2:
		v2Data, ok := scan.Data.(map[string]any)
		if !ok {
			return DecodedScan{}, fmt.Errorf("invalid V2 data format")
		}

		response, ok = v2Data["response_str"].(string)
		if !ok {
			return DecodedScan{}, fmt.Errorf("missing or invalid response_str")
		}

	default:
		return DecodedScan{}, fmt.Errorf("unknown data version: %d", scan.DataVersion)
	}

	return DecodedScan{
		IP:           scan.Ip,
		Port:         scan.Port,
		Service:      scan.Service,
		Timestamp:    scan.Timestamp,
		ResponseText: response,
	}, nil
}
