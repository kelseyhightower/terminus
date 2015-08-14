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
	"strconv"
	"strings"
	"text/template"
)

var (
	externalFactsDir string
	format           string
	formatFile       string
	httpAddr         string
	printVersion     bool
)

func init() {
	log.SetFlags(0)
	flag.StringVar(&externalFactsDir, "external-facts-dir", defaultExternalFacts, "Path to external facts directory.")
	flag.StringVar(&format, "format", "", "Format the output using the given go template.")
	flag.StringVar(&formatFile, "format-file", "", "Format the output using the given go template file.")
	flag.StringVar(&httpAddr, "http", "", "HTTP service address (e.g., ':6060')")
	flag.BoolVar(&printVersion, "version", false, "print version and exit")
}

var defaultExternalFacts = "/etc/terminus/facts.d"

func main() {
	flag.Parse()

	if printVersion {
		fmt.Printf("terminus %s\n", Version)
		os.Exit(0)
	}

	if httpAddr != "" {
		http.Handle("/facts", httpHandler(factsHandler))
		log.Fatal(http.ListenAndServe(httpAddr, nil))
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

	// If there are arguments left over, use the first argument as a fact query.
	if len(flag.Args()) > 0 {
		var data interface{}
		var value interface{}
		path := flag.Args()[0]
		path_pieces := strings.Split(path, ".")

		// Convert the Fact structure into a generic interface{}
		// by first converting it to JSON and then unmarshaling it.
		j, err := json.Marshal(&f.Facts)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(j, &data)
		if err != nil {
			log.Fatal(err)
		}

		// Walk through the given path.
		// If there's a result, print it.
		for _, p := range path_pieces {
			i, err := strconv.Atoi(p)
			if err != nil {
				if _, ok := data.(map[string]interface{}); ok {
					value = data.(map[string]interface{})[p]
				}
			} else {
				if _, ok := data.([]interface{}); ok {
					if len(data.([]interface{})) >= i {
						value = data.([]interface{})[i]
					}
				}
			}
			data = value
		}
		if value != nil {
			// if the value is an []interface{} or map[string]interface{}, convert it to JSON
			if _, ok := value.([]interface{}); ok {
				if j, err := json.Marshal(value); err == nil {
					fmt.Println(string(j))
				}
			} else {
				if _, ok := value.(map[string]interface{}); ok {
					if j, err := json.Marshal(value); err == nil {
						fmt.Println(string(j))
					}
				} else {
					fmt.Println(value)
				}
			}
		}
		os.Exit(0)
	}

	data, err := json.MarshalIndent(&f.Facts, " ", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", data)
}
