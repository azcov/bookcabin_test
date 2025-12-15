package errorz

type WrappedError struct {
	StatusCode int            `json:"-"`
	ErrCode    string         `json:"error_code"`
	Msg        string         `json:"message"`
	Err        error          `json:"-"`
	Detail     map[string]any `json:"detail,omitempty"`
}

func (we *WrappedError) Error() string {
	return we.Msg
}

func New(statusCode int, errCode string, msg string) error {
	return &WrappedError{
		StatusCode: statusCode,
		ErrCode:    errCode,
		Msg:        msg,
	}
}

func (we *WrappedError) WithDetail(key string, value any) *WrappedError {
	if we.Detail == nil {
		we.Detail = make(map[string]any)
	}
	we.Detail[key] = value
	return we
}

func (we *WrappedError) WithDetails(details map[string]any) *WrappedError {
	if we.Detail == nil {
		we.Detail = make(map[string]any)
	}
	for k, v := range details {
		we.Detail[k] = v
	}
	return we
}
