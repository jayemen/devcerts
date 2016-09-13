//go:generate go run assets/assets_generate.go

package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jayemen/devcerts/assets"
	"github.com/jayemen/devcerts/certutil"
)

type certRequest struct {
	CommonName string   `json:"commonName"`
	Names      []string `json:"names"`
	Ips        []string `json:"ips"`
}

type caHandler struct {
	store   certutil.Store
	handler func(certutil.Store, http.ResponseWriter, *http.Request)
}

func (c caHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.handler(c.store, w, r)
}

func bail(w http.ResponseWriter, e error, code int) {
	log.Printf("%s", e.Error())
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(code)
	w.Write([]byte(e.Error()))
}

func certReqHandler(s certutil.Store, w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		bail(w, err, 400)
		return
	}

	data := r.Form.Get("data")
	var request certRequest
	err = json.Unmarshal([]byte(data), &request)
	if err != nil {
		bail(w, err, 400)
		return
	}

	cert, err := s.New(request.CommonName, request.Names, request.Ips)
	if err != nil {
		bail(w, err, 500)
		return
	}

	w.Header().Add("Content-Type", "application/zip")
	w.Header().Add("Content-Disposition", "attachment; filename=cert.zip")
	w.WriteHeader(200)
	err = zipContents(cert, w)
	if err != nil {
		log.Printf("Error: %v", err)
	}
}

func createCa() (certutil.Store, error) {
	cert, err := ioutil.ReadFile("ca.crt")
	if err != nil {
		return certutil.Store{}, err
	}

	key, err := ioutil.ReadFile("ca.key")
	if err != nil {
		return certutil.Store{}, err
	}

	return certutil.Load(cert, key)
}

func zipContents(cert certutil.Store, w io.Writer) error {
	zipper := zip.NewWriter(w)
	file, err := zipper.Create("cert.crt")
	if err != nil {
		return err
	}

	err = cert.WriteCert(file)
	if err != nil {
		return err
	}

	err = cert.WriteIntermediates(file)
	if err != nil {
		return err
	}

	file, err = zipper.Create("cert.key")
	if err != nil {
		return err
	}

	err = cert.WriteKey(file)
	if err != nil {
		return err
	}

	file, err = zipper.Create("ca.crt")
	if err != nil {
		return err
	}

	err = cert.WriteRoot(file)
	if err != nil {
		return err
	}

	err = zipper.Close()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	ca, err := createCa()
	if err != nil {
		log.Fatal(err)
		return
	}

	r := mux.NewRouter()

	r.Path("/create-certificate").
		Methods("POST").
		Handler(caHandler{ca, certReqHandler})

	r.PathPrefix("/").
		Handler(http.FileServer(assets.Assets))

	port := os.Getenv("PORT")

	if port == "" {
		port = "80"
	}

	fmt.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe("127.0.0.1:"+port, r))
}
