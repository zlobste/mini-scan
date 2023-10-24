package scanning

const (
	DataFormat = iota
	BinaryFormat
	JsonFormat
)

type Scan struct {
	Ip         string `json:"ip"`
	Port       uint32 `json:"port"`
	Service    string `json:"service"`
	Timestamp  int64  `json:"timestamp"`
	DataFormat int    `json:"data_format"`
	Data       []byte `json:"data"`
}
