package main

import (
	"log"
	"net/http"
)

// // MyHandler handles
// type MyHandler struct {
// }

// // ServeHTTP does the file serve
// func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	path := r.URL.Path[1:]
// 	log.Println(path)

// 	data, err := ioutil.ReadFile(string(path))

// 	if err == nil {
// 		w.Write(data)
// 	} else {
// 		w.WriteHeader(404)
// 		w.Write([]byte("404 Something went wrong - " + http.StatusText(404)))
// 	}
// }

func main() {
	fs := http.FileServer(http.Dir("data/source"))
	http.Handle("/", fs)

	log.Println("Serving weather data...")
	http.ListenAndServe(":7000", nil)
}
