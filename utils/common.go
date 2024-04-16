package utils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/jinzhu/now"
	uuid "github.com/satori/go.uuid"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
)

func Or(s ...string) string {
	if len(s) == 0 {
		return ""
	}
	for _, str := range s[:len(s)-1] {
		if str != "" {
			return str
		}
	}
	return s[len(s)-1]
}

func DecodeURIComponent(str string) string {
	result, _ := url.QueryUnescape(str)

	return result
}

func MD5Hash(data string) string {
	md5hash := md5.New()
	md5hash.Write([]byte(data))

	// convert the hash value to a string
	return hex.EncodeToString(md5hash.Sum(nil))
}

func Salt() string {
	return MD5Hash(Cfg.AppSecret + now.BeginningOfMonth().String())
}

func VisitSalt() string {
	return MD5Hash(Cfg.AppSecret + now.BeginningOfHour().String())
}

func GenerateUUIDWithSeed(seed string) string {
	// generate the MD5 hash
	md5string := MD5Hash(seed)

	// generate the UUID from the
	// first 16 bytes of the MD5 hash
	u, err := uuid.FromBytes([]byte(md5string[0:16]))
	if err != nil {
		log.Fatal(err)
	}

	return u.String()
}

func UUID(args ...string) string {
	seed := strings.Join(args, ".")

	return GenerateUUIDWithSeed(seed + Salt())
}

func SoftTouch(data interface{}) string {
	if data == nil {
		return ""
	}
	if dataSlice, ok := data.([]string); ok && len(dataSlice) > 0 {
		return dataSlice[0]
	}
	return ""
}

// getIP returns the ip address from the http request
func GetIP(r *http.Request) (string, error) {
	ips := r.Header.Get("X-Forwarded-For")
	splitIps := strings.Split(ips, ",")

	if len(splitIps) > 0 {
		// get last IP in list since ELB prepends other user defined IPs, meaning the last one is the actual client IP.
		netIP := net.ParseIP(splitIps[len(splitIps)-1])
		if netIP != nil {
			return netIP.String(), nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	netIP := net.ParseIP(ip)
	if netIP != nil {
		ip := netIP.String()
		if ip == "::1" {
			return "127.0.0.1", nil
		}
		return ip, nil
	}

	return "", errors.New("IP not found")
}
