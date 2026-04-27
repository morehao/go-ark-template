package code

import "github.com/morehao/golib/gerror"

const (
	OIDCClientIDRequiredError    = 101200
	OIDCClientInvalidError       = 101201
	OIDCRedirectURIMismatchError = 101202
	OIDCGenerateCodeError        = 101203
	OIDCInvalidScopeError        = 101204
	OIDCInvalidStateError        = 101205
	OIDCInvalidPKCEError         = 101206
	OIDCInvalidCodeVerifierError = 101207
	OIDCInvalidAuthCodeError     = 101208
	OIDCAuthCodeExpiredError     = 101209
	OIDCAuthCodeUsedError        = 101210
)

var oidcErrorMsgMap = gerror.CodeMsgMap{
	OIDCClientIDRequiredError:    "client_id is required",
	OIDCClientInvalidError:       "invalid client_id",
	OIDCRedirectURIMismatchError: "redirect_uri mismatch",
	OIDCGenerateCodeError:        "generate authorization code failed",
	OIDCInvalidScopeError:        "invalid scope",
	OIDCInvalidStateError:        "invalid state",
	OIDCInvalidPKCEError:         "invalid PKCE parameters",
	OIDCInvalidCodeVerifierError: "invalid code verifier",
	OIDCInvalidAuthCodeError:     "invalid authorization code",
	OIDCAuthCodeExpiredError:     "authorization code expired",
	OIDCAuthCodeUsedError:        "authorization code already used",
}