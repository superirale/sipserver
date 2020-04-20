package sip

import (
	"strconv"
	"strings"
)

type SIPRequest struct {
	RequestsMethod SIPRequestLine
	RequestHeader  SIPHeaders
}

type SIPRequestLine struct {
	Method     string
	SIPURI     string
	SIPVersion string
}

// RequestHandler handle all incoming sip request
func RequestHandler(request []byte, length int, addresses map[string]string) string {

	reqString := string(request[:length])
	headersVars := []string{"To:", "From:", "Via:", "Call-ID:",
		"Contact:", "CSeq:", "User-Agent:", "Max-Forwards:",
		"Allow:", "Content-Length:", "Expires:", "Authorization:",
		"Allow-Events:", "Supported:"}
	ReqHeadersVals := make(map[string]string)

	reqArr := strings.Split(reqString, "\n")

	cmd := ""

	if len(reqArr) > 0 {
		isValidSIPMessage := strings.Contains(reqArr[0], "SIP/2.0")
		if isValidSIPMessage {
			for _, headerVar := range headersVars {
				for i := 1; i < len(reqArr); i++ {

					if strings.Contains(reqArr[i], headerVar) {
						key := strings.Replace(headerVar, ":", "", 1)
						key = strings.ToLower(key)
						ReqHeadersVals[key] = reqArr[i]
					}
				}
			}

			sipReqHeaders := BuildSIPRequestHeaders(ReqHeadersVals)

			// this line is either request line or status line
			reqStatusLine := strings.Split(reqArr[0], " ")

			if reqStatusLine[0] == "SIP/2.0" && len(reqStatusLine[1]) == 3 {
				// process status line
			} else {
				// process request line
				//set requestline
				sipRequestLine := BuildSIPRequestLine(reqStatusLine)

				// setup request Object
				sipRequest := BuildSIPRequest(sipRequestLine, sipReqHeaders)

				method := GetMethodFromRequestLine(reqStatusLine[0])
				switch method {
				case "REGISTER":
					cmd = RegistrationHandler(sipRequest, ReqHeadersVals, addresses, method)

				case "INVITE":

					cmd = InviteHandler(reqString)
				}
			}
		}
	}
	return cmd
}

func BuildSIPRequest(rl *SIPRequestLine, h *SIPHeaders) *SIPRequest {
	sipRequest := new(SIPRequest)
	sipRequest.RequestsMethod = *rl
	sipRequest.RequestHeader = *h
	return sipRequest
}

func BuildSIPRequestHeaders(headers map[string]string) *SIPHeaders {

	sipHeaders := new(SIPHeaders)
	sipHeaders.To = GetHeaderValue(headers["to"], "To:", 1)
	sipHeaders.From = GetHeaderValue(headers["from"], "From:", 1)
	sipHeaders.Contact = GetHeaderValue(headers["contact"], "Contact:", 1)
	sipHeaders.CallID = GetHeaderValue(headers["call-id"], " ", 1)
	sipHeaders.UserAgent = GetHeaderValue(headers["user-agent"], "User-Agent:", 1)
	if _, ok := headers["allow"]; ok {
		sipHeaders.Allow = GetHeaderValue(headers["allow"], "Allow:", 1)
	}
	exp := GetHeaderValue(headers["expires"], ":", 1)
	exp = strings.Trim(exp, " ")
	exp = strings.Trim(exp, "\r\n")
	expires, _ := strconv.Atoi(exp)
	sipHeaders.Expires = expires

	if _, ok := headers["authorization"]; ok {
		authStr := GetHeaderValue(headers["authorization"], "Authorization:", 1)
		authObj := BuildAuthorizationHeader(authStr)
		sipHeaders.Authorization = authObj
	}

	maxForwards, _ := strconv.Atoi(GetHeaderValue(headers["max-forwards"], " ", 1))
	sipHeaders.MaxForwards = maxForwards
	contentLength, _ := strconv.Atoi(GetHeaderValue(headers["content-length"], " ", 1))
	sipHeaders.ContentLength = contentLength
	CseqObj := BuildCSeq(headers["cseq"])
	sipHeaders.CSeq = *CseqObj
	return sipHeaders
}

func BuildSIPRequestLine(strArr []string) *SIPRequestLine {
	SIPRequestLine := new(SIPRequestLine)
	SIPRequestLine.Method = strArr[0]
	SIPRequestLine.SIPURI = strArr[1]
	SIPRequestLine.SIPVersion = strArr[2]
	return SIPRequestLine
}

func GetMethodFromRequestLine(line string) string {

	lineArr := strings.Split(line, " ")
	method := ""

	if len(lineArr) > 0 {
		method = lineArr[0]
	}

	return method
}