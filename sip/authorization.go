package sip

import (
	"strings"
	uriLib "github.com/superirale/sipserver/uri"
	"github.com/superirale/sipserver/utils"
	// "fmt"
)

// Authorization struct
type Authorization struct {
	AuthType  string
	Username  string
	Realm     string
	Nonce     string
	Uri       uriLib.URI
	Response  string
	Algorithm string
	isAuthSet bool
}

// IsAuthSet method checks if isAuthSet is set
func (auth *Authorization) IsAuthSet() bool {
	return auth.isAuthSet
}

// BuildAuthorizationHeader function builds authorization header
func BuildAuthorizationHeader(auth string) Authorization {

	var authorize Authorization
	authorize.isAuthSet = true
	authArr := strings.Split(auth, " ")
	// username="usmanirale",realm="0.0.0.0:5060",nonce="8b78cfb87e909e14c61a6cceeb5a0c7c",uri="sip:192.168.8.102",response="b0b1d50e70616c8a73dc8a8da52b8f99",algorithm=MD5
	// fmt.Println(authArr[2])
	// fmt.Println(authArr)
	for _, prop := range authArr {

		prop = strings.Replace(prop, ",", "", 1)
		// fmt.Println(prop)
		isURIString := strings.Contains(prop, "uri")

		if isURIString {
			tags := make(map[string]string)
			uriText := strings.Replace(prop, "\"", "", 2)
			uriText = strings.Replace(uriText, "uri=", "", 1)
			// build Uri
			// uri="sip:192.168.8.101;transport=UDP"
			// fmt.Println(uriText)
			// fmt.Println(prop)
			uriArr := strings.Split(prop, ";")
			// uri="sip:192.168.8.101
			uriStrArr := strings.Split(prop, "=")
			// "sip:192.168.8.101
			uri := strings.Replace(uriStrArr[1], "\"", "", 1)

			if len(uriArr) > 1 {
				for c := 1; c < len(uriArr); c++ {
					// transport=UDP"
					tagArr := strings.Split(uriArr[c], "=")
					if len(tagArr) == 2 {
						key := tagArr[0]
						value := tagArr[1]
						tags[key] = strings.Replace(value, "\"", "", 1)
					}
				}
			}
			uriObj := uriLib.BuildURI(uri, uriText, tags)
			authorize.Uri = *uriObj
		} else {
			propArr := strings.Split(prop, "=")

			if len(propArr) > 1 {

				propTitleCase := strings.Title(propArr[0])
				if propTitleCase != "uri" {
					propValue := strings.Replace(propArr[1], "\"", "", 2)
					utils.SetField(&authorize, propTitleCase, propValue)
				}

			}
		}
	}
	return authorize
}