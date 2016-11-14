package logredis

// LogstashMessageV0 represents v0 format
type LogstashMessageV0 struct {
	Type       string `json:"@type,omitempty"`
	Timestamp  string `json:"@timestamp"`
	Sourcehost string `json:"@source_host"`
	Message    string `json:"@message"`
	Fields     struct {
		Application string `json:"application"`
		File        string `json:"file"`
		Level       string `json:"level"`
		Labels      map[string]interface{}
	} `json:"@fields"`
}

// LogstashMessageV1 represents v1 format
type LogstashMessageV1 struct {
	Type        string `json:"@type,omitempty"`
	Timestamp   string `json:"@timestamp"`
	Sourcehost  string `json:"host"`
	Message     string `json:"message"`
	Application string `json:"application"`
	File        string `json:"file"`
	Level       string `json:"level"`
}
