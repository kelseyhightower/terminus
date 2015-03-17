// Copyright (c) 2014 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/kelseyhightower/facter/facts/coreos"
)

var m = make(map[string]interface{})

func main() {
	coreosFacts := coreos.Run()
	data, err := json.MarshalIndent(&coreosFacts, " ", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", data)
}
