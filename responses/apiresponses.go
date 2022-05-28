package response

import (
	"bytes"
	"text/template"
)

type ErrorCallResponse struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	MoreInfo string `json:"more_info"`
	Status   int    `json:"status"`
}

type ErrorTemplate struct {
	Host   string
	Number string
}

var VoiceValidation = map[string]ErrorCallResponse{
	"TO-REQUIRED": {
		Code:     21201,
		Message:  "No 'To' number is specified",
		MoreInfo: "https://{{.Host}}",
		Status:   400,
	},
	"FROM-REQUIRED": {
		Code:     21213,
		Message:  "No 'From' number is specified.",
		MoreInfo: "https://{{.Host}}",
		Status:   400,
	},
	"URL-REQUIRED": {
		Code:     21205,
		Message:  "Url parameter is required.",
		MoreInfo: "https://{{.Host}}",
		Status:   400,
	},
	"BAD-PLAY-URL": {
		Code:     00005,
		Message:  "Play url is not correct.",
		MoreInfo: "https://{{.Host}}",
		Status:   400,
	},
	"TO-INVALID": {
		Code:     21211,
		Message:  "The phone number you are attempting to call, {{.Number}}, is not valid.",
		MoreInfo: "https://{{.Host}}",
		Status:   400,
	},
	"FROM-INVALID": {
		Code:     21212,
		Message:  " \"From is not a valid phone number: {{.Number}}.\"",
		MoreInfo: "https://{{.Host}}",
		Status:   400,
	},
	"FROM-VERIFY": {
		Code: 21210,
		Message: "The source phone number provided, {{.Number}}, is not yet verified for your account. " +
			"You may only make calls from phone numbers that you've verified or purchased from Siprtc.",
		MoreInfo: "https://{{.Host}}",
		Status:   400,
	},
	"MALFORMED": {
		Code:     21210,
		Message:  "Request is not able to decode, Please check and send again",
		MoreInfo: "https://{{.Host}}",
		Status:   400,
	},
	"ACCOUNT-DISABLE": {
		Code:     10001,
		Message:  "Account is not active",
		MoreInfo: "https://{{.Host}}",
		Status:   400,
	},
	"TO-UNVERIFIED": {
		Code:     21219,
		Message:  "The number {{.Number}} is unverified. Trial accounts may only make calls to verified numbers.",
		MoreInfo: "https://{{.Host}}",
		Status:   400,
	},
	"PAYMENT-REQUIRED": {
		Code:     000402,
		Message:  "You don't have sufficient balance, please recharge",
		MoreInfo: "https://{{.Host}}",
		Status:   402,
	},
	"FROM-SIP-USER-REQUIRED": {
		Code:     21212,
		Message:  "The caller address is not a user name: {{.Number}}",
		MoreInfo: "https://{{.Host}}",
		Status:   400,
	},
	"DELETE-IN-PROGRESS": {
		Code:     20009,
		Message:  "Call cannot be deleted because it is still in progress.",
		MoreInfo: "https://{{.Host}}",
		Status:   409,
	},
	"CALL-UPDATE-FAILED-CDR": {
		Code:     10000,
		Message:  "Status Cannot update a completed call.",
		MoreInfo: "https://{{.Host}}",
		Status:   400,
	},
	"RESOURCE-CALL-NOT-FOUND": {
		Code:     20404,
		Message:  "The requested resource {{.Number}} was not found",
		MoreInfo: "https://{{.Host}}",
		Status:   404,
	},
	"BODY-PARAM-INCORRECT": {
		Code:     20001,
		Message:  "{{.Number}} is not a valid choice",
		MoreInfo: "https://{{.Host}}",
		Status:   400,
	},
}

func GetErrorStruct(state string, errorTemplate interface{}) ErrorCallResponse {
	validationError := VoiceValidation[state]
	validationError.Message = parseTemplate(validationError.Message, errorTemplate)
	validationError.MoreInfo = parseTemplate(validationError.MoreInfo, errorTemplate)
	return validationError
}

func parseTemplate(tmpl string, errorTemplate interface{}) string {
	var err error
	// with name passed as argument
	voiceValidationTmpl := template.New("VoiceValidation")
	// "Parse" parses a string into a template
	voiceValidationTmpl, err = voiceValidationTmpl.Parse(tmpl)
	if err != nil {
		return tmpl
	}
	var tmplBytes bytes.Buffer

	if err = voiceValidationTmpl.Execute(&tmplBytes, errorTemplate); err != nil {
		return tmpl
	}
	return tmplBytes.String()
}

/*

{
"code": "10000",
"message": "Status Cannot update a completed call.",
"more_info": "https://signalwire.com",
"status": 400
}*/
