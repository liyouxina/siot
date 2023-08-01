package notify

type DeviceBreakMsg struct {
	MsgBase
	DeviceId   int64  `json:"device_id"`
	DeviceName string `json:"device_name"`
	Reportor   string `json:"reportor"`
	TargetUser int64  `json:"target_user"`
	TargetOrg  int64  `json:"target_org"`
}
