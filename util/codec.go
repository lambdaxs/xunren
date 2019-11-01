package util

import (
    "crypto/md5"
    "encoding/base64"
    "encoding/hex"
    "fmt"
    "math/rand"
    "regexp"
    "time"
    cryptorand "crypto/rand"
)

func GenerateRandomBytes(n int) ([]byte, error) {
    b := make([]byte, n)
    _, err := cryptorand.Read(b)
    // Note that err == nil only if we read len(b) bytes.
    if err != nil {
        return nil, err
    }

    return b, nil
}

func GenerateRandomString(s int) (string, error) {
    b, err := GenerateRandomBytes(s)
    return base64.URLEncoding.EncodeToString(b), err
}

//mobile verify
func VerifyPhone(mobileNum string) bool {
    regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
    reg := regexp.MustCompile(regular)
    return reg.MatchString(mobileNum)
}

// create md5 str
func MD5Str(str string) string  {
    h := md5.New()
    h.Write([]byte(str))
    return hex.EncodeToString(h.Sum(nil))
}

// create base64 str
func Base64Str(str string) string {
    return base64.StdEncoding.EncodeToString([]byte(str))
}

// parse base64 str
func UnBase64Str(str string) (res string,err error) {
    buf,err := base64.StdEncoding.DecodeString(str)
    return string(buf),err
}

// create 4 captcha code
func CaptchaCode() string {
    return fmt.Sprintf("%04d", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000))
}
