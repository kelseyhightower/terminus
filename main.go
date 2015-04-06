// Copyright (c) 2014 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/kelseyhightower/terminus/facts"
)

var (
	externalFactsDir string
	format           string
	formatFile       string
)

func init() {
	log.SetFlags(0)
	flag.StringVar(&externalFactsDir, "external-facts-dir", defaultExternalFacts, "Path to external facts directory.")
	flag.StringVar(&format, "format", "", "Format the output using the given go template.")
	flag.StringVar(&formatFile, "format-file", "", "Format the output using the given go template file.")
}

var defaultExternalFacts = "/etc/terminus/facts.d"

func main() {
	flag.Parse()

	systemFacts := getSystemFacts()

	f := facts.New()
	f.Add("System", systemFacts)

	processExternalFacts(externalFactsDir, f)

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

func processExternalFacts(factsDir string, f *facts.Facts) {
	dir, err := os.Open(factsDir)
	if err != nil {
		log.Println(err)
		return
	}
	defer dir.Close()

	files, err := dir.Readdir(0)
	if err != nil {
		log.Println(err)
		return
	}

	executableFacts := make([]string, 0)
	staticFacts := make([]string, 0)

	for _, fi := range files {
		name := filepath.Join(factsDir, fi.Name())
		if isExecutable(fi) {
			executableFacts = append(executableFacts, name)
			continue
		}
		if strings.HasSuffix(name, ".json") {
			staticFacts = append(staticFacts, name)
		}
	}

	var wg sync.WaitGroup
	for _, p := range staticFacts {
		p := p
		wg.Add(1)
		go factsFromFile(p, f, &wg)
	}
	for _, p := range executableFacts {
		p := p
		wg.Add(1)
		go factsFromExec(p, f, &wg)
	}
	wg.Wait()
}

func factsFromFile(path string, f *facts.Facts, wg *sync.WaitGroup) {
	defer wg.Done()
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
		return
	}
	var result interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Println(err)
		return
	}
	f.Add(strings.TrimSuffix(filepath.Base(path), ".json"), result)
}

func factsFromExec(path string, f *facts.Facts, wg *sync.WaitGroup) {
	defer wg.Done()
	out, err := exec.Command(path).Output()
	if err != nil {
		log.Println(err)
		return
	}
	var result interface{}
	err = json.Unmarshal(out, &result)
	if err != nil {
		log.Println(err)
		return
	}
	f.Add(filepath.Base(path), result)
}

func isExecutable(fi os.FileInfo) bool {
	if m := fi.Mode(); !m.IsDir() && m&0111 != 0 {
		return true
	}
	return false
}
