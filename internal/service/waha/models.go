package waha

type SessionStatus string

const (
	StatusStopped    SessionStatus = "STOPPED"
	StatusStarting   SessionStatus = "STARTING"
	StatusScanQRCode SessionStatus = "SCAN_QR_CODE"
	StatusWorking    SessionStatus = "WORKING"
	StatusFailed     SessionStatus = "FAILED"
)

type Me struct {
	ID       string `json:"id"`
	PushName string `json:"pushName"`
}

type Engine struct {
	Engine string `json:"engine"`
}

type WahaSession struct {
	Name   string        `json:"name"`
	Status SessionStatus `json:"status"`
	Me     Me            `json:"me"`
	Engine Engine        `json:"engine"`
}
