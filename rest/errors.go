package rest

const (
	ErrUnauthorized        = "ErrUnauthorized"
	ErrCodeNetwork         = "ErrNetwork"
	ErrCodeInternalError   = "ErrInternalError"
	ErrCodeBadResponseBody = "ErrBadResponseBody"
	ErrCodeUnknown         = "Unknown"
	ErrCodeOk              = "Ok"
)

type Error struct {
	StatusCode int    `json:"-"`
	ErrCode    string `json:"err_code,omitempty"`
	Message    string `json:"message,omitempty"`
}

func (e *Error) Status(status int) *Error {
	e.StatusCode = status
	return e
}

func (e *Error) Err(err string) *Error {
	e.ErrCode = err
	return e
}

func (e *Error) Error() string {
	return e.String()
}

func (e *Error) String() string {
	if e == nil {
		return "<nil>"
	}
	//return fmt.Sprintf("status=%d, err_code=%s, err_msg=%s", e.StatusCode, e.ErrCode, e.Message)
	return e.Message
}

func IsRestError(err error) (*Error, bool) {
	if err == nil {
		return nil, false
	}

	ce, ok := err.(*Error)
	return ce, ok
}

func NewError(statusCode int, errMsg string) *Error {
	return &Error{
		StatusCode: statusCode,
		ErrCode:    "",
		Message:    errMsg,
	}
}

func NewErrorWithCode(statusCode int, errCode string, errMsg string) *Error {
	return &Error{
		StatusCode: statusCode,
		ErrCode:    errCode,
		Message:    errMsg,
	}
}
