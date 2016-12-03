package main

import (
	fmt "fmt"
	http "net/http"
	strings "strings"
	time "time"
)

const (
	listenAddress = ":8080"
)

var (
	sleepDuration = 60 * time.Second
)

func init() {
	fmt.Printf("Listening on %s ...\n", listenAddress)
}

func handler(w http.ResponseWriter, r *http.Request) {
	var queryStrings = make([]string, 0)

	for queryKey, queryValue := range r.URL.Query() {
		queryStrings = append(queryStrings, fmt.Sprintf("%s=%s", queryKey, queryValue))
	}

	fmt.Printf("Request  -  [%s %s?%s]  -  (Sleeping for %s ...)\n", r.Method, r.URL.Path, strings.Join(queryStrings, "&"), sleepDuration)
	time.Sleep(sleepDuration)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(listenAddress, nil)
}
