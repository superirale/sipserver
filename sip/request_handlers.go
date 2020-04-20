package sip

import (
	"github.com/superirale/sipserver/utils"
	"github.com/superirale/sipserver/users"
	"fmt"
)

// RegistrationHandler handles REGISTER Requests
func RegistrationHandler(r *SIPRequest, rmap map[string]string, addresses map[string]string, method string) string {

	var sipMessage string
	crlf := "\n"
	authCheck := false
	username := ""

	if r.RequestHeader.Authorization.Username == "" {
		username = GetUserNameFromTag(rmap["to"])
	} else {
		username = r.RequestHeader.Authorization.Username
	}

	aorMap := map[string]string{
		"username": username,
		"physicalAddress": addresses["remote"],
	}
	if r.RequestHeader.Expires == 0 {
		aor := NewAOR(aorMap)
		_,err := RemoveAOR(aor)
		utils.CheckError(err)
		return sipMessage
	}

	isAuthAvailable := r.RequestHeader.Authorization.IsAuthSet()

	if isAuthAvailable {
		users := users.MakeUsers()
		userName := r.RequestHeader.Authorization.Username
		user := users[userName]
		authObj :=r.RequestHeader.Authorization
		digestData := map[string]string{
			"realm": authObj.Realm,
			"nonce": authObj.Nonce,
			"response": authObj.Response,
			"username": authObj.Username,
			"uriText": authObj.Uri.UriText,
		}
		// mapss := map[string]interface{}{
		// 	"realm": authObj.Realm,
		// 	"nonce": authObj.Nonce,
		// 	"response": authObj.Response,
		// 	"username": authObj.Username,
		// 	"uriText": authObj.Uri.UriText,
		// }
		// utils.PrettyPrintMap(mapss)
		authCheck = utils.VerifyDigestAuth(user, digestData, method)
	}

	if isAuthAvailable && authCheck {

		aor := NewAOR(aorMap)
		isAORExists := isAORExists(aor)

		if isAORExists == false {
			_,err := SaveAOR(aor)
			utils.CheckError(err)
		}

		sipMessage += "SIP/2.0 200 OK" + crlf
	} else {
		data := make(map[string]string)
		data["realm"] = addresses["local"] // or use the domain name
		data["clientIp"] = addresses["remote"]
		sipMessage += "SIP/2.0 401 Unauthorized" + crlf
		sipMessage += GenerateWWWAuthHeader(data)
	}
	sipMessage += rmap["via"] + crlf
	sipMessage += rmap["to"] + crlf
	sipMessage += rmap["from"] + crlf
	sipMessage += rmap["call-id"] + crlf
	sipMessage += rmap["cseq"] + crlf
	sipMessage += rmap["expires"] + crlf
	if _, ok := rmap["supported"]; ok {
		sipMessage += rmap["supported"] + crlf
	}
	sipMessage += rmap["content-length"] + crlf
	sipMessage += crlf
	return sipMessage
}

// InviteHandler handles INVITE Requests
func InviteHandler(headers string) string {
	headersByte := []byte(headers)
	aor := GetAOR("superirale")
	response := ForwardConn(headersByte, aor[0], "udp")
	// fmt.Println(response)
	return response
}

// MessageHandler handle MESSAGE Requests
func MessageHandler(headers string) string  {
	headersByte := []byte(headers)
	aor := GetAOR("superirale")
	response := ForwardConn(headersByte, aor[0], "udp")
	fmt.Println(response)
	return response
}