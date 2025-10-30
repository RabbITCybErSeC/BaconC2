package models

type HttpBeaconResponse struct {
	Command        Command       `json:"command"`
	Status         CommandStatus `json:"status,omitempty"`
	NextBeacon     int           `json:"nextBeacon,omitempty"`
	RequestResults bool          `json:"requestResults"`
}
