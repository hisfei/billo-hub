package helper

import "fmt"

// CodeDetail encapsulates an error code and its corresponding static description.
type CodeDetail struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}

// WithMessage creates a new CodeDetail instance,
// which appends dynamic context information to the original message.
// This does not modify the original global CodeDetail variable.
func (cd CodeDetail) WithMessage(dynamicMsg string) CodeDetail {
	return CodeDetail{
		Code: cd.Code,
		Msg:  fmt.Sprintf("%s: %s", cd.Msg, dynamicMsg),
	}
}

// APIResponse defines the standard API response structure.
type APIResponse struct {
	CodeDetail
	Body interface{} `json:"body,omitempty"` // omitempty prevents the field from being displayed if the body is nil
}

// NewSuccessResponse creates an API response that indicates success.
func NewSuccessResponse(body interface{}) *APIResponse {
	return &APIResponse{
		CodeDetail: OK,
		Body:       body,
	}
}

// NewErrorResponse creates an API response that indicates failure.
func NewErrorResponse(codeDetail CodeDetail, body interface{}) *APIResponse {
	return &APIResponse{
		CodeDetail: codeDetail,
		Body:       body,
	}
}

var (
	/*************************** Middleware Error Returns ****************************/
	ErrHeaderParam    = CodeDetail{401, "Authentication request parameter error"}
	ErrTimestamp      = CodeDetail{402, "Inconsistent with server time"}
	ErrUserAuth       = CodeDetail{403, "Failed to get identity information"}
	ErrBody           = CodeDetail{405, "Failed to get request body"}
	ErrSign           = CodeDetail{406, "Signature error"}
	ErrDecrypt        = CodeDetail{406, "Data decryption failed"}
	ErrUpdateAuthInfo = CodeDetail{406, "Failed to update verification information"}
	ErrSaveUserInfo   = CodeDetail{411, "Failed to store personal information"}
	ErrCreateRedisKey = CodeDetail{412, "Failed to create key"}

	ErrGetSign     = CodeDetail{407, "Failed to get signature"}
	ErrTokenExpire = CodeDetail{408, "Token expired"}
	ErrToken       = CodeDetail{409, "Token error"}
	ErrComm        = CodeDetail{410, "General error"}

	/*************************** General ****************************/

	OK      = CodeDetail{200, "Request successful"}
	ErrOK   = CodeDetail{0, "Correct"}
	ErrCode = CodeDetail{10000000, "Code error"}

	ErrInner      = CodeDetail{10000003, "Internal error"}
	ErrWrongParam = CodeDetail{10000004, "Invalid parameter"}
	ErrSystemBusy = CodeDetail{10000005, "System is busy, please try again later"}

	/****** Upload ********/
	ErrUploadSize = CodeDetail{11000001, "File size is too large"}
	ErrUpload     = CodeDetail{11000002, "Upload error"}
	ErrUploadType = CodeDetail{11000003, "Incorrect upload data type"}
	ErrFilePath   = CodeDetail{11000004, "Failed to create file path"}
	ErrFileSave   = CodeDetail{11000005, "File saving error"}
)
