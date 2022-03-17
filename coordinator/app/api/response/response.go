package response

type Response struct {
	Error     string      `json:"error"`
	ErrorCode int         `json:"error_code"`
	Data      interface{} `json:"data"`
}

