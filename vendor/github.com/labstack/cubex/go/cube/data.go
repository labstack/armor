package cube

type Data struct {
	Uptime         int64             `json:"uptime"`
	Request        int64             `json:"request"`
	ActiveRequest  int64             `json:"active_request"`
	BytesIn        int64             `json:"bytes_in"`
	BytesOut       int64             `json:"bytes_out"`
	AverageLatency int64             `json:"average_latency"`
	Endpoint       map[string]int64  `json:"endpoint"`
	UserAgent      map[string]int64  `json:"user_agent"`
	RemoteIP       map[string]int64  `json:"remote_ip"`
	Status         map[int32]int64   `json:"status"`
	Tags           map[string]string `json:"tags"`
}
