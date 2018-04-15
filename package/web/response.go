package web

type Response struct {
	Error Error       `json:"error"`
	Date  interface{} `json:"data"`
}

type Error struct {
	Error      string      `json:"error"`
	Code       int         `json:"code"`
	StackTrace interface{} `json:"stack_trace,omitempty"`
}
