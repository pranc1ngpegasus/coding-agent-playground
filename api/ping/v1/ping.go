package pingv1

const (
	ServiceName = "example.ping.v1.PingService"
	Procedure   = "/" + ServiceName + "/Ping"
)

type PingRequest struct {
	Message string `json:"message"`
}

type PingResponse struct {
	Message string `json:"message"`
}
