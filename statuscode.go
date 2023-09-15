package libgemini

import "fmt"

type StatusCode uint16 // two bytes

const (
	Input             StatusCode = 10
	SensitiveInput    StatusCode = 11
	Success           StatusCode = 20
	RedirectTemporary            = 30
	RedirectPermanent            = 31
	TemporaryFailure             = 40
	ServerUnavailable            = 41
	CGIError                     = 42
	ProxyError                   = 43
	SlowDown                     = 44
	PermanentFailure             = 50
	NotFound                     = 51
)

func (code StatusCode) IsSuccess() bool {
	return code >= Success && code < RedirectTemporary
}

func (code StatusCode) String() string {
	name := "Unknown (Not Predefined)"

	switch code {
	case Input:
		name = "Input"
	case SensitiveInput:
		name = "Sensitive Input"
	case Success:
		name = "Success"
	case RedirectTemporary:
		name = "Redirect (Temporary)"
	case RedirectPermanent:
		name = "Redirect (Permanent)"
	case TemporaryFailure:
		name = "Temporary Failure"
	case ServerUnavailable:
		name = "Server Unavailable"
	case CGIError:
		name = "CGI Error"
	case ProxyError:
		name = "Proxy Error"
	case SlowDown:
		name = "Slow Down"
	case PermanentFailure:
		name = "Permanent Failure"
	case NotFound:
		name = "Not Found"
	}

	return fmt.Sprintf("%d - %s", code, name)
}
