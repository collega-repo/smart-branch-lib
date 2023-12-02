package dto

type CfgDevice struct {
	DeviceId       string `json:"deviceId"`
	DeviceIdParent string `json:"deviceIdParent"`
	Name           string `json:"name"`
	Seq            int    `json:"seq"`
	UserLogin      string `json:"userLogin"`
	FullNameLogin  string `json:"fullNameLogin,omitempty"`
	UserRecv       string `json:"userRecv"`
	FullNameRecv   string `json:"fullNameRecv,omitempty"`
	Status         int    `json:"status"`
	Descr          string `json:"descr"`
}
