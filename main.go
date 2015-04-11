// Copyright (c) 2014 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
)

var (
	externalFactsDir string
	format           string
	formatFile       string
	daemon           bool
)

func init() {
	log.SetFlags(0)
	flag.BoolVar(&daemon, "daemon", false, "Run in daemon mode.")
	flag.StringVar(&externalFactsDir, "external-facts-dir", defaultExternalFacts, "Path to external facts directory.")
	flag.StringVar(&format, "format", "", "Format the output using the given go template.")
	flag.StringVar(&formatFile, "format-file", "", "Format the output using the given go template file.")
}

var defaultExternalFacts = "/etc/terminus/facts.d"

func main() {
	flag.Parse()

	if daemon {
		http.Handle("/facts", httpHandler(factsHandler))
		log.Fatal(http.ListenAndServe(":8080", nil))
	}

	f := getFacts()

	if format != "" {
		tmpl, err := template.New("format").Parse(format)
		if err != nil {
			log.Fatal(err)
		}
		err = tmpl.Execute(os.Stdout, &f.Facts)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	if formatFile != "" {
		tmpl, err := template.ParseFiles(formatFile)
		if err != nil {
			log.Fatal(err)
		}
		err = tmpl.Execute(os.Stdout, &f.Facts)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	data, err := json.MarshalIndent(&f.Facts, " ", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", data)
}
