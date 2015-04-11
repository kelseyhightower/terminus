package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"text/template"
)

type httpError struct {
	Error   error
	Message string
	Code    int
}

type httpHandler func(http.ResponseWriter, *http.Request) *httpError

func (fn httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		log.Println(err)
		http.Error(w, err.Message, err.Code)
	}
}

func factsHandler(w http.ResponseWriter, r *http.Request) *httpError {
	f := getFacts()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return &httpError{err, "Can't process template string", 500}
	}
	defer r.Body.Close()

	if string(body) != "" {
		tmpl, err := template.New("format").Parse(string(body))
		if err != nil {
			return &httpError{err, "Can't process template string", 500}
		}
		err = tmpl.Execute(w, &f.Facts)
		if err != nil {
			return &httpError{err, "Can't process template string", 500}
		}
		return nil
	}

	data, err := json.MarshalIndent(&f.Facts, " ", "  ")
	if err != nil {
		return &httpError{err, "Error processing facts", 500}
	}
	w.Header().Set("Server", "Teminus 1.0.0")
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
	return nil
}
