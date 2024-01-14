package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
)

type LogInfo struct {
	clientIp  string
	reqURL    string
	reqMethod string
	userAgent string
}

func (l *LogInfo) New(req *http.Request) {

	l.clientIp = req.Header.Get("X-FORWARDED-FOR")
	if l.clientIp == "" {
		l.clientIp = req.RemoteAddr
	}

	l.userAgent = req.Header.Get("User-Agent")
	if l.userAgent == "" {
		l.userAgent = "NO-USER-AGENT"
	}

	l.reqURL = req.URL.String()
	l.reqMethod = req.Method

}

func (l *LogInfo) Log() {
	data := fmt.Sprintf("%s - %s - %s - %s\n",l.clientIp, l.reqMethod, l.reqURL, l.userAgent)
	log.Printf("%s",data)
	data = time.Now().Format("2006.01.02 15:04:05.000000 ") + data
	writeToFile("log.txt", data)

}

func writeToFile(filename, data string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err, "error open log file: ")
	}
	defer f.Close()
	_,err = f.WriteString(data)
	if err != nil {
		return errors.Wrap(err, "error write log file:")
	}
	return nil
}