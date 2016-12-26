package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

type Trans bool

var remote string

func (t Trans) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	url_string := remote + req.URL.Path + "?" + req.URL.RawQuery
	zombie_request, err := http.NewRequest("GET", url_string, nil)
	if err != nil {
		log.Fatal(err)
	}
	zombie_request.Header = req.Header
	zombie_request.Header.Del("Host")
	zombie_request.Header.Add("Host", "www.google.com")
	zombie_request.Header.Del("Refer")
	zombie_request.Header.Add("Refer", "https://www.google.com")
	client := &http.Client{}

	res, err := client.Do(zombie_request)
	if err != nil {
		log.Fatal(err)
	}
	robots, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	actual_remote := res.Request.URL.Host
	remote = "https://" + actual_remote
	log.Println(res.Request.URL.String())
	log.Println("-------------------------------------------\n")

	for k, _ := range res.Header {
		var v string
		if k == "Host" {
			v = "www.google.com"
		} else {
			v = res.Header.Get(k)
		}

		w.Header().Set(k, v)
	}
	w.Write(robots)

}

func main() {
	var trans Trans
	remote = "https://www.google.com"
	log.Fatal(http.ListenAndServe(":80", trans))
}
