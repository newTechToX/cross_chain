package config

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestDir(t *testing.T) {
	files, _ := ioutil.ReadDir("./")

	for _, f := range files {
		fmt.Println(f.Name())
	}
}
