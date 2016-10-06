package main

import (
	"flag"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
)

var (
	hostport = flag.String("http", ":"+os.Getenv("PORT"), "")
	baseurl  = flag.String("base", os.Getenv("IMG_PROXY_BASEURL"), "")
)

func main() {
	flag.Parse()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		imgPath := r.RequestURI
		log.Printf("Serving image %v\n", *baseurl+imgPath)
		res, err := http.Get(*baseurl + imgPath)
		defer res.Body.Close()
		if err != nil {
			log.Printf("Could not get base image %v: %v\n", *baseurl+imgPath, err)
			http.Error(w, "Could not get base image", res.StatusCode)
			return
		}
		img, err := jpeg.Decode(res.Body)
		if err != nil {
			log.Printf("Could not decode image: %v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if err := png.Encode(w, img); err != nil {
			log.Printf("Could not encode image: %v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})
	log.Fatalln(http.ListenAndServe(*hostport, nil))
}
