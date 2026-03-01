package templates

type _baseResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
}
