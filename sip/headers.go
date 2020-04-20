package sip

import (
	"strings"

	"github.com/superirale/sipserver/utils"
)

// SIPHeaders SIP Request headers
type SIPHeaders struct {
	Contact         string
	CallID          string
	To              string
	From            string
	MinExpires      int
	Expires         int
	CSeq            CSeqStruct
	MaxForwards     int
	UserAgent       string
	Allow           string
	ContentLength   int
	Content         []interface{}
	WWWAuthenticate string
	Authorization   Authorization
	RecordRoute     string
}

// GenerateWWWAuthHeader function generate the auth challenge header
func GenerateWWWAuthHeader(data map[string]string) string {
	wwwHeader := "WWW-Authenticate: Digest "
	nonce := utils.GenerateNonce(data["clientIp"])
	nonceStr := "nonce=\"" + nonce + "\""
	realmStr := "realm=\"" + data["realm"] + "\","
	algoStr := ",algorithm=md5"
	wwwHeader += realmStr + nonceStr +algoStr + "\n"

	return wwwHeader
}

//GetHeaderValue Get Header values
func GetHeaderValue(header, s string, i int) string {
	str := strings.Split(header, s)
	return str[i]
}

// GetUserNameFromTag function extract username from the FROM tag
func GetUserNameFromTag(fromTag string) string {
	var username string
	// var res string
	//  reimplement with regex

	//<sip:usmanirale@192.168.8.102>;tag=bf39fb91
	str := strings.Split(fromTag, ";")

	// <sip:usmanirale@192.168.8.102>;
	if strings.Contains(str[0], "<sip:") {
		secStrArr := strings.Split(str[0], "<sip:")
		secStr := strings.Trim(secStrArr[1], ">")
		res := strings.Split(secStr, "@")
		username = res[0]

	} else {
		secStr := strings.Split(str[0], "sip:")
		res := strings.Split(secStr[0], "@")
		username = res[0]
	}
	return username
}
