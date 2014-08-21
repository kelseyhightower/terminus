// Copyright (c) 2014 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/kelseyhightower/facter/facts/coreos"
	"gopkg.in/yaml.v1"
)

var (
	outputYAML bool
)

func init() {
	flag.BoolVar(&outputYAML, "yaml", false, "Output YAML")
}

func main() {
	flag.Parse()
	facts := make(map[string]interface{})
	facts["coreos"] = coreos.Run()
	var data []byte
	var err error
	if outputYAML {
		data, err = yaml.Marshal(facts)
	} else {
		data, err = json.MarshalIndent(facts, "", "  ")
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))
}
