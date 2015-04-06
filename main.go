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

	"github.com/kelseyhightower/terminus/facts/system"
)

var (
	externalFactsDir string
)

func init() {
	flag.StringVar(&externalFactsDir, "external-facts-dir", defaultExternalFacts, "path to external facts directory")
}

var defaultExternalFacts = "/etc/terminus/facts.d"

var m map[string]interface{}

func main() {
	flag.Parse()

	log.SetFlags(0)
	systemFacts := system.Run()
	m = make(map[string]interface{})
	m["system"] = systemFacts

	processExternalFacts(externalFactsDir)
	data, err := json.MarshalIndent(&m, " ", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", data)
}

func processExternalFacts(factsDir string) {
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

	for _, fact := range staticFacts {
		data, err := ioutil.ReadFile(fact)
		if err != nil {
			log.Println(err)
			continue
		}

		var result interface{}
		err = json.Unmarshal(data, &result)
		if err != nil {
			log.Println(err)
			continue
		}
		m[strings.TrimSuffix(filepath.Base(fact), ".json")] = result
	}

	for _, fact := range executableFacts {
		out, err := exec.Command(fact).Output()
		if err != nil {
			log.Println(err)
			continue
		}
		var result interface{}
		err = json.Unmarshal(out, &result)
		if err != nil {
			log.Println(err)
			continue
		}
		m[filepath.Base(fact)] = result
	}
}

func isExecutable(fi os.FileInfo) bool {
	if m := fi.Mode(); !m.IsDir() && m&0111 != 0 {
		return true
	}
	return false
}
