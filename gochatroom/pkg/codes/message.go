package codes

var codeMessage = map[int]string{
	Success:         "success",
	Error:           "unknown error",
	PageNotFound:    "page not found",
	ErrorParameter:  "parameter error",
	ErrorValidation: "validation error",
}

func GetCodeMessage(code int) string {
	if msg, ok := codeMessage[code]; ok {
		return msg
	}
	return codeMessage[Error]
}
