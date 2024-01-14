package main

import (
	"github.com/google/uuid"
	"crypto/sha256"
	"github.com/pkg/errors"
	"net/url"
	"math/rand"
	"strings"
	"fmt"
	"net/http"
	"time"
)

func GenRandData(t string, l int) string {

	var letters string

	switch t {
	case "txt":
		letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	case "num":
		letters = "1234567890"
	case "uuid":
		return uuid.New().String()
	case "mixed":
		letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	case "dict":
		letters = "xz"
	default:
		letters = "10"
	}

	dict := []rune(letters)
	res := make([]rune, l)
	for i := range res {
		res[i] = dict[rand.Intn(len(dict))]
	}
	return string(res)

}

func GenerateID(s ...string) string {
	str := strings.Join(s, "")
	if len(s) == 0 {
		str = time.Now().String() + salt
	}

	h := sha256.New()
	h.Write([]byte(str))
	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs)
}

func getUrlKey(req *http.Request, key string) (string, error) {
	u, err := url.Parse(req.URL.String())
	if err != nil {
		return "", errors.Wrap(err, "url parsing error")
	}
	m, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "", errors.Wrap(err, "url query parsing error")
	}
	val, ok := m[key]
	if ok {
		if val[0] != "" {
			return val[0], nil
		} else {
			return "", errors.New("Empty " + key)
		}
	} else {
		return "", errors.New("No " + key)
	}

}