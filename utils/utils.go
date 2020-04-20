package utils

import (
	"fmt"
	"os"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"reflect"
	"time"
	"strings"
	"github.com/superirale/sipserver/users"

)

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func SetField(v interface{}, name string, value string) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return errors.New("v must be pointer to struct")
	}
	rv = rv.Elem()
	fv := rv.FieldByName(name)
	if !fv.IsValid() {
		return fmt.Errorf("not a field name: %s", name)
	}
	if !fv.CanSet() {
		return fmt.Errorf("cannot set field %s", name)
	}
	if fv.Kind() != reflect.String {
		return fmt.Errorf("%s is not a string field", name)
	}
	fv.SetString(value)
	return nil
}

func GenerateNonce(clientIp string) string {
	// recommended implementation from ietf rfc2069
	privateKey := "the_8mile_cookie"
	timeStamp := time.Now().String()
	str := clientIp + ":" + timeStamp + ":" + privateKey
	nonce := getMD5Hash(str)

	return nonce
}

func VerifyDigestAuth(user *users.User, digest map[string]string, method string) bool {

	h1 := digest["username"] + ":" + digest["realm"] + ":" + user.Password
	h1h := getMD5Hash(h1)

	h2 := method + ":" + digest["uriText"]
	h2h := getMD5Hash(h2)
	responseStr := h1h + ":" + digest["nonce"] + ":" + h2h
	response := getMD5Hash(responseStr)

	validity := false

	respStr := strings.Trim(digest["response"], "\r\n")

	if response == respStr {
		validity = true
	}

	return validity
}

// PrettyPrintMap pretty print map data
func PrettyPrintMap(data map[string]interface{})  {
	// map[string]interface{}{"a": 1, "b": 2}
	x := data
	b, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b))
}