package libgemini

import "fmt"

type StatusCode uint16

const (
	Unset                      StatusCode = 0
	Input                      StatusCode = 10
	SensitiveInput             StatusCode = 11
	Success                    StatusCode = 20
	RedirectTemporary          StatusCode = 30
	RedirectPermanent          StatusCode = 31
	TemporaryFailure           StatusCode = 40
	ServerUnavailable          StatusCode = 41
	CGIError                   StatusCode = 42
	ProxyError                 StatusCode = 43
	SlowDown                   StatusCode = 44
	PermanentFailure           StatusCode = 50
	NotFound                   StatusCode = 51
	ClientCertificatedRequired StatusCode = 60
	CertificateNotAuthorized   StatusCode = 61
	CertificateNotValid        StatusCode = 62
)

func (code StatusCode) IsSuccess() bool {
	return code >= Success && code < RedirectTemporary
}

func (code StatusCode) String() string {
	name := ""

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
	case ClientCertificatedRequired:
		name = "Client Certificate Required"
	case CertificateNotAuthorized:
		name = "Certificate Not Authorized"
	case CertificateNotValid:
		name = "Certificate Not Valid"
	default:
		name = "Unknown"
	}

	return fmt.Sprintf("%d - %s", code, name)
}
