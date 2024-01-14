package main

import (
	"encoding/json"
	"fmt"
	
	"log"
	"net/http"
	"time"
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"sync"
)

type ResponseData struct {
	RequestId string `json:"requestid"`
	RandData  string `json:"data"`
}

func GenerateURI(w http.ResponseWriter, req *http.Request) {
	var l LogInfo
	defaultGenType := "txt"
	l.New(req)
	l.Log()

	switch req.Method {
	case "POST":
		genType, err := getUrlKey(req, "type")
		if err != nil {
			log.Println("err get gen type: ", err)
			genType = defaultGenType
		}

		len := 10
		genLenStr, err := getUrlKey(req, "len")
		if err != nil {
			log.Println("no gen length: ", err)
		} else {
			genLen, err := strconv.Atoi(genLenStr)
			if err != nil {
				log.Println("length not int: ", err)
			} else {
				len = genLen
			}
		}

		w.Header().Set("content-type", "application/json")

		var respData ResponseData
		respData.RequestId = GenerateID()
		respData.RandData = GenRandData(genType, len)

		data, err := json.Marshal(respData)
		if err != nil {
			log.Println("err marshaling: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if _, err := w.Write(data); err != nil {
			log.Println("err writing to client: ", err)
		}

		mu.Lock()
		storage[respData.RequestId] = respData.RandData
		mu.Unlock()

	default:
		log.Printf("method %s not allowed", req.Method)
		w.WriteHeader(http.StatusInternalServerError)

	}
}

func RetrieveURI(w http.ResponseWriter, req *http.Request) {
	var l LogInfo
	var r ResponseData
	l.New(req)
	l.Log()

	id, err := getUrlKey(req, "id")
	if err != nil {
		fmt.Println("err:", err)
	} else {
		r.RequestId = id
		r.RandData = storage[id]
	}
	w.Header().Set("content-type", "application/json")

	data, err := json.Marshal(r)
	if err != nil {
		log.Println("err marshaling: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(data); err != nil {
		log.Println("err writing to client: ", err)
	}


}

const salt = "896f3f70f61bc3fb19bef"

var storage = make(map [string]string)
var mu sync.Mutex
//https://github.com/avito-tech/pro-backend-trainee-assignment
func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/api/generate/", GenerateURI)
	mux.HandleFunc("/api/retrieve/", RetrieveURI)

	api := http.Server{
		Addr:         "localhost:8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("main: started on %s\n", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("error: listening and serving: %s", err)
	case <-shutdown:
		log.Println("main: shuting down")
		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("main: Graceful shutdown did not complete in %v : %v", timeout, err)
			err = api.Close()
		}
		if err != nil {
			log.Fatalf("main: could not stop server gracefully: %v", err)
		}

	}

}
