package common

type Sync struct {
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}
