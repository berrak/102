/*
Copyright 2018 The 101 Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestMain(t *testing.T) {

	// Be nice save Stdout
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	main()
	w.Close()
	byteStdOutArray, _ := ioutil.ReadAll(r)
	// Restore stdout
	os.Stdout = rescueStdout

	got := fmt.Sprintf("%s", byteStdOutArray)
	// Want's LF also
	want := fmt.Sprintln("Hello 123")
	if got != want {
		t.Errorf("Want %s but got %s", want, got)
	}
}
