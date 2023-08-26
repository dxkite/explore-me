package actions

const (
	ParamError    = "param_error"
	InternalError = "internal_error"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
